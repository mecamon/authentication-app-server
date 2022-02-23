package users

import (
	"fmt"
	"github.com/authentication-app-server/api-services/models"
	"github.com/authentication-app-server/api-services/pkg"
	"github.com/authentication-app-server/helpers"
)

func EvaluateEditUserCredentials(user models.User) helpers.ErrorsMap {
	errResponse := helpers.ErrorsMap{
		Success: false,
		Message: map[string]string{},
	}

	if user.Email != "" {
		if isValidEmail := pkg.ValidEmail(user.Email); !isValidEmail {
			errResponse.Message["email"] = "Invalid email address"
		}
	}

	if user.Password != "" {
		if isValidPassword := pkg.ValidPassword(user.Password); !isValidPassword {
			errResponse.Message["password"] =
				"Password must least contain one of these a-z, A-Z, 0-9 and be longer than 8 characters"
		}
	}

	if user.Telephone != "" {
		if isValidPhoneNumber := pkg.ValidTelephone(user.Telephone); !isValidPhoneNumber {
			errResponse.Message["telephone"] = "Telephone must be longer than 10 and contain only numbers"
		}
	}

	return errResponse
}

func EvaluateFile(contentType string, fileSize int64, errorMap helpers.ErrorsMap) helpers.ErrorsMap {

	var (
		fileMaxSize      int64 = 3145728
		allowedFileTypes       = []string{"image/jpg", "image/jpeg", "image/png"}
		isAValidTypeFile bool
	)

	if fileSize > fileMaxSize {
		errorMap.Message["file_size"] = fmt.Sprintf("File size can't be higher than %d bytes", fileMaxSize)
	}

	for _, allowedType := range allowedFileTypes {
		if isAValidTypeFile = allowedType == contentType; isAValidTypeFile {
			break
		}
	}

	if !isAValidTypeFile {
		errorMap.Message["file type"] = "Wrong file type. Only jpg, jpeg and png accepted"
	}

	return errorMap
}
