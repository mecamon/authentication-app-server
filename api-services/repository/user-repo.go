package repository

import (
	"context"
	"github.com/authentication-app-server/api-services/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserRepoImp struct {
	Database *mongo.Database
}

func NewUserRepo(conn *mongo.Client, dbName string) UserRepo {
	return &UserRepoImp{
		Database: conn.Database(dbName),
	}
}

func (u *UserRepoImp) UpdateUser(ID string, validUser models.User) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	id, _ := primitive.ObjectIDFromHex(ID)
	filter := bson.D{{"_id", id}}

	userToUpdate := models.User{
		Name:      validUser.Name,
		Bio:       validUser.Bio,
		Email:     validUser.Email,
		Password:  validUser.Password,
		Telephone: validUser.Telephone,
		PhotoURL:  validUser.PhotoURL,
		UpdatedAt: primitive.Timestamp{
			T: uint32(time.Now().Unix()),
		},
	}

	result, err := u.Database.Collection("users").UpdateOne(ctx, filter, bson.D{{"$set", userToUpdate}})
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

func (u *UserRepoImp) UserInfo(ID string) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	id, _ := primitive.ObjectIDFromHex(ID)
	filter := bson.D{{"_id", id}}

	userFromDB := models.User{}

	err := u.Database.Collection("users").FindOne(ctx, filter).Decode(&userFromDB)
	if err != nil {
		return nil, err
	}

	var userSecureData = struct {
		Name      string `bson:"name,omitempty" json:"name"`
		Bio       string `bson:"bio,omitempty" json:"bio"`
		Email     string `bson:"email,omitempty" json:"email"`
		Telephone string `bson:"telephone,omitempty" json:"telephone"`
		PhotoURL  string `bson:"photo_url,omitempty" json:"photoURL"`
	}{
		Name:      userFromDB.Name,
		Bio:       userFromDB.Bio,
		Email:     userFromDB.Email,
		Telephone: userFromDB.Telephone,
		PhotoURL:  userFromDB.PhotoURL,
	}

	return userSecureData, nil
}
