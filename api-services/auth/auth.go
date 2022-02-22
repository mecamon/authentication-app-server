package auth

import (
	"github.com/authentication-app-server/api-services/pkg"
	"github.com/authentication-app-server/helpers"
)

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m *Auth) EvaluateNewUserCredentials() helpers.ErrorsMap {
	errResponse := helpers.ErrorsMap{
		Success: false,
		Message: map[string]string{},
	}

	isValidEmail := pkg.ValidEmail(m.Email)
	if !isValidEmail {
		errResponse.Message["email"] = "Invalid email address"
	}

	isValidPassword := pkg.ValidPassword(m.Password)
	if !isValidPassword {
		errResponse.Message["password"] = "Password must least contain one of these a-z, A-Z, 0-9 and be longer than 8 characters"
	}

	return errResponse
}
