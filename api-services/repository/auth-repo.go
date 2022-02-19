package repository

import (
	"context"
	"github.com/authentication-app-server/api-services/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type AuthRepoImp struct {
	Database *mongo.Database
}

func NewAuthRepo(conn *mongo.Client, dbName string) AuthRepo {
	return &AuthRepoImp{
		Database: conn.Database(dbName),
	}
}

func (m *AuthRepoImp) Register(email string, password string) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	user := models.User{
		Email:    email,
		Password: password,
		IsActive: true,
		CreatedAt: primitive.Timestamp{
			T: uint32(time.Now().Unix()),
		},
		UpdatedAt: primitive.Timestamp{
			T: uint32(time.Now().Unix()),
		},
	}

	res, err := m.Database.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	//Object id to string
	insertedID := res.InsertedID.(primitive.ObjectID).Hex()

	return insertedID, err
}

func (m *AuthRepoImp) Login(email string, password string) (models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	filter := bson.D{{"email", email}}
	defer cancel()

	var user models.User

	err := m.Database.Collection("users").FindOne(ctx, filter).Decode(&user)

	return user, err
}
