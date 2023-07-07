package main

import (
	"context"
	"net/http"
	"tds/shared/configs"
	"tds/shared/models"
	"tds/shared/responses"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/bson"
)

func GetTrainingRuns(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := configs.GetCollection(configs.ConnectDB(), "training-runs")
	var trainingRuns []models.TrainingRun
	cursor, _ := collection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var trainingRun models.TrainingRun
		cursor.Decode(&trainingRun)
		trainingRuns = append(trainingRuns, trainingRun)
	}
	return c.Status(http.StatusOK).JSON(responses.TrainingRunResponse{
		Status: http.StatusOK,
		Data:   trainingRuns,
	})
}

func GetTrainingRunsByModelName(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	modelName := c.Params("modelName")
	defer cancel()
	collection := configs.GetCollection(configs.ConnectDB(), "training-runs")
	var trainingRuns []models.TrainingRun
	cursor, _ := collection.Find(ctx, bson.M{
		"name": modelName,
	})
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var trainingRun models.TrainingRun
		cursor.Decode(&trainingRun)
		trainingRuns = append(trainingRuns, trainingRun)
	}
	return c.Status(http.StatusOK).JSON(responses.TrainingRunResponse{
		Status: http.StatusOK,
		Data:   trainingRuns,
	})
}

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Get("/training-runs", GetTrainingRuns)
	app.Get("/training-runs/:modelName", GetTrainingRunsByModelName)

	// logger.Info("Server Running")
	app.Listen(":8081")
}
