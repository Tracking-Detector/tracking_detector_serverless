package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	ADMIN  Role = "admin"
	CLIENT Role = "client"
)

type UserData struct {
	Id    primitive.ObjectID `json:"_id,omitempty"`
	Role  Role               `json:"role"`
	Email string             `json:"email"`
	Key   string             `json:"key"`
}

type CreateUserData struct {
	Email string `json:"email"`
}

type UserDataRepresentation struct {
	Email string `json:"email"`
	Role  Role   `json:"role"`
}
