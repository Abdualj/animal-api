package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Animal struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Species   string             `json:"species" bson:"species"`
	Breed     string             `json:"breed,omitempty" bson:"breed,omitempty"`
	Age       int                `json:"age" bson:"age"`
	Weight    float64            `json:"weight" bson:"weight"`
	Adopted   bool               `json:"adopted" bson:"adopted"`
	CreatedBy string             `json:"createdBy" bson:"createdBy"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// AnimalInput is the request body for POST and PUT, with validation rules.
type AnimalInput struct {
	Name    string  `json:"name" binding:"required,min=2,max=50"`
	Species string  `json:"species" binding:"required,oneof=dog cat bird rabbit horse fish reptile other"`
	Breed   string  `json:"breed" binding:"omitempty,max=50"`
	Age     *int    `json:"age" binding:"required,gte=0,lte=100"`
	Weight  float64 `json:"weight" binding:"required,gt=0,lte=2000"`
	Adopted bool    `json:"adopted"`
}

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username     string             `json:"username" bson:"username"`
	PasswordHash string             `json:"-" bson:"passwordHash"`
	Role         string             `json:"role" bson:"role"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
}

type Credentials struct {
	Username string `json:"username" binding:"required,alphanum,min=3,max=30"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)
