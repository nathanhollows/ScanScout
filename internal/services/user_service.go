package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
	"github.com/nathanhollows/Rapua/pkg/security"
)

// ErrPasswordsDoNotMatch is returned when the passwords do not match
var (
	ErrPasswordsDoNotMatch = errors.New("passwords do not match")
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User, passwordConfirm string) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{
		userRepository: repositories.NewUserRepository(),
	}
}

// GetUserByEmail retrieves a user by their email address
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepository.GetUserByEmail(ctx, email)
}

// CreateUser creates a new user in the database
func (s *userService) CreateUser(ctx context.Context, user *models.User, passwordConfirm string) error {
	// Confirm passwords match
	if user.Password != passwordConfirm {
		return ErrPasswordsDoNotMatch
	}

	// Hash the password
	hashedPassword, err := security.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Generate UUID for user
	user.ID = uuid.New().String()

	return s.userRepository.CreateUser(ctx, user)
}
