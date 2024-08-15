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
var ErrPasswordsDoNotMatch = errors.New("passwords do not match")

func CreateUser(ctx context.Context, user *models.User, passwordConfirm string) error {
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

	return repositories.CreateUser(ctx, user)
}
