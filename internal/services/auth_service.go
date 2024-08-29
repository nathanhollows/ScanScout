package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
	"github.com/nathanhollows/Rapua/internal/sessions"
	"github.com/nathanhollows/Rapua/pkg/security"
)

var (
	ErrUserNotAuthenticated = errors.New("user not authenticated")
	ErrSessionNotFound      = errors.New("session not found")
)

type AuthService interface {
	AuthenticateUser(ctx context.Context, email, password string) (*models.User, error)
	GetAuthenticatedUser(r *http.Request) (*models.User, error)
	OAuthLogin(ctx context.Context, provider string, user goth.User) (*models.User, error)
	CheckUserRegisteredWithOAuth(ctx context.Context, provider, userID string) (*models.User, error)
	CreateUserWithOAuth(ctx context.Context, user goth.User) (*models.User, error)
	CompleteUserAuth(w http.ResponseWriter, r *http.Request) (*models.User, error)
}

type authService struct {
	userRepository repositories.UserRepository
	emailService   EmailService
}

func NewAuthService(userRepository repositories.UserRepository) AuthService {
	return &authService{
		userRepository: userRepository,
		emailService:   NewEmailService(),
	}
}

// AuthenticateUser authenticates the user with the given email and password.
func (s *authService) AuthenticateUser(ctx context.Context, email, password string) (*models.User, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}

	if !security.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

// GetAuthenticatedUser retrieves the authenticated user from the session.
func (s *authService) GetAuthenticatedUser(r *http.Request) (*models.User, error) {
	session, err := sessions.Get(r, "admin")
	if err != nil {
		return nil, err
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		return nil, errors.New("user not authenticated")
	}

	user, err := s.userRepository.FindUserByID(r.Context(), userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// OAuthLogin handles User Login via OAuth
func (s *authService) OAuthLogin(ctx context.Context, provider string, oauthUser goth.User) (*models.User, error) {
	existingUser, err := s.userRepository.GetUserByEmail(ctx, oauthUser.Email)
	if err != nil {
		// User doesn't exist, create a new one
		newUser, err := s.CreateUserWithOAuth(ctx, oauthUser)
		if err != nil {
			return nil, fmt.Errorf("creating user with OAuth: %w", err)
		}
		return newUser, nil
	}

	return existingUser, nil
}

// CheckUserRegisteredWithOAuth looks for user already registered with OAuth
func (s *authService) CheckUserRegisteredWithOAuth(ctx context.Context, provider, email string) (*models.User, error) {
	user, err := s.userRepository.GetUserByEmailAndProvider(ctx, email, provider)
	if err != nil {
		return nil, fmt.Errorf("getting user by email and provider: %w", err)
	}

	return user, nil
}

// CreateUserWithOAuth creates a new user if logging in with OAuth for the first time
func (s *authService) CreateUserWithOAuth(ctx context.Context, user goth.User) (*models.User, error) {
	uuid := uuid.New()
	newUser := models.User{
		ID:       uuid.String(),
		Name:     user.Name,
		Email:    user.Email,
		Password: "",
		Provider: user.Provider,
	}

	err := s.userRepository.CreateUser(ctx, &newUser)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	return &newUser, nil
}

// CompleteUserAuth completes the user authentication process
func (s *authService) CompleteUserAuth(w http.ResponseWriter, r *http.Request) (*models.User, error) {
	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		return nil, fmt.Errorf("completing user auth: %w", err)
	}

	user, err := s.OAuthLogin(r.Context(), gothUser.Provider, gothUser)
	if err != nil {
		return nil, fmt.Errorf("OAuth login: %w", err)
	}

	return user, nil
}
