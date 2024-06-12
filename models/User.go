package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	CreatedAt *time.Time `json:"createAt"`
	UpdateAt  *time.Time `json:"updateAt"`
	Name      string     `gorm:"unique" json:"username"`
	Password  string     `json:"password"`
	Email     string     `json:"email"`
	Level     int64      `json:"level"`
	ID        uuid.UUID  `gorm:"primaryKey;type:char(36)" json:"id"`
}

func (u *User) IsValidPassword() bool {
	return u.Password == ""
}

func (u *User) IsValidUsername() bool {
	return u.Name == ""
}

func (u *User) HashPassword() {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	u.Password = string(hashed)
}

func (u *User) PasswordIsCorrect(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
