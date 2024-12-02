package repositories_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/internal/migrate"
	"github.com/nathanhollows/Rapua/models"
	"github.com/nathanhollows/Rapua/repositories"
	"github.com/stretchr/testify/assert"
)

func setupUserRepo(t *testing.T) (repositories.UserRepository, func()) {
	t.Helper()
	os.Setenv("DB_CONNECTION", "file::memory:?cache=shared")
	os.Setenv("DB_TYPE", "sqlite3")
	db := db.MustOpen()

	// Create tables
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	migrate.CreateTables(logger, db)

	userRepository := repositories.NewUserRepository(db)
	return userRepository, func() {
		db.Close()
	}
}

func TestUserRepository_Create(t *testing.T) {
	repo, cleanup := setupUserRepo(t)
	defer cleanup()
	ctx := context.Background()

	user := &models.User{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "hashed_password",
		Provider: "local",
	}

	err := repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	repo, cleanup := setupUserRepo(t)
	defer cleanup()
	ctx := context.Background()

	// Seed user
	user := &models.User{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "hashed_password",
		Provider: "local",
	}
	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	// Test GetUserByEmail
	fetchedUser, err := repo.FindUserByEmail(ctx, "john.doe@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user.Email, fetchedUser.Email)
	assert.Equal(t, user.Name, fetchedUser.Name)
}

func TestUserRepository_Update(t *testing.T) {
	repo, cleanup := setupUserRepo(t)
	defer cleanup()
	ctx := context.Background()

	// Seed user
	user := &models.User{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "hashed_password",
		Provider: "local",
	}
	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	// Update user
	user.Name = "John Updated"
	err = repo.Update(ctx, user)
	assert.NoError(t, err)

	// Verify update
	fetchedUser, err := repo.FindUserByEmail(ctx, "john.doe@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "John Updated", fetchedUser.Name)
}

func TestUserRepository_FindUserByID(t *testing.T) {
	repo, cleanup := setupUserRepo(t)
	defer cleanup()
	ctx := context.Background()

	// Seed user
	user := &models.User{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "hashed_password",
		Provider: "local",
	}
	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	// Test FindUserByID
	fetchedUser, err := repo.FindUserByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, fetchedUser.ID)
	assert.Equal(t, user.Email, fetchedUser.Email)
}

func TestUserRepository_GetUserByEmailAndProvider(t *testing.T) {
	repo, cleanup := setupUserRepo(t)
	defer cleanup()
	ctx := context.Background()

	// Seed user
	user := &models.User{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "hashed_password",
		Provider: "local",
	}
	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	// Test GetUserByEmailAndProvider
	fetchedUser, err := repo.FindUserByEmailAndProvider(ctx, "john.doe@example.com", "local")
	assert.NoError(t, err)
	assert.Equal(t, user.Email, fetchedUser.Email)
	assert.Equal(t, user.Provider, fetchedUser.Provider)
}
