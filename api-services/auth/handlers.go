package auth

import (
	"encoding/json"
	"github.com/authentication-app-server/api-services/repository"
	"github.com/authentication-app-server/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strings"
)

var handlers *Handlers

type Handlers struct {
	AuthRepo repository.AuthRepo
}

func NewHandlers(conn *mongo.Client, dbName string) *Handlers {
	handlers = &Handlers{
		AuthRepo: repository.NewAuthRepo(conn, dbName),
	}
	return handlers
}

func (m *Handlers) login(w http.ResponseWriter, r *http.Request) {

	var (
		auth     Auth
		errorMap helpers.ErrorsMap
	)

	err := json.NewDecoder(r.Body).Decode(&auth)

	result, err := m.AuthRepo.Login(auth.Email, auth.Password)
	hasAValidPass := helpers.CheckPassword(auth.Password, result.Password)

	if err != nil || !hasAValidPass {
		errorMap = helpers.ErrorsMap{
			Success: false,
			Message: map[string]string{
				"credentials": "incorrect email or password",
			},
		}
		_, output := helpers.CustomResponse(nil, errorMap)
		helpers.ResGenerator(w, http.StatusBadRequest, output)
		return
	}

	signedToken, err := helpers.GenerateToken(primitive.ObjectID.Hex(result.ID), result.Email)
	tokenMap := map[string]string{"token": signedToken}

	_, output := helpers.CustomResponse(tokenMap, errorMap)
	helpers.ResGenerator(w, http.StatusOK, output)
}

func (m *Handlers) register(w http.ResponseWriter, r *http.Request) {
	var auth Auth

	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		log.Println(err)
	}

	errorMap := auth.EvaluateNewUserCredentials()

	if len(errorMap.Message) > 0 {
		_, output := helpers.CustomResponse(nil, errorMap)
		helpers.ResGenerator(w, http.StatusBadRequest, output)
		return
	}

	hashed, _ := helpers.HashPassword(auth.Password)
	insertedID, err := m.AuthRepo.Register(auth.Email, string(hashed))
	if err != nil {
		if strings.Contains(err.Error(), "E11000") {
			errorMap.Message["email"] = "inserted email address is already taken"
			_, output := helpers.CustomResponse(nil, errorMap)
			helpers.ResGenerator(w, http.StatusConflict, output)
			return
		}
	}

	signedToken, err := helpers.GenerateToken(insertedID, auth.Email)

	resultMap := map[string]string{
		"token":       signedToken,
		"inserted_id": insertedID,
	}

	_, output := helpers.CustomResponse(resultMap, errorMap)
	helpers.ResGenerator(w, http.StatusCreated, output)
}
