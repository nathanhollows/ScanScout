package services

import (
	"context"
	"errors"

	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
	"github.com/nathanhollows/Rapua/pkg/security"
)

// AuthenticateUser authenticates the user with the given email and password.
func AuthenticateUser(ctx context.Context, email, password string) (*models.User, error) {
	user, err := repositories.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !security.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
