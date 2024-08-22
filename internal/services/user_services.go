package services

import (
	"github.com/nathanhollows/Rapua/internal/repositories"
)

type UserServices struct {
	AuthService AuthService
	UserService UserService
}

func NewUserServices(repo repositories.UserRepository) *UserServices {
	return &UserServices{
		AuthService: NewAuthService(repo),
		UserService: NewUserService(repo),
	}
}
