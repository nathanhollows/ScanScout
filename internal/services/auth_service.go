package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
	"github.com/nathanhollows/Rapua/internal/sessions"
	"github.com/nathanhollows/Rapua/pkg/security"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// AuthenticateUser authenticates the user with the given email and password.
func (s *AuthService) AuthenticateUser(ctx context.Context, email, password string) (*models.User, error) {
	user, err := repositories.GetUserByEmail(ctx, email)
	if err != nil {
		// Wrap the error
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}

	if !security.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

// GetAuthenticatedUser retrieves the authenticated user from the session.
func (s *AuthService) GetAuthenticatedUser(r *http.Request) (*models.User, error) {
	session, err := sessions.Get(r, "admin")
	if err != nil {
		return nil, err
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		return nil, errors.New("user not authenticated")
	}

	user, err := repositories.FindUserByID(r.Context(), userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
