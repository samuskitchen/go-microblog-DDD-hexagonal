package domain

import (
	"errors"
	"github.com/badoux/checkmail"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User of the system.
type User struct {
	ID           uint      `json:"id,omitempty"`
	FirstName    string    `json:"first_name,omitempty"`
	LastName     string    `json:"last_name,omitempty"`
	Username     string    `json:"username,omitempty"`
	Email        string    `json:"email,omitempty"`
	Picture      string    `json:"picture,omitempty"`
	Password     string    `json:"password,omitempty"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

// HashPassword generates a hash of the password and places the result in PasswordHash.
func (u *User) HashPassword() error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(passwordHash)

	return nil
}

// PasswordMatch compares HashPassword with the password and returns true if they match.
func (u User) PasswordMatch(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	return err == nil
}

// Validate is the validation method for mandatory fields
func (u *User) Validate(action string) error {

	const msgErrorPass = "required password"
	const msgErrorEmail = "required Email"
	const msgErrorEmailRequired = "invalid email"

	switch strings.ToLower(action) {

	case "login":
		if u.Password == "" {
			return errors.New(msgErrorPass)
		}

		if u.Email == "" {
			return errors.New(msgErrorEmail)
		}

		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New(msgErrorEmailRequired)
		}

		return nil

	default:
		if u.Username == "" {
			return errors.New("required nickname")
		}

		if u.Password == "" {
			return errors.New(msgErrorPass)
		}

		if u.Email == "" {
			return errors.New(msgErrorEmail)
		}

		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New(msgErrorEmailRequired)
		}

		return nil
	}

}
