package repositories

import (
	"context"

	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/pkg/db"
	"github.com/uptrace/bun"
)

func CreateUser(ctx context.Context, user *models.User) error {
	_, err := db.DB.NewInsert().Model(user).Exec(ctx)
	return err
}

// GetUserByEmail retrieves a user by their email address
func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
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
func FindUserByID(ctx context.Context, userID string) (*models.User, error) {
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
		return nil, err
	}
	return &user, nil
}
