package database

import "time"

type AuditLog struct {
	ID        int `gorm:"primaryKey"`
	UserID    string
	MessageID string
	Action    string
	Data      string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
