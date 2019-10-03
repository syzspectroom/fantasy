package model

import "golang.org/x/crypto/bcrypt"

// User db struct
type User struct {
	ID             string `json:"_key,omitempty"`
	Email          string `json:"email" validate:"required,email"`
	HashedPassword string `json:"hashed_password" validate:"required"`
}

// HashPassword hash plaintext password with bcrypt
func (u *User) HashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	u.HashedPassword = string(hash)
	return nil
}
