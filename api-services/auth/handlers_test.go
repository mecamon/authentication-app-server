package auth

import (
	"bytes"
	"encoding/json"
	"log"
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

var registerTests = []struct {
	testName           string
	email              string
	password           string
	expectedStatusCode int
}{
	{testName: "valid-credentials", email: "pipi@test.com", password: "Validpass1234", expectedStatusCode: http.StatusCreated},
	{testName: "email-in-user", email: "pipi@test.com", password: "Validpass1234", expectedStatusCode: http.StatusConflict},
	{testName: "invalid-credentials", email: "pipicte.com", password: "invalidpass", expectedStatusCode: http.StatusBadRequest},
}

func TestHandlers_Register(t *testing.T) {
	testServer := httptest.NewTLSServer(testMainRouter)
	defer testServer.Close()

	for _, tt := range registerTests {
		t.Log(tt.testName)
		{
			body := auth{Email: tt.email, Password: tt.password}

			jsonBody, err := json.Marshal(body)
			if err != nil {
				log.Println("Error marshaling json")
			}

			resp, err := testServer.Client().Post(
				testServer.URL+"/api/auth/register",
				"application/json",
				bytes.NewReader(jsonBody),
			)

			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf(
					"Wrong status code. Expected %d but got %d",
					tt.expectedStatusCode,
					resp.StatusCode,
				)
			}
		}
	}
}

const defaultLoginUser = "login@user.com"
const defaultLoginPassword = "Password12345"

var loginTests = []struct {
	testName           string
	email              string
	password           string
	expectedStatusCode int
}{
	{testName: "valid-credentials", email: defaultLoginUser, password: defaultLoginPassword, expectedStatusCode: http.StatusOK},
	{testName: "invalid-credentials", email: "user@notloged.com", password: "Password12345", expectedStatusCode: http.StatusBadRequest},
}

func TestHandlers_Login(t *testing.T) {
	testServer := httptest.NewTLSServer(testMainRouter)
	defer testServer.Close()

	body := auth{
		Email:    defaultLoginUser,
		Password: defaultLoginPassword,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Println("Error marshaling json")
	}

	testServer.Client().Post(
		testServer.URL+"/api/auth/register",
		"application/json",
		bytes.NewReader(jsonBody),
	)

	for _, tt := range loginTests {
		t.Log(tt.testName)
		{
			body := auth{
				Email:    tt.email,
				Password: tt.password,
			}

			jsonBody, err := json.Marshal(body)
			if err != nil {
				log.Println("Error marshaling json")
			}

			resp, err := testServer.Client().Post(
				testServer.URL+"/api/auth/login",
				"application/json",
				bytes.NewReader(jsonBody),
			)

			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("Wrong status code. Expected %d but got %d",
					tt.expectedStatusCode,
					resp.StatusCode,
				)
			}
		}
	}

}
