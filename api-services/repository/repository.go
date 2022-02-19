package repository

import "github.com/authentication-app-server/api-services/models"

type AuthRepo interface {
	Register(email string, password string) (string, error)
	Login(email string, password string) (models.User, error)
}
