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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var requestDataCollection *mongo.Collection = configs.GetCollection(configs.DB, "requests")

func RunDataExport(fe extractor.Extractor) {
	pr, pw := io.Pipe()

	gzipWriter := gzip.NewWriter(pw)
	go func() {
		defer pw.Close()
		defer gzipWriter.Close()

		cursor, err := requestDataCollection.Find(context.Background(), nil)
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

			arr, err := json.Marshal(fe.Encode(doc))
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cursor, _ := requestDataCollection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var requestData models.RequestData
		cursor.Decode(&requestData)

	}
	return c.Status(http.StatusCreated).JSON(responses.RequestDataResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    &fiber.Map{"data": "adsasd"}})
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

func main() {

	app := fiber.New()

	app.Post("/export/:extractorName/run", ExportData)
	app.Get("/export", GetAllPossibleExports)

	app.Listen(":8081")
}
