package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/emoji"
	"gorm.io/gorm"
)

type StateCommand struct{}

const (
	StateCommandName = "!state"
)

func (c StateCommand) Execute(db *gorm.DB, s *discordgo.Session, m *discordgo.MessageCreate) {
	turn := database.GetCurrentTurn(db, m.ChannelID)
	highScoreTurn := database.GetHighScoreTurn(db, m.ChannelID)

	highScoreMessage := "How do you not have a high score yet?! Get counting!"
	if highScoreTurn.Turn > 0 {
		highScoreMessage = fmt.Sprintf("The last number was %d, and the high score is: %d %s (turn %d of that game)",
			turn.Guess,
			highScoreTurn.Guess,
			emoji.HIGH_SCORE,
			highScoreTurn.Turn,
		)
	}

	_, err := s.ChannelMessageSendReply(
		m.ChannelID,
		highScoreMessage,
		m.Message.Reference(),
	)
	if err != nil {
		panic("Could not create new game: " + err.Error())
	}
}

func (c StateCommand) Describe() string {
	return "I shall tell you where we left off."
}
