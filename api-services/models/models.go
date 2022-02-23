package models

import (
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name      string              `bson:"name,omitempty" json:"name"`
	Bio       string              `bson:"bio,omitempty" json:"bio"`
	Email     string              `bson:"email,omitempty" json:"email"`
	Password  string              `bson:"password,omitempty" json:"password"`
	Telephone string              `bson:"telephone,omitempty" json:"telephone"`
	PhotoURL  string              `bson:"photo_url,omitempty" json:"photoURL"`
	IsActive  bool                `bson:"is_active,omitempty" json:"is_active"`
	CreatedAt primitive.Timestamp `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt primitive.Timestamp `bson:"updated_at,omitempty" json:"updated_at"`
}

type CustomClaims struct {
	TokenType string
	ID        string `json:"id"`
	*jwt.RegisteredClaims
}
