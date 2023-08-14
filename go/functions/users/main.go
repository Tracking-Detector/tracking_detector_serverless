package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"tds/shared/configs"
	"tds/shared/models"
	"tds/shared/responses"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
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

func CreateAPIKey(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var userData models.CreateUserData
	defer cancel()

	if err := c.BodyParser(&userData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.RequestDataResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()}})
	}
	userCollection := configs.GetCollection(configs.ConnectDB(), "users")
	result := userCollection.FindOne(ctx, bson.M{
		"email": userData.Email,
	})
	if result != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.RequestDataResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": "A user with this email address already exists."}})
	}
	key := uuid.New().String()
	hash, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.RequestDataResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": "Could not generate hash of api key."}})
	}
	newUser := models.UserData{
		Id:    primitive.NewObjectID(),
		Role:  models.CLIENT,
		Email: userData.Email,
		Key:   string(hash),
	}
	userCollection.InsertOne(ctx, newUser)
	return c.Status(http.StatusOK).JSON(responses.RequestDataResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &fiber.Map{"data": &fiber.Map{
			"key": string(hash),
		}}})

}

func GetUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userCollection := configs.GetCollection(configs.ConnectDB(), "users")
	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.RequestDataResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err}})
	}
	var userData []models.UserDataRepresentation
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user models.UserDataRepresentation
		cursor.Decode(&user)
		userData = append(userData, user)
	}
	fmt.Println(userData)
	return c.Status(http.StatusOK).JSON(responses.RequestDataResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": userData}})
}

func main() {
	InitAdmin()
	app := fiber.New()
	app.Use(cors.New())
	app.Get("/users", GetUsers)
	app.Post("/users", CreateAPIKey)
	app.Listen(":8081")
}
