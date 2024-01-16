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
	ChannelID string
	Rules     int
	Guess     int
	Correct   bool
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func GetCurrentTurn(db *gorm.DB, channelID string) Turn {
	var currentTurn Turn
	db.Order("game desc, turn desc").Where("channel_id = ?", channelID).First(&currentTurn)
	return currentTurn
}

func GetHighScoreTurn(db *gorm.DB, channelID string) Turn {
	var highScoreTurn Turn
	db.Order("turn desc").Where("channel_id = ? AND correct = ?", channelID, true).First(&highScoreTurn)
	return highScoreTurn
}

func GetNextGame(db *gorm.DB, channelID string) int {
	var lastTurn Turn
	db.Last(&lastTurn).Where("channel_id = ?", channelID)
	return lastTurn.Game + 1
}

func CreateTurnFromContext(db *gorm.DB, s *discordgo.Session, m *discordgo.MessageCreate, lastTurn Turn, guess int, correct bool) Turn {
	newTurn := Turn{
		Game:      lastTurn.Game,
		Turn:      lastTurn.Turn + 1,
		UserID:    m.Author.ID,
		MessageID: m.Message.ID,
		ChannelID: m.ChannelID,
		Rules:     lastTurn.Rules,
		Guess:     guess,
		Correct:   correct,
	}

	db.Create(&newTurn)

	return newTurn
}
