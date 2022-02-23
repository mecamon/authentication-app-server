package users

import (
	"github.com/authentication-app-server/api-services/models"
	"github.com/authentication-app-server/api-services/repository"
	"github.com/authentication-app-server/helpers"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strings"
)

var handlers *Handlers

type Handlers struct {
	UserRepo repository.UserRepo
}

func NewHandlers(conn *mongo.Client, dbName string) *Handlers {
	handlers = &Handlers{
		UserRepo: repository.NewUserRepo(conn, dbName),
	}

	return handlers
}

func (m *Handlers) updateUserInfo(w http.ResponseWriter, r *http.Request) {

	var (
		contentType string
		fileSize    int64
	)

	errorMap := helpers.ErrorsMap{
		Success: false,
		Message: map[string]string{},
	}

	ID := r.Context().Value("ID").(string)

	err := r.ParseMultipartForm(128)
	if err != nil {
		log.Println("error parsing form:", err)
	}

	user := models.User{
		Name:      r.Form.Get("name"),
		Bio:       r.Form.Get("bio"),
		Email:     r.Form.Get("email"),
		Password:  r.Form.Get("password"),
		Telephone: r.Form.Get("telephone"),
	}

	errorMap = EvaluateEditUserCredentials(user)
	file, fileHeader, err := r.FormFile("file")

	if file != nil {
		contentType = fileHeader.Header.Get("Content-Type")
		fileSize = fileHeader.Size

		errorMap = EvaluateFile(contentType, fileSize, errorMap)
		if len(errorMap.Message) > 0 {
			_, output := helpers.CustomResponse(nil, errorMap)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(output)
			return
		}
	}

	hashed, err := helpers.HashPassword(user.Password)
	user.Password = string(hashed)
	
	modified, err := m.UserRepo.UpdateUser(ID, user)
	if err != nil {
		if strings.Contains(err.Error(), "E11000") {
			errorMap.Message["email"] = "inserted email address is already taken"
		}
		_, output := helpers.CustomResponse(nil, errorMap)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write(output)
		return
	}

	modifiedCountMap := map[string]int64{
		"modified": modified,
	}

	_, output := helpers.CustomResponse(modifiedCountMap, errorMap)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
