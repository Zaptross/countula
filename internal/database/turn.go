package database

import (
	"time"

	"github.com/bwmarrin/discordgo"
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
	db.Order("game desc, turn desc").First(&currentTurn)
	return currentTurn
}

func GetNextGame(db *gorm.DB) int {
	var lastTurn Turn
	db.Last(&lastTurn)
	return lastTurn.Game + 1
}

func CreateTurnFromContext(db *gorm.DB, s *discordgo.Session, m *discordgo.MessageCreate, lastTurn Turn, guess int, correct bool) Turn {
	newTurn := Turn{
		Game:      lastTurn.Game,
		Turn:      lastTurn.Turn + 1,
		UserID:    m.Author.ID,
		MessageID: m.Message.ID,
		Rules:     lastTurn.Rules,
		Guess:     guess,
		Correct:   correct,
	}

	db.Create(&newTurn)

	return newTurn
}
