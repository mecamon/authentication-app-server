package auth

import (
	"encoding/json"
	"github.com/authentication-app-server/api-services/repository"
	"github.com/authentication-app-server/helpers"
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
		auth       Auth
		out        []byte
		errorMap   helpers.ErrorsMap
		statusCode int
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
	}

	signedToken, err := helpers.GenerateToken(auth.Email)
	tokenMap := map[string]string{"token": signedToken}

	hasError, out := helpers.CustomResponse(tokenMap, errorMap)

	if !hasError {
		statusCode = http.StatusOK
	} else {
		if statusCode == 0 {
			statusCode = http.StatusBadRequest
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(out)
}

func (m *Handlers) register(w http.ResponseWriter, r *http.Request) {
	var auth Auth
	var statusCode int

	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		log.Println(err)
	}

	errorMap := auth.EvaluateNewUserCredentials()
	hashed, _ := helpers.HashPassword(auth.Password)
	result, err := m.AuthRepo.Register(auth.Email, string(hashed))
	if err != nil {
		if strings.Contains(err.Error(), "E11000") {
			errorMap.Message["email"] = "inserted email address is already taken"
			statusCode = 409
		}
	}

	resultMap := map[string]string{"inserted_id": result}
	hasError, out := helpers.CustomResponse(resultMap, errorMap)

	if !hasError {
		statusCode = http.StatusCreated
	} else {
		if statusCode == 0 {
			statusCode = http.StatusBadRequest
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(out)
}
