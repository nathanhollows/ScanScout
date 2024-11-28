package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/models"
	"github.com/uptrace/bun"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	// Create creates a new user in the database
	Create(ctx context.Context, user *models.User) error
	// Update updates a user in the database
	Update(ctx context.Context, user *models.User) error
	// FindUserByEmail retrieves a user by their email address
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	// FindUserByEmailToken retrieves a user by their email token
	FindUserByEmailToken(ctx context.Context, token string) (*models.User, error)
	// FindUserByID fetches a user by their ID from the database.
	FindUserByID(ctx context.Context, userID string) (*models.User, error)
	// FindUserByEmailAndProvider retrieves a user by their email address and provider
	FindUserByEmailAndProvider(ctx context.Context, email, provider string) (*models.User, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

// Update the user in the database
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	_, err := db.DB.NewUpdate().
		Model(user).
		WherePK().
		Exec(ctx)
	return err
}

// Create creates a new user in the database
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	if user.ID == "" {
		uuid := uuid.New()
		user.ID = uuid.String()
	}

	_, err := db.DB.NewInsert().Model(user).Exec(ctx)
	return err
}

// FindUserByEmail retrieves a user by their email address
func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := db.DB.NewSelect().
		Model(user).
		Where("email = ?", email).
		Relation("CurrentInstance").
		Relation("Instances").
		Scan(ctx)
	return user, err
}

// FindUserByEmailToken retrieves a user by their email token
func (r *userRepository) FindUserByEmailToken(ctx context.Context, token string) (*models.User, error) {
	user := &models.User{}
	err := db.DB.NewSelect().
		Model(user).
		Where("email_token = ?", token).
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

// FindUserByEmailAndProvider retrieves a user by their email address and provider
func (r *userRepository) FindUserByEmailAndProvider(ctx context.Context, email, provider string) (*models.User, error) {
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
