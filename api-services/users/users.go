package users

import (
	"github.com/authentication-app-server/api-services/models"
	"github.com/authentication-app-server/api-services/pkg"
	"github.com/authentication-app-server/helpers"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func EvaluateEditUserCredentials(user models.User, locales *i18n.Localizer) helpers.ErrorsMap {
	errResponse := helpers.ErrorsMap{
		Success: false,
		Message: map[string]string{},
	}

	if user.Email != "" {
		if isValidEmail := pkg.ValidEmail(user.Email); !isValidEmail {
			errResponse.Message["email"] = locales.MustLocalize(&i18n.LocalizeConfig{MessageID: "InvalidEmail"})
		}
	}

	if user.Password != "" {
		if isValidPassword := pkg.ValidPassword(user.Password); !isValidPassword {
			errResponse.Message["password"] =
				locales.MustLocalize(&i18n.LocalizeConfig{MessageID: "InvalidPassword"})
		}
	}

	if user.Telephone != "" {
		if isValidPhoneNumber := pkg.ValidTelephone(user.Telephone); !isValidPhoneNumber {
			errResponse.Message["telephone"] = locales.MustLocalize(&i18n.LocalizeConfig{MessageID: "InvalidTelephone"})
		}
	}

	return errResponse
}

func EvaluateFile(contentType string, fileSize int64, errorMap helpers.ErrorsMap, locales *i18n.Localizer) helpers.ErrorsMap {

	var (
		fileMaxSize      int64 = 3145728
		allowedFileTypes       = []string{"image/jpg", "image/jpeg", "image/png"}
		isAValidTypeFile bool
	)

	if fileSize > fileMaxSize {
		errorMap.Message["file_size"] = locales.MustLocalize(&i18n.LocalizeConfig{
			MessageID:    "ImageFileTooBig",
			TemplateData: map[string]int64{"Size": fileMaxSize}})
	}

	for _, allowedType := range allowedFileTypes {
		if isAValidTypeFile = allowedType == contentType; isAValidTypeFile {
			break
		}
	}

	if !isAValidTypeFile {
		errorMap.Message["file_type"] = locales.MustLocalize(&i18n.LocalizeConfig{MessageID: "WrongImageType"})
	}

	return errorMap
}
