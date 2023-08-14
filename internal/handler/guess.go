package handler

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/emoji"
	"github.com/zaptross/countula/internal/game"
	"github.com/zaptross/countula/internal/rules"
	"github.com/zaptross/countula/internal/verbeage"
	"gorm.io/gorm"
)

func handleGuess(db *gorm.DB, s *discordgo.Session, m *discordgo.MessageCreate, config Config) {
	turn := database.GetCurrentTurn(db)
	rules := rules.GetRulesForTurn(turn)

	var guess int
	for _, pvr := range rules.PreValidateRules {
		var err error
		guess, err = pvr.PreValidate(db, s, *m.Message)
		if err != nil {
			return // assume that rules which do not match are not guesses
		}
	}

	for _, vr := range rules.ValidateRules {
		if !vr.Validate(db, turn, *m.Message, guess) {
			println(fmt.Sprintf("Failed validation at rule: %s, with guess: %d, after last guess: %d", vr.Name(), guess, turn.Guess))
			failValidate(db, s, m, turn, guess, config.CountingChannel)
			return
		}
	}

	database.CreateTurnFromContext(db, s, m, turn, guess, true)
	go s.MessageReactionAdd(m.ChannelID, m.Message.ID, emoji.CHECK)
}

func failValidate(db *gorm.DB, s *discordgo.Session, m *discordgo.MessageCreate, lastTurn database.Turn, guess int, channelID string) {
	go s.MessageReactionAdd(m.ChannelID, m.Message.ID, emoji.CROSS)
	failMessageSend(s, m)
	database.CreateTurnFromContext(db, s, m, lastTurn, guess, false)
	game.CreateNewGame(db, s, channelID)
}

func failMessageSend(s *discordgo.Session, m *discordgo.MessageCreate) {
	failMessage := verbeage.GetRandomFail()

	username := m.Author.Username
	gm, err := s.GuildMember(m.GuildID, m.Author.ID)

	if err == nil {
		username = gm.Nick
	}

	tf := verbeage.TemplateFields{
		Username: username,
	}

	if failMessage.Reply != nil {
		t, err := failMessage.Reply(tf)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "**ERROR**: DOES NOT COMPUTE")
		}

		go s.ChannelMessageSendReply(m.ChannelID, t, m.Message.Reference())
	}

	if failMessage.Message != nil {
		t, err := failMessage.Message(tf)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "**ERROR**: DOES NOT COMPUTE")
		}

		go s.ChannelMessageSend(m.ChannelID, t)
	}
}
