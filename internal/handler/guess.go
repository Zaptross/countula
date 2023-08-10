package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/rules"
	"github.com/zaptross/countula/internal/verbeage"
	"gorm.io/gorm"
)

func handleGuess(db *gorm.DB, s *discordgo.Session, m *discordgo.MessageCreate) {
	turn := database.GetCurrentTurn(db)
	rules := rules.GetRulesForTurn(turn)

	var guess int
	for _, pvr := range rules.PreValidateRules {
		var err error
		guess, err = pvr.PreValidate(db, s, *m.Message, guess)
		if err != nil {
			failPreValidate(s, m)
			return
		}
	}
}

func failPreValidate(s *discordgo.Session, m *discordgo.MessageCreate) {
	failMessage := verbeage.GetRandomFail()

	tf := verbeage.TemplateFields{
		Username: m.Author.Username,
	}

	if failMessage.Reply != nil {
		t, err := failMessage.Reply(tf)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "**ERROR**: DOES NOT COMPUTE")
		}

		s.ChannelMessageSendReply(m.ChannelID, t, m.Message.Reference())
	}

	if failMessage.Message != nil {
		t, err := failMessage.Message(tf)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "**ERROR**: DOES NOT COMPUTE")
		}

		s.ChannelMessageSend(m.ChannelID, t)
	}
}
