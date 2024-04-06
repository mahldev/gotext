package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	CreatedAt *time.Time `json:"createAt"`
	UdpateAt  *time.Time `json:"updateAt"`
	Name      string     `gorm:"unique" json:"username"`
	Password  string     `json:"password"`
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
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(hashed) == u.Password
}
