package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/models"
	"github.com/nathanhollows/Rapua/repositories"
	"github.com/nathanhollows/Rapua/security"
)

// ErrPasswordsDoNotMatch is returned when the passwords do not match
var (
	ErrPasswordsDoNotMatch = errors.New("passwords do not match")
)

type UserService interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, user *models.User, passwordConfirm string) error

	// GetUserByEmail retrieves a user by their email address
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	// UpdateUser updates a user
	UpdateUser(ctx context.Context, user *models.User) error

	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, userID string) error
}

type userService struct {
	transactor         db.Transactor
	instanceRepository repositories.InstanceRepository
	userRepository     repositories.UserRepository
}

func NewUserService(transactor db.Transactor, userRepository repositories.UserRepository, instanceRepository repositories.InstanceRepository) UserService {
	return &userService{
		transactor:         transactor,
		instanceRepository: instanceRepository,
		userRepository:     userRepository,
	}
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

	return s.userRepository.Create(ctx, user)
}

// UpdateUser updates a user in the database
func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.userRepository.Update(ctx, user)
}

// GetUserByEmail retrieves a user by their email address
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepository.GetByEmail(ctx, email)
}

// DeleteUser deletes a user from the database
func (s *userService) DeleteUser(ctx context.Context, userID string) error {
	tx, err := s.transactor.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	// Ensure rollback on failure
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	err = s.userRepository.Delete(ctx, tx, userID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
	}

	err = s.instanceRepository.DeleteByUser(ctx, tx, userID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
	}

	return tx.Commit()
}
