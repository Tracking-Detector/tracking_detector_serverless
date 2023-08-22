package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"tds/shared/configs"
	"tds/shared/extractor"
	"tds/shared/models"
	"tds/shared/responses"
	"tds/shared/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	workers "github.com/jrallison/go-workers"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var modelCollection *mongo.Collection = configs.GetCollection(configs.DB, configs.EnvModelCollection())

func SetupWorkers() {
	workers.Configure(map[string]string{
		"server":   "redis:6379",
		"database": "0",
		"pool":     "30",
		"process":  fmt.Sprintf("worker-%d", os.Getpid()),
	})
}

func EnqueueExportJob(exportName string) {
	workers.Enqueue("exports", "Export", []string{exportName})
}

func EnqueueTrainingJob(modelName string, dataSetName string) {
	workers.Enqueue("training", "train_model", []string{modelName, dataSetName})
}

func RegisterExportJob(c *fiber.Ctx) error {
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
	EnqueueExportJob(ext.GetName())
	return c.Status(http.StatusOK).JSON(responses.ExportJobStartResponse{
		Status:  http.StatusOK,
		Message: "The export is queued.",
	})
}

func RegisterTrainingJob(c *fiber.Ctx) error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	modelName := c.Params("modelName")
	dataSetName := c.Params("dataSetName")
	var ext extractor.Extractor
	for _, ex := range extractor.EXTRACTORS {
		if ex.GetFileName() == dataSetName {
			ext = ex
		}
	}
	if ext.GetFileName() == "" {
		return c.Status(http.StatusNotFound).JSON(responses.ExportJobStartResponse{
			Status:  http.StatusNotFound,
			Message: "Could not find the extractor you want to trigger.",
		})
	}
	EnqueueTrainingJob(modelName, dataSetName)
	return c.Status(http.StatusOK).JSON(responses.ExportJobStartResponse{
		Status:  http.StatusOK,
		Message: "The training job is queued.",
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

func GetAllPossibleModels(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var modelDataValues []models.ModelData
	cursor, _ := modelCollection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var modelData models.ModelData
		cursor.Decode(&modelData)
		modelDataValues = append(modelDataValues, modelData)
	}

	return c.Status(http.StatusOK).JSON(responses.ModelDataResponse{
		Status: http.StatusOK,
		Data:   modelDataValues,
	})
}

func StartCronJobs() {
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("0 0 1,15 * *", func() {
		for _, export := range extractor.EXTRACTORS {
			EnqueueExportJob(export.GetName())
		}
	})
	if err != nil {
		log.Fatalf("Could not schedule training job: %v", err)
	}
	c.Start()
}

func main() {
	SetupWorkers()
	StartCronJobs()
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/dispatch/health", utils.GetHealth)
	app.Get("/dispatch/export", GetAllPossibleExports)
	app.Get("/dispatch/model", GetAllPossibleModels)
	app.Post("/dispatch/export/:extractorName", RegisterExportJob)
	app.Post("/dispatch/train/:modelName/run/:dataSetName", RegisterTrainingJob)

	app.Listen(":8081")
}
