package main

import (
	"context"
	"net/http"
	"tds/shared/configs"
	"tds/shared/extractor"
	"tds/shared/models"
	"tds/shared/responses"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var requestDataCollection *mongo.Collection = configs.GetCollection(configs.DB, "requests")
var featExtractor *extractor.Extractor = extractor.NewExtractor()

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

func main() {
	// init extractor
	featExtractor.URL(extractor.URL_EXTRACTOR)
	featExtractor.FrameType(extractor.FRAME_TYPE_EXTRACTOR)
	featExtractor.Method(extractor.METHOD_EXTRACTOR)
	featExtractor.Type(extractor.TYPE_EXTRACTOR)
	featExtractor.RequestHeaders(extractor.REQUEST_HEADER_REFERER_EXTRACTOR)
	featExtractor.Label(extractor.LABEL_EXTRACTOR)

	app := fiber.New()

	app.Post("/export/204", ExportData)

	app.Listen(":8081")
}
