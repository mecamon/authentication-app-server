package auth

import (
	"github.com/authentication-app-server/helpers"
	"regexp"
)

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m *Auth) ValidEmail() bool {
	pattern := "^.{2,}@.{2,}\\..{2,}$"
	output, err := regexp.MatchString(pattern, m.Email)
	if err != nil {
		return false
	}

	return output
}

func (m *Auth) ValidPassword() bool {
	regx1 := regexp.MustCompile("(.*[a-z])")
	regx2 := regexp.MustCompile("(.*[A-Z])")
	regx3 := regexp.MustCompile("(.*[\\d])")
	regx4 := regexp.MustCompile("^.{8,}$")

	output1 := regx1.MatchString(m.Password)
	output2 := regx2.MatchString(m.Password)
	output3 := regx3.MatchString(m.Password)
	output4 := regx4.MatchString(m.Password)

	return output1 && output2 && output3 && output4
}

func (m *Auth) EvaluateNewUserCredentials() helpers.ErrorsMap {
	errResponse := helpers.ErrorsMap{
		Success: false,
		Message: map[string]string{},
	}

	isValidEmail := m.ValidEmail()
	if !isValidEmail {
		errResponse.Message["email"] = "Invalid email address"
	}

	isValidPassword := m.ValidPassword()
	if !isValidPassword {
		errResponse.Message["password"] = "Password must least contain one of these a-z, A-Z, 0-9 and be longer than 8 characters"
	}

	return errResponse
}
