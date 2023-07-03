package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"tds/shared/configs"
	"tds/shared/extractor"
	"tds/shared/models"
	"tds/shared/responses"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var requestDataCollection *mongo.Collection = configs.GetCollection(configs.DB, "requests")

func RunDataExport(fe extractor.Extractor) {
	pr, pw := io.Pipe()

	gzipWriter := gzip.NewWriter(pw)
	// TODO write export data into mongo
	go func() {
		defer pw.Close()
		defer gzipWriter.Close()

		cursor, err := requestDataCollection.Find(context.Background(), bson.M{})
		if err != nil {
			log.Fatal("Failed to query MongoDB collection:", err)
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var doc models.RequestData
			if err := cursor.Decode(&doc); err != nil {
				log.Println("Failed to decode MongoDB document:", err)
				continue
			}
			encoded, encodeErr := fe.Encode(doc)
			if encodeErr != nil {
				continue
			}
			arr, err := json.Marshal(encoded)
			if err != nil {
				log.Fatal("Could not convert int[] to string")
				continue
			}
			data := strings.Trim(string(arr), "[]") + "\n"

			if _, err := gzipWriter.Write([]byte(data)); err != nil {
				log.Println("Failed to write to gzip writer:", err)
				break
			}
		}

		if cursor.Err() != nil {
			log.Println("Error occurred while iterating MongoDB cursor:", cursor.Err())
		}
	}()
	_, putErr := configs.MINIO.PutObject(context.Background(), configs.EnvExportBucketName(), fe.GetFileName(), pr, -1, minio.PutObjectOptions{
		ContentType: "application/gzip",
	})
	if putErr != nil {
		log.Fatal("Failed to upload data to MinIO:", putErr)
	}

	log.Println("Data compression and upload for " + fe.GetName() + " completed successfully!")

}

func ExportData(c *fiber.Ctx) error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	extractorName := c.Params("extractorName")
	var ext extractor.Extractor
	for _, ex := range extractor.EXTRACTORS {
		if ex.GetName() == extractorName {
			ext = ex
		}
	}
	if ext.GetName() == "" {
		return c.Status(http.StatusNotFound).JSON(responses.ExportJobStartResponse{
			Status:  http.StatusNotFound,
			Message: "Could not find the extractor you want to trigger.",
		})
	}
	go RunDataExport(ext)
	return c.Status(http.StatusOK).JSON(responses.ExportJobStartResponse{
		Status:  http.StatusOK,
		Message: "The export has been started.",
	})
}

func GetAllPossibleExports(c *fiber.Ctx) error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	possibleExports := make([]*fiber.Map, 0)
	for _, ext := range extractor.EXTRACTORS {
		possibleExports = append(possibleExports, &fiber.Map{
			"name":        ext.GetName(),
			"location":    ext.GetFileName(),
			"description": ext.GetDescription(),
		})
	}
	return c.Status(http.StatusOK).JSON(responses.ExportTypesResponse{
		Status: http.StatusOK,
		Data:   possibleExports,
	})
}

func RunAllExtractors() {
	for _, ex := range extractor.EXTRACTORS {
		RunDataExport(ex)
	}
}

func SetupCron() {
	c := cron.New()
	_, err := c.AddFunc("0 0 */14 * *", RunAllExtractors)
	if err != nil {
		log.Fatalln("Failed to add cron job", err)
	}
	c.Start()
}

func main() {
	SetupCron()
	app := fiber.New()
	configs.VerifyBucketExists(context.Background(), configs.MINIO, configs.EnvExportBucketName())
	app.Post("/export/:extractorName/run", ExportData)
	app.Get("/export", GetAllPossibleExports)

	app.Listen(":8081")
}
