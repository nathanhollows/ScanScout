package models

import (
	"context"
	"errors"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/Rapua/internal/sessions"
	"github.com/nathanhollows/Rapua/pkg/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	baseModel

	UserID            string    `bun:",unique,pk,type:varchar(36)" json:"user_id"`
	Name              string    `bun:",type:varchar(255)" json:"name"`
	Email             string    `bun:",unique,pk" json:"email"`
	Password          string    `bun:",type:varchar(255)" json:"password"`
	Instances         Instances `bun:"rel:has-many,join:user_id=user_id" json:"instances"`
	CurrentInstanceID string    `bun:",type:varchar(36)" json:"current_instance_id"`
	CurrentInstance   *Instance `bun:"rel:has-one,join:current_instance_id=id" json:"current_instance"`
}

type Users []*User

// Save the user to the database
func (u *User) Save(ctx context.Context) error {
	_, err := db.DB.NewInsert().Model(u).Exec(ctx)
	return err
}

// Update the user in the database
func (u *User) Update(ctx context.Context) error {
	_, err := db.DB.NewUpdate().Model(u).WherePK("user_id").Exec(ctx)
	return err
}

// AuthenticateUser checks the user's credentials and returns the user if they are valid
func AuthenticateUser(ctx context.Context, email, password string) (*User, error) {
	// Find the user by email
	user, err := FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Check the password
	if !user.checkPassword(password) {
		return nil, errors.New("Invalid password")
	} else {
		return user, nil
	}
}

// FindUserByEmail finds a user by their email address
func FindUserByEmail(ctx context.Context, email string) (*User, error) {
	// Find the user by email
	user := &User{}
	err := db.DB.NewSelect().
		Model(user).
		Where("email = ?", email).
		Relation("CurrentInstance").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindUserByID finds a user by their user id
func FindUserByID(ctx context.Context, userID string) (*User, error) {
	// Find the user by user id
	user := &User{}
	err := db.DB.NewSelect().
		Model(user).
		Where("User.user_id = ?", userID).
		Relation("CurrentInstance").
		Relation("Instances").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CheckPassword checks if the given password is correct
func (u *User) checkPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		log.Error("Error comparing password: ", err)
		return false
	}
	return true
}

// hashAndSalt hashes and salts the given password
func hashAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Error hashing password: ", err)
	}

	return string(hash)
}

// FindUserBySession finds the user by the session
func FindUserBySession(r *http.Request) (*User, error) {
	// Get the session
	session, err := sessions.Get(r, "admin")
	if err != nil {
		return nil, err
	}

	// Get the user id from the session
	userID, ok := session.Values["user_id"].(string)
	if !ok {
		return nil, errors.New("User not found")
	}

	// Find the user by the user id
	user, err := FindUserByID(r.Context(), userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
