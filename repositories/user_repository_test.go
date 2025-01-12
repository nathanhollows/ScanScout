package repositories_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/models"
	"github.com/nathanhollows/Rapua/repositories"
	"github.com/stretchr/testify/assert"
)

func setupUserRepo(t *testing.T) (repositories.UserRepository, db.Transactor, func()) {
	t.Helper()
	dbc, cleanup := setupDB(t)

	transactor := db.NewTransactor(dbc)

	userRepository := repositories.NewUserRepository(dbc)
	return userRepository, transactor, cleanup
}

func TestUserRepository_Create(t *testing.T) {
	repo, _, cleanup := setupUserRepo(t)
	defer cleanup()
	ctx := context.Background()

	user := &models.User{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, true, 12),
		Provider: "local",
	}

	err := repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	repo, _, cleanup := setupUserRepo(t)
	defer cleanup()
	ctx := context.Background()

	// Seed user
	user := &models.User{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, true, 12),
		Provider: "local",
	}
	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	// Test GetUserByEmail
	fetchedUser, err := repo.FindUserByEmail(ctx, user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, fetchedUser.Email)
	assert.Equal(t, user.Name, fetchedUser.Name)
}

func TestUserRepository_Update(t *testing.T) {
	repo, _, cleanup := setupUserRepo(t)
	defer cleanup()
	ctx := context.Background()

	// Seed user
	user := &models.User{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, true, 12),
		Provider: "local",
	}
	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	// Update user
	newName := gofakeit.Name()
	user.Name = newName
	err = repo.Update(ctx, user)
	assert.NoError(t, err)

	// Verify update
	fetchedUser, err := repo.FindUserByEmail(ctx, user.Email)
	assert.NoError(t, err)
	assert.Equal(t, newName, fetchedUser.Name)
}

func TestUserRepository_FindUserByID(t *testing.T) {
	repo, _, cleanup := setupUserRepo(t)
	defer cleanup()
	ctx := context.Background()

	// Seed user
	user := &models.User{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, true, 12),
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
	repo, _, cleanup := setupUserRepo(t)
	defer cleanup()
	ctx := context.Background()

	// Seed user
	user := &models.User{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, true, 12),
		Provider: "local",
	}
	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	// Test GetUserByEmailAndProvider
	fetchedUser, err := repo.FindUserByEmailAndProvider(ctx, user.Email, "local")
	assert.NoError(t, err)
	assert.Equal(t, user.Email, fetchedUser.Email)
	assert.Equal(t, user.Provider, fetchedUser.Provider)
}

func TestUserRepository_Delete(t *testing.T) {
	repo, transactor, cleanup := setupUserRepo(t)
	defer cleanup()
	ctx := context.Background()

	// Seed user
	user := &models.User{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, true, 12),
		Provider: "local",
	}
	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	// Delete user
	tx, err := transactor.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		t.Fatal(err)
	}
	err = repo.Delete(ctx, tx, user.ID)
	assert.NoError(t, err)
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	// Verify user is deleted
	fetchedUser, err := repo.FindUserByEmail(ctx, user.Email)
	assert.Error(t, err)
	assert.Empty(t, fetchedUser)
}
