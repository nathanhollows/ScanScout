package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/pkg/db"
	"github.com/uptrace/bun"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	// CreateUser creates a new user in the database
	CreateUser(ctx context.Context, user *models.User) error
	// UpdateUser updates a user in the database
	UpdateUser(ctx context.Context, user *models.User) error
	// GetUserByEmail retrieves a user by their email address
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	// FindUserByID fetches a user by their ID from the database.
	FindUserByID(ctx context.Context, userID string) (*models.User, error)
	// GetUserByEmailAndProvider retrieves a user by their email address and provider
	GetUserByEmailAndProvider(ctx context.Context, email, provider string) (*models.User, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

// Update the user in the database
func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := db.DB.NewUpdate().
		Model(user).
		WherePK().
		Exec(ctx)
	return err
}

// CreateUser creates a new user in the database
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	if user.ID == "" {
		uuid := uuid.New()
		user.ID = uuid.String()
	}

	_, err := db.DB.NewInsert().Model(user).Exec(ctx)
	return err
}

// GetUserByEmail retrieves a user by their email address
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := db.DB.NewSelect().
		Model(user).
		Where("email = ?", email).
		Relation("CurrentInstance").
		Relation("Instances").
		Scan(ctx)
	return user, err
}

// FindUserByID fetches a user by their ID from the database.
func (r *userRepository) FindUserByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	err := db.DB.NewSelect().
		Model(&user).
		Where("user.id = ?", userID).
		Relation("CurrentInstance").
		Relation("CurrentInstance.Settings").
		Relation("CurrentInstance.Teams").
		// Locations ordered by Order
		Relation("CurrentInstance.Locations", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("order ASC")
		}).
		Relation("CurrentInstance.Locations.Marker").
		Relation("Instances").
		Scan(ctx)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

// GetUserByEmailAndProvider retrieves a user by their email address and provider
func (r *userRepository) GetUserByEmailAndProvider(ctx context.Context, email, provider string) (*models.User, error) {
	user := &models.User{}
	err := db.DB.NewSelect().
		Model(user).
		Where("email = ?", email).
		Where("provider = ? OR provider = ''", provider).
		Relation("CurrentInstance").
		Relation("Instances").
		Scan(ctx)
	return user, err
}
