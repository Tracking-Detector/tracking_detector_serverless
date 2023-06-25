package main

import (
	"context"
	"tds/shared/configs"
	"tds/shared/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var requestDataCollection *mongo.Collection = configs.GetCollection(configs.DB, "requests")

func ExportData(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cursor, _ := requestDataCollection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var requestData models.RequestData
		cursor.Decode(&requestData)

	}
}

func main() {
	app := fiber.New()

	app.Post("/export/204", ExportData)

	app.Listen(":8081")
}
