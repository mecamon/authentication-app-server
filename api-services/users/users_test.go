package users

import (
	"github.com/authentication-app-server/api-services/models"
	"github.com/authentication-app-server/helpers"
	i18n_app "github.com/authentication-app-server/i18n-app"
	"testing"
)

func TestEvaluateEditUserCredentials(t *testing.T) {
	user := models.User{
		Name:      "Random name",
		Bio:       "This is just a sample bio.",
		Email:     "wrongmail.com",
		Password:  "wrongpass",
		Telephone: "wrongphone",
	}

	locales := i18n_app.GetLocales("en-US")

	output := EvaluateEditUserCredentials(user, locales)

	if len(output.Message) != 3 {
		t.Errorf("Expected 3 errors and got %d", len(output.Message))
	}
}

var fileTests = []struct {
	testName               string
	fileSize               int64
	contentType            string
	expectedErrorMsgLength int
}{
	{"Wrong size and content type", 4000345, "text/pdf", 2},
	{"Proper size and content type", 2000345, "image/jpeg", 0},
}

func TestEvaluateFile(t *testing.T) {

	locales := i18n_app.GetLocales("en-US")

	for _, tt := range fileTests {
		t.Log(tt.testName)
		{
			emptyErrMap := helpers.ErrorsMap{
				Success: false,
				Message: map[string]string{},
			}
			errMap := EvaluateFile(tt.contentType, tt.fileSize, emptyErrMap, locales)

			if len(errMap.Message) != tt.expectedErrorMsgLength {
				t.Errorf("Expected %d error messages but got %d", tt.expectedErrorMsgLength, len(errMap.Message))
			}
		}
	}
}
