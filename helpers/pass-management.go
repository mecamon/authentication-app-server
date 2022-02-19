package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

var saltLevel = 14

func HashPassword(password string) ([]byte, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), saltLevel)
	return hashedPass, err
}

func CheckPassword(insertedPass, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(insertedPass))
	return err == nil
}
