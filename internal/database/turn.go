package database

import (
	"time"

	"gorm.io/gorm"
)

type Turn struct {
	Game      int `gorm:"primaryKey"`
	Turn      int `gorm:"primaryKey"`
	UserID    string
	MessageID string
	Rules     int
	Guess     int
	Correct   bool
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func GetCurrentTurn(db *gorm.DB) Turn {
	var currentTurn Turn
	db.Last(&currentTurn)
	return currentTurn
}
