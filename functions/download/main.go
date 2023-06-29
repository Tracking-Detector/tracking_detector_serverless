package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"tds/shared/configs"
	"tds/shared/responses"

	"github.com/gofiber/fiber/v2"
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

func main() {
	app := fiber.New()
	configs.VerifyBucketExists(context.Background(), configs.MINIO, configs.EnvExportBucketName())
	app.Get("/download/export/:fileName", DownloadExport)
	app.Listen(":8081")
}
