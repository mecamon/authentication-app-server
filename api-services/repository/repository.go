package repository

import "github.com/authentication-app-server/api-services/models"

type AuthRepo interface {
	Register(email string, password string) (string, error)
	Login(email string, password string) (models.User, error)
	LoginWithGithub(githubUser models.GithubUser) (models.User, error)
}

type UserRepo interface {
	UpdateUser(ID string, validUser models.User) (int64, error)
	UserInfo(ID string) (interface{}, error)
}
