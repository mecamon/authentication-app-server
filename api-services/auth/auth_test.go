package auth

import (
	i18n_app "github.com/authentication-app-server/i18n-app"
	"testing"
)

func TestAuth_EvaluateNewUserCredentials(t *testing.T) {

	locales := i18n_app.GetLocales("en-US")

	t.Log("invalid user credentials")
	{
		auth := Auth{
			Email:    "invalidmail.co",
			Password: "123145",
		}
		errResponse := auth.EvaluateNewUserCredentials(locales)

		if len(errResponse.Message) != 2 {
			t.Error("has not the expected quantity of errors")
		}
	}
}
