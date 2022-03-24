package auth

import (
	"github.com/authentication-app-server/api-services/pkg"
	"github.com/authentication-app-server/helpers"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m *Auth) EvaluateNewUserCredentials(locales *i18n.Localizer) helpers.ErrorsMap {
	errResponse := helpers.ErrorsMap{
		Success: false,
		Message: map[string]string{},
	}

	isValidEmail := pkg.ValidEmail(m.Email)
	if !isValidEmail {
		errResponse.Message["email"] = locales.MustLocalize(&i18n.LocalizeConfig{MessageID: "InvalidEmail"})
	}

	isValidPassword := pkg.ValidPassword(m.Password)
	if !isValidPassword {
		errResponse.Message["password"] = locales.MustLocalize(&i18n.LocalizeConfig{MessageID: "InvalidPassword"})
	}

	return errResponse
}
