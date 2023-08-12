package game

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
)

func CreateNewGame(db *gorm.DB, s *discordgo.Session, channelID string) database.Turn {
	newGame := database.Turn{
		UserID:  s.State.User.ID,
		Game:    database.GetNextGame(db),
		Rules:   getRulesForNewGame(),
		Turn:    0,
		Guess:   0,
		Correct: true,
	}

	msg, err := s.ChannelMessageSend(channelID, "A new game has begun!")
	if err != nil {
		panic("Could not create new game: " + err.Error())
	}

	newGame.MessageID = msg.ID
	db.Create(&newGame)

	return newGame
}
