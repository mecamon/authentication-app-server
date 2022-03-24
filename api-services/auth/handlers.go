package auth

import (
	"encoding/json"
	"fmt"
	"github.com/authentication-app-server/api-services/models"
	"github.com/authentication-app-server/api-services/repository"
	"github.com/authentication-app-server/helpers"
	i18napp "github.com/authentication-app-server/i18n-app"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strings"
)

var (
	handlers     *Handlers
	ClientID     string
	ClientSecret string
)

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

	lang := r.Header.Get("Accept-Language")
	locales := i18napp.GetLocales(lang)

	err := json.NewDecoder(r.Body).Decode(&auth)

	result, err := m.AuthRepo.Login(auth.Email, auth.Password)
	hasAValidPass := helpers.CheckPassword(auth.Password, result.Password)

	if err != nil || !hasAValidPass {
		incorrectCredentials := locales.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "IncorrectCredentials",
		})
		errorMap = helpers.ErrorsMap{
			Success: false,
			Message: map[string]string{
				"credentials": incorrectCredentials,
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

	lang := r.Header.Get("Accept-Language")
	locales := i18napp.GetLocales(lang)

	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		log.Println(err)
	}

	errorMap := auth.EvaluateNewUserCredentials(locales)

	if len(errorMap.Message) > 0 {
		_, output := helpers.CustomResponse(nil, errorMap)
		helpers.ResGenerator(w, http.StatusBadRequest, output)
		return
	}

	hashed, _ := helpers.HashPassword(auth.Password)
	insertedID, err := m.AuthRepo.Register(auth.Email, string(hashed))
	if err != nil {
		if strings.Contains(err.Error(), "E11000") {
			errorMap.Message["email"] = locales.MustLocalize(&i18n.LocalizeConfig{MessageID: "EmailAddressTaken"})
			_, output := helpers.CustomResponse(nil, errorMap)
			helpers.ResGenerator(w, http.StatusConflict, output)
			return
		}
		errorMap.Message["register"] = locales.MustLocalize(&i18n.LocalizeConfig{MessageID: "GenericRegisterError"})
		_, output := helpers.CustomResponse(nil, errorMap)
		helpers.ResGenerator(w, http.StatusBadRequest, output)
		return
	}

	signedToken, err := helpers.GenerateToken(insertedID, auth.Email)

	resultMap := map[string]string{
		"token":       signedToken,
		"inserted_id": insertedID,
	}

	_, output := helpers.CustomResponse(resultMap, errorMap)
	helpers.ResGenerator(w, http.StatusCreated, output)
}

func (m *Handlers) accessRequestForGithub(w http.ResponseWriter, r *http.Request) {

	scope := "read:user user:email"
	state := helpers.RandomString(20)

	resultMap := map[string]string{
		"requestURL": fmt.Sprintf(
			"https://github.com/login/oauth/authorize?state=%s&client_id=%s&scope=%s",
			state,
			ClientID,
			scope,
		),
	}

	output, _ := json.Marshal(resultMap)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

func (m *Handlers) loginWithGithub(w http.ResponseWriter, r *http.Request) {
	var data = struct {
		Code  string
		State string
	}{}

	var errorMap = helpers.ErrorsMap{
		Success: false,
		Message: map[string]string{},
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
	}

	requestURL := fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&state=%s&client_secret=%s&code=%s",
		ClientID,
		data.State,
		ClientSecret,
		data.Code,
	)

	req, _ := http.NewRequest(http.MethodPost, requestURL, nil)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errorMap.Message["request"] = err.Error()
		_, output := helpers.CustomResponse(nil, errorMap)
		helpers.ResGenerator(w, http.StatusBadRequest, output)
		return
	}

	var githubAuthInfo = struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&githubAuthInfo)
	if err != nil {
		log.Println(err)
	}

	bearerToken := fmt.Sprintf("Bearer %s", githubAuthInfo.AccessToken)

	requestUserDataURL := "https://api.github.com/user"
	req, err = http.NewRequest(http.MethodGet, requestUserDataURL, nil)

	req.Header.Set("Authorization", bearerToken)
	req.Header.Set("Accept", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		errorMap.Message["request"] = err.Error()
		_, output := helpers.CustomResponse(nil, errorMap)
		helpers.ResGenerator(w, http.StatusBadRequest, output)
		return
	}

	githubUser := models.GithubUser{}

	err = json.NewDecoder(resp.Body).Decode(&githubUser)
	if err != nil {
		log.Println(err)
	}

	lang := r.Header.Get("Accept-Language")
	locales := i18napp.GetLocales(lang)

	result, err := m.AuthRepo.LoginWithGithub(githubUser)
	if err != nil {
		log.Println(err)
		errorMap.Message["login"] = locales.MustLocalize(&i18n.LocalizeConfig{MessageID: "GenericLoginError"})
		_, output := helpers.CustomResponse(nil, errorMap)
		helpers.ResGenerator(w, http.StatusBadRequest, output)
		return
	}

	signedToken, err := helpers.GenerateToken(primitive.ObjectID.Hex(result.ID), result.Email)
	tokenMap := map[string]string{"token": signedToken}

	_, output := helpers.CustomResponse(tokenMap, errorMap)
	helpers.ResGenerator(w, http.StatusOK, output)
}
