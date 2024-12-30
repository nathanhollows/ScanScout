package services_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/nathanhollows/Rapua/internal/services"
	"github.com/nathanhollows/Rapua/models"
	"github.com/nathanhollows/Rapua/repositories"
	"github.com/stretchr/testify/assert"
)

func setupUserService(t *testing.T) (services.UserService, func()) {
	t.Helper()
	db, cleanup := setupDB(t)

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	return userService, cleanup
}

func setupNotificationRepo(t *testing.T) (repositories.NotificationRepository, func()) {
	t.Helper()
	db, cleanup := setupDB(t)

	notificationRepo := repositories.NewNotificationRepository(db)
	return notificationRepo, cleanup
}

func TestCreateUser(t *testing.T) {
	service, cleanup := setupUserService(t)
	defer cleanup()

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 12)

	user := &models.User{
		Email:    email,
		Password: password,
	}
	err := service.CreateUser(context.Background(), user, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.NotEqual(t, password, user.Password) // Ensure password is transformed/hashed
}

func TestCreateUser_PasswordsDoNotMatch(t *testing.T) {
	service, cleanup := setupUserService(t)
	defer cleanup()

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 12)
	user := &models.User{
		Email:    email,
		Password: password,
	}
	err := service.CreateUser(context.Background(), user, "differentPassword")

	assert.Error(t, err)
	assert.Equal(t, services.ErrPasswordsDoNotMatch, err)
}

func TestGetUserByEmail(t *testing.T) {
	service, cleanup := setupUserService(t)
	defer cleanup()

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 12)

	user := &models.User{
		Email:    email,
		Password: password,
	}
	err := service.CreateUser(context.Background(), user, password)
	assert.NoError(t, err)

	retrievedUser, err := service.GetUserByEmail(context.Background(), email)

	assert.NoError(t, err)
	assert.Equal(t, user.Email, retrievedUser.Email)
	assert.Equal(t, user.ID, retrievedUser.ID)
}

func TestUpdateUser(t *testing.T) {
	service, cleanup := setupUserService(t)
	defer cleanup()

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 12)

	user := &models.User{
		Email:    email,
		Password: password,
	}
	err := service.CreateUser(context.Background(), user, password)
	assert.NoError(t, err)

	newName := gofakeit.Name()
	user.Name = newName
	err = service.UpdateUser(context.Background(), user)
	assert.NoError(t, err)

	retrievedUser, err := service.GetUserByEmail(context.Background(), email)
	assert.NotEmpty(t, retrievedUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, newName, retrievedUser.Name)
}

func TestDeleteUser(t *testing.T) {
	service, cleanup := setupUserService(t)
	defer cleanup()

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 12)

	user := &models.User{
		Email:    email,
		Password: password,
	}
	err := service.CreateUser(context.Background(), user, password)
	assert.NoError(t, err)

	err = service.DeleteUser(context.Background(), user)
	assert.NoError(t, err)

	_, err = service.GetUserByEmail(context.Background(), email)
	assert.Error(t, err)
}
