package models

import (
	"time"

	"github.com/google/uuid"
)

type TestStats struct {
	CreatedAt time.Time `json:"created_at"`
	User      *User     `gorm:"foreignKey:UserID" json:"user"`
	WPM       float64   `json:"wpm"`
	Accuracy  float64   `json:"accuracy"`
	TimeTaken float64   `json:"timeTaken"`
	Errors    int       `json:"errors"`
	Words     int       `json:"words"`
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uuid.UUID `gorm:"type:char(36)" json:"user_id"`
}
