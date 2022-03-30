package users

import (
	"github.com/authentication-app-server/api-services/models"
	"github.com/authentication-app-server/api-services/repository"
	"github.com/authentication-app-server/helpers"
	i18n_app "github.com/authentication-app-server/i18n-app"
	"github.com/authentication-app-server/services"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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

	lang := r.Header.Get("Accept-Language")
	locales := i18n_app.GetLocales(lang)

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

	errorMap = EvaluateEditUserCredentials(user, locales)
	if len(errorMap.Message) > 0 {
		_, output := helpers.CustomResponse(nil, errorMap)
		helpers.ResGenerator(w, http.StatusBadRequest, output)
		return
	}
	file, fileHeader, err := r.FormFile("file")

	if file != nil {
		defer file.Close()

		contentType = fileHeader.Header.Get("Content-Type")
		fileSize = fileHeader.Size

		errorMap = EvaluateFile(contentType, fileSize, errorMap, locales)

		imageURL, err := services.UploadImage(file, ID)
		if err != nil {
			errorMap.Message["file"] = err.Error()
		}

		if len(errorMap.Message) > 0 {
			_, output := helpers.CustomResponse(nil, errorMap)
			helpers.ResGenerator(w, http.StatusBadRequest, output)
			return
		}

		user.PhotoURL = imageURL
	}

	hashed, err := helpers.HashPassword(user.Password)
	user.Password = string(hashed)

	modified, err := m.UserRepo.UpdateUser(ID, user)
	if err != nil {
		if strings.Contains(err.Error(), "E11000") {
			errorMap.Message["email"] = locales.MustLocalize(&i18n.LocalizeConfig{MessageID: "EmailAddressTaken"})
		}
		_, output := helpers.CustomResponse(nil, errorMap)
		helpers.ResGenerator(w, http.StatusConflict, output)
		return
	}

	modifiedCountMap := map[string]int64{
		"modified": modified,
	}

	_, output := helpers.CustomResponse(modifiedCountMap, errorMap)
	helpers.ResGenerator(w, http.StatusOK, output)
	return
}

func (m *Handlers) userInfo(w http.ResponseWriter, r *http.Request) {

	errorMap := helpers.ErrorsMap{
		Success: false,
		Message: map[string]string{},
	}

	lang := r.Header.Get("Accept-Language")
	locales := i18n_app.GetLocales(lang)

	ID := r.Context().Value("ID").(string)

	userData, err := m.UserRepo.UserInfo(ID)
	if err != nil {
		errorMap.Message["user"] = locales.MustLocalize(&i18n.LocalizeConfig{MessageID: "UserInfoError"})
	}

	_, output := helpers.CustomResponse(userData, errorMap)
	helpers.ResGenerator(w, http.StatusOK, output)
}
