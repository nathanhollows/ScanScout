package models

import (
	"context"
	"errors"

	"github.com/charmbracelet/log"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	baseModel

	UserID    string    `bun:",pk,type:varchar(36)" json:"user_id"`
	Email     string    `bun:",unique,pk" json:"email"`
	Password  string    `bun:",type:varchar(255)" json:"password"`
	Instances Instances `bun:"rel:has-many,join:user_id=user_id" json:"instances"`
}

type Users []*User

// Save the user to the database
func (u *User) Save() error {
	ctx := context.Background()
	_, err := db.NewInsert().Model(u).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// AuthenticateUser checks the user's credentials and returns the user if they are valid
func AuthenticateUser(email, password string) (*User, error) {
	// Find the user by email
	user, err := FindUserByEmail(email)
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
func FindUserByEmail(email string) (*User, error) {
	ctx := context.Background()
	// Find the user by email
	user := &User{}
	err := db.NewSelect().Model(user).Where("email = ?", email).Scan(ctx)
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
