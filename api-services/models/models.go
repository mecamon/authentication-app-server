package models

import (
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty"`
	Name      string              `bson:"name,omitempty"`
	Bio       string              `bson:"bio,omitempty"`
	Email     string              `bson:"email,omitempty"`
	Password  string              `bson:"password,omitempty"`
	Telephone string              `bson:"telephone,omitempty"`
	PhoneURL  string              `bson:"phone_url,omitempty"`
	IsActive  bool                `bson:"is_active,omitempty"`
	CreatedAt primitive.Timestamp `bson:"created_at,omitempty"`
	UpdatedAt primitive.Timestamp `bson:"updated_at,omitempty"`
}

type CustomClaims struct {
	TokenType string
	Email     string `json:"email"`
	*jwt.RegisteredClaims
}
