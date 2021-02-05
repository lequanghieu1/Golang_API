package models

import (
	"html"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Manager struct {
	Username  string    `json:"user_name" bson:"user_name"`
	Password  string    `json:"password" bson:"password"`
	Fullname  string    `json:"full_name" bson:"full_name"`
	IsLogin   bool      `json:"is_login" bson:"is_login"`
	CreatedAt time.Time `json:"created_at, omitempty"`
	UpdatedAt time.Time `json:"updated_at, omitempty"`
}

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func Santize(data string) string {
	data = html.EscapeString(strings.TrimSpace(data))
	return data
}
