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
