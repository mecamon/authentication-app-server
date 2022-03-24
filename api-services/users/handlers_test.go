package users

import (
	"bytes"
	"github.com/authentication-app-server/api-services/repository"
	"github.com/authentication-app-server/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHandlers(t *testing.T) {
	var i interface{}

	i = NewHandlers(mongoTestClient, mongoTestDBName)

	if _, ok := i.(*Handlers); !ok {
		t.Errorf("Expected a auth handler type but got another")
	}
}

var userTestStruct = []struct {
	testName           string
	email              string
	newEmail           string
	password           string
	name               string
	telephone          string
	bio                string
	expectedStatusCode int
}{
	{"valid-changes", "usertest1@mail.com", "", "Password12345", "name test user 1", "1234567890", "test bio 1", http.StatusOK},
	{"email-taken", "usertest2@mail.com", "usertest1@mail.com", "Password12345", "name test user 2", "1234567890", "test bio 2", http.StatusConflict},
}

func TestUpdateUserInfo(t *testing.T) {

	authRepo := repository.NewAuthRepo(mongoTestClient, mongoTestDBName)
	createTestUsers(authRepo)

	for _, tt := range userTestStruct {
		t.Log(tt.testName)
		{

			result, _ := authRepo.Login(tt.email, tt.password)
			tokenString, _ := helpers.GenerateToken(primitive.ObjectID.Hex(result.ID), result.Email)

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("email", tt.newEmail)
			writer.WriteField("password", tt.password)
			writer.WriteField("name", tt.name)
			writer.WriteField("telephone", tt.telephone)
			writer.WriteField("bio", tt.bio)
			writer.Close()

			req := httptest.NewRequest(http.MethodPut, "/api/users/update", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.Header.Set("Authorization", tokenString)

			rr := httptest.NewRecorder()
			testMainRouter.ServeHTTP(rr, req)
			if statusCode := rr.Code; statusCode != tt.expectedStatusCode {
				t.Errorf("Expected statusCode: %d but got %d", tt.expectedStatusCode, rr.Code)
			}
		}
	}
}

func createTestUsers(authRepo repository.AuthRepo) {
	passwordHashed, _ := helpers.HashPassword(userTestStruct[0].password)
	authRepo.Register(userTestStruct[0].email, string(passwordHashed))
	authRepo.Register(userTestStruct[1].email, string(passwordHashed))
}

func TestUserInfo(t *testing.T) {

	t.Log("get-user-info")
	{
		authRepo := repository.NewAuthRepo(mongoTestClient, mongoTestDBName)

		result, _ := authRepo.Login("userinfotest@mail.com", "Password1234")
		tokenString, _ := helpers.GenerateToken(primitive.ObjectID.Hex(result.ID), result.Email)

		req := httptest.NewRequest(http.MethodGet, "/api/users/info", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", tokenString)

		rr := httptest.NewRecorder()
		testMainRouter.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected statusCode %d but got %d", http.StatusOK, rr.Code)
		}
	}
}
