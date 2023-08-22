package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"strings"
	"tds/shared/configs"
	"tds/shared/extractor"
	"tds/shared/models"

	workers "github.com/jrallison/go-workers"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var requestDataCollection *mongo.Collection = configs.GetCollection(configs.DB, configs.EnvRequestCollection())

func SetupWorkers() {
	workers.Configure(map[string]string{
		"server":   "redis:6379",
		"database": "0",
		"pool":     "30",
		"process":  fmt.Sprintf("worker-%d", os.Getpid()),
	})
}

func Export(message *workers.Msg) {
	exportName, _ := message.Args().Array()
	log.WithFields(log.Fields{
		"service": "export",
	}).Info("Starting export with extractor: ", exportName[0])
	var ext extractor.Extractor
	for _, ex := range extractor.EXTRACTORS {
		if ex.GetName() == exportName[0] {
			ext = ex
		}
	}
	if ext.GetName() == "" {
		log.WithFields(log.Fields{
			"service": "export",
		}).Error("Error identifing export job: ", exportName[0])
		return
	}
	RunDataExport(ext)
}

func RunDataExport(fe extractor.Extractor) {
	pr, pw := io.Pipe()

	gzipWriter := gzip.NewWriter(pw)
	go func() {
		defer pw.Close()
		defer gzipWriter.Close()

		cursor, err := requestDataCollection.Find(context.Background(), bson.M{})
		if err != nil {
			log.WithFields(log.Fields{
				"service": "export",
				"error":   err.Error(),
			}).Error("Failed to query MongoDB collection.")
			return
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var doc models.RequestData
			if err := cursor.Decode(&doc); err != nil {
				log.WithFields(log.Fields{
					"service": "export",
					"error":   err.Error(),
				}).Error("Failed to decode MongoDB document.")
				continue
			}
			encoded, encodeErr := fe.Encode(doc)
			if encodeErr != nil {
				continue
			}
			arr, err := json.Marshal(encoded)
			if err != nil {
				log.WithFields(log.Fields{
					"service": "export",
					"error":   err.Error(),
				}).Error("Could not convert int[] to string.")
				continue
			}
			data := strings.Trim(string(arr), "[]") + "\n"

			if _, err := gzipWriter.Write([]byte(data)); err != nil {
				log.WithFields(log.Fields{
					"service": "export",
					"error":   err.Error(),
				}).Error("Failed to write to gzip writer.")
				break
			}
		}

		if cursor.Err() != nil {
			log.WithFields(log.Fields{
				"service": "export",
				"error":   cursor.Err().Error(),
			}).Error("Error occurred while iterating MongoDB cursor.")
		}
		log.WithFields(log.Fields{
			"service": "export",
		}).Info("Data compression and upload for " + fe.GetName() + " completed successfully.")
	}()
	_, putErr := configs.MINIO.PutObject(context.Background(), configs.EnvExportBucketName(), fe.GetFileName(), pr, -1, minio.PutObjectOptions{
		ContentType: "application/gzip",
	})
	if putErr != nil {
		log.WithFields(log.Fields{
			"service": "export",
			"error":   putErr.Error(),
		}).Error("Failed to upload data to MinIO.")
	}

}

func main() {
	configs.VerifyBucketExists(context.Background(), configs.MINIO, configs.EnvExportBucketName())
	SetupWorkers()
	workers.Process("exports", Export, 2)
	workers.Run()
}
