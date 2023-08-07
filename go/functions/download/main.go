package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"tds/shared/configs"
	"tds/shared/responses"
	"tds/shared/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/minio/minio-go/v7"
)

func DownloadExport(c *fiber.Ctx) error {
	log.Println("Started download")
	filename := c.Params("filename")
	object, err := configs.MINIO.GetObject(context.Background(), configs.EnvExportBucketName(), filename, minio.GetObjectOptions{})
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(&responses.DownloadResponse{
			Status:  http.StatusNotFound,
			Message: "The requested does not exist or has not been exported.",
		})
	}
	defer object.Close()
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Set(fiber.HeaderContentType, "application/gzip")

	// Stream the object content as the response
	if _, err = io.Copy(c.Response().BodyWriter(), object); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&responses.DownloadResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to download file.",
		})
	}
	return nil
}

func GetDownloadModels(c *fiber.Ctx) error {
	bucketStruc, err := utils.GetBucketStructure(configs.MINIO, configs.EnvModelBucketName(), "")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&responses.DownloadResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to list buckets.",
		})
	}
	return c.Status(http.StatusOK).JSON(bucketStruc)
}

func GetDownloadExport(c *fiber.Ctx) error {
	bucketStruc, err := utils.GetBucketStructure(configs.MINIO, configs.EnvExportBucketName(), "")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&responses.DownloadResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to list buckets.",
		})
	}
	return c.Status(http.StatusOK).JSON(bucketStruc)
}

func GetZippedModel(c *fiber.Ctx) error {
	log.Println("Started download")
	modelName := c.Params("modelName")
	zippedModelName := c.Params("zippedModelName")
	object, err := configs.MINIO.GetObject(context.Background(), configs.EnvModelBucketName(), modelName+"/"+zippedModelName, minio.GetObjectOptions{})
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(&responses.DownloadResponse{
			Status:  http.StatusNotFound,
			Message: "The requested does not exist or has not been exported.",
		})
	}
	defer object.Close()
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", zippedModelName))
	c.Set(fiber.HeaderContentType, "application/gzip")

	// Stream the object content as the response
	if _, err = io.Copy(c.Response().BodyWriter(), object); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&responses.DownloadResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to download file.",
		})
	}
	return nil
}

func GetModelData(c *fiber.Ctx) error {
	log.Println("Started download")
	modelName := c.Params("modelName")
	trainingSet := c.Params("trainingSet")
	fileName := c.Params("filename")
	object, err := configs.MINIO.GetObject(context.Background(), configs.EnvModelBucketName(), modelName+"/"+trainingSet+"/"+fileName, minio.GetObjectOptions{})
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(&responses.DownloadResponse{
			Status:  http.StatusNotFound,
			Message: "The requested does not exist or has not been exported.",
		})
	}
	defer object.Close()
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	contentType := mime.TypeByExtension(filepath.Ext(fileName))
	c.Set(fiber.HeaderContentType, contentType)

	// Stream the object content as the response
	if _, err = io.Copy(c.Response().BodyWriter(), object); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&responses.DownloadResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to download file.",
		})
	}
	return nil
}

func main() {
	app := fiber.New()
	app.Use(cors.New())
	configs.VerifyBucketExists(context.Background(), configs.MINIO, configs.EnvExportBucketName())
	app.Get("/transfer/health", utils.GetHealth)
	app.Get("/transfer/export/:fileName", DownloadExport)
	app.Get("/transfer/export", GetDownloadExport)
	app.Get("/transfer/models", GetDownloadModels)
	app.Get("/transfer/models/:modelName/:zippedModelName", GetZippedModel)
	app.Get("/transfer/models/:modelName/:trainingSet/:filename", GetModelData)
	app.Listen(":8081")
}
