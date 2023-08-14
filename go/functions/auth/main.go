package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"tds/shared/configs"
	"tds/shared/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var adminURIs = []string{
	"/api/requests",
	"/api/training-runs",
	"/api/export",
	"/api/train",
	"/api/users",
}

func ValidateToken(c *fiber.Ctx) error {
	originalURI := c.Get("X-Original-URI")
	apiKey := c.Get("X-API-Key")
	isAdminUri := IsAdminURI(originalURI)
	res := ValidateApiKey(apiKey, isAdminUri)
	if res {
		return c.SendStatus(http.StatusOK)
	} else {
		return c.SendStatus(http.StatusForbidden)
	}
}

func ValidateApiKey(apiKey string, isAdmin bool) bool {
	split := strings.Split(apiKey, " ")
	if len(split) != 2 || split[0] != "Bearer" {
		return false
	}
	userCollection := configs.GetCollection(configs.ConnectDB(), configs.EnvUserCollection())
	cursor, err := userCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal("Failed to query MongoDB collection:", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var doc models.UserData
		if err := cursor.Decode(&doc); err != nil {
			log.Println("Failed to decode MongoDB document:", err)
			continue
		}
		if err := bcrypt.CompareHashAndPassword([]byte(doc.Key), []byte(split[1])); err == nil {
			if isAdmin {
				if doc.Role == models.ADMIN {
					return true
				}
			} else {
				return true
			}
		}
	}
	return false
}

func IsAdminURI(uri string) bool {
	for _, adminURI := range adminURIs {
		if strings.HasPrefix(uri, adminURI) {
			return true
		}
	}
	return false
}

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(ValidateToken)
	app.Listen(":8081")
}
