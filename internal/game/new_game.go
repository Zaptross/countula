package game

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/rules"
	"github.com/zaptross/countula/internal/verbeage"
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

	ruleMessage := verbeage.GetRandomRuleMessage()

	rm, err := ruleMessage.Message(verbeage.TemplateFields{})

	if err != nil {
		panic("Could not create new game: " + err.Error())
	}

	rules := rules.GetRuleTextsForGame(newGame)

	msg, err := s.ChannelMessageSend(channelID, fmt.Sprintf("%s\n%s", rm, strings.Join(rules, "\n")))
	if err != nil {
		panic("Could not create new game: " + err.Error())
	}

	newGame.MessageID = msg.ID
	db.Create(&newGame)

	return newGame
}
