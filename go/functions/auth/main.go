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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var adminURIs = []string{
	"/api/requests",
	"/api/training-runs",
	"/api/export",
	"/api/train",
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

func InitAdmin() {
	adminApiKey := configs.EnvAdminApiKey()
	hash, err := bcrypt.GenerateFromPassword([]byte(adminApiKey), bcrypt.DefaultCost)

	if err != nil {
		log.Fatalln("Could not generate Hash for Admin key.")
	}

	userCollection := configs.GetCollection(configs.ConnectDB(), "users")

	result := userCollection.FindOne(context.Background(), bson.M{
		"role": "admin",
	})
	var adminUser models.UserData
	if err := result.Decode(&adminUser); err != nil {
		admin := models.UserData{
			Id:    primitive.NewObjectID(),
			Role:  models.ADMIN,
			Email: "henry.schwerdtner@web.de",
			Key:   string(hash),
		}
		userCollection.InsertOne(context.Background(), admin)
		return
	}

	comparingResult := bcrypt.CompareHashAndPassword([]byte(adminUser.Key), []byte(adminApiKey))
	if comparingResult != nil {
		userCollection.UpdateOne(context.Background(), bson.M{
			"role": "admin",
		}, bson.M{
			"$set": bson.M{
				"key": string(hash),
			},
		})
	}
}

func main() {
	InitAdmin()
	app := fiber.New()
	app.Use(cors.New())
	app.Use(ValidateToken)
	app.Listen(":8081")
}
