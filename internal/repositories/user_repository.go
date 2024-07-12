package repositories

import (
	"context"

	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/pkg/db"
)

func CreateUser(ctx context.Context, user *models.User) error {
	_, err := db.DB.NewInsert().Model(user).Exec(ctx)
	return err
}

// GetUserByEmail retrieves a user by their email address
func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := db.DB.NewSelect().Model(user).Where("email = ?", email).Scan(ctx)
	return user, err
}
