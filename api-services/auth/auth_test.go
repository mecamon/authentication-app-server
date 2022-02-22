package auth

import "testing"

func TestAuth_EvaluateNewUserCredentials(t *testing.T) {
	t.Log("invalid user credentials")
	{
		auth := Auth{
			Email:    "invalidmail.co",
			Password: "123145",
		}
		errResponse := auth.EvaluateNewUserCredentials()

		if len(errResponse.Message) != 2 {
			t.Error("has not the expected quantity of errors")
		}
	}
}
