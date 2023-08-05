package database

import "time"

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
