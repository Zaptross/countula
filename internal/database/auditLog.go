package database

import "time"

type AuditLog struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	UserID    string
	Username  string // just for human readability
	MessageID string
	Action    string
	Data      string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
