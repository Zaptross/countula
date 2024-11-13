package handler

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/emoji"
	"github.com/zaptross/countula/internal/game"
	"github.com/zaptross/countula/internal/rules"
	"github.com/zaptross/countula/internal/statistics"
	"github.com/zaptross/countula/internal/verbeage"
	"gorm.io/gorm"
)

func handleGuess(db *gorm.DB, s *discordgo.Session, m *discordgo.MessageCreate, config *database.ServerConfig) {
	turn := database.GetCurrentTurn(db, m.ChannelID)
	currentTurnRules := rules.GetRulesForTurn(turn)

	var guess int
	for _, pvr := range currentTurnRules.PreValidateRules {
		var err error
		guess, err = pvr.PreValidate(db, s, *m.Message)
		if err != nil {
			return // assume that rules which do not match are not guesses
		}
	}

	successEmoji := emoji.CHECK

	for _, vr := range currentTurnRules.ValidateRules {
		if !vr.Validate(db, turn, *m.Message, guess) {
			slog.Info("Incorrect guess", "rule", vr.Name(), "guess", guess, "last_guess", turn.Guess)
			failValidate(
				vr.OnFailure(&rules.FailureContext{
					DB:             db,
					DG:             s,
					Msg:            m,
					FailureMessage: verbeage.GetRandomFail(),
					LastTurn:       turn,
					Guess:          guess,
					ChannelID:      config.CountingChannelID,
					Emoji:          emoji.CROSS,
				}),
			)
			return
		}
		successEmoji = rules.OverrideSuccessEmoji(vr.Id())
	}

	hst := database.GetHighScoreTurn(db, m.ChannelID)
	ct := database.CreateTurnFromContext(db, s, m, turn, guess, true)

	go checkHighScore(s, m, ct, hst)
	go s.MessageReactionAdd(m.ChannelID, m.Message.ID, successEmoji)

	go statistics.Collect(db, s, m, config.CountingChannelID, ct)
}

func checkHighScore(s *discordgo.Session, m *discordgo.MessageCreate, ct database.Turn, hs database.Turn) {
	if ct.Turn > hs.Turn {
		go s.MessageReactionAdd(m.Message.ChannelID, ct.MessageID, emoji.HIGH_SCORE)
		go s.MessageReactionRemove(m.Message.ChannelID, hs.MessageID, emoji.HIGH_SCORE, s.State.User.ID)
	}
}

func failValidate(ctx *rules.FailureContext) {
	go ctx.DG.MessageReactionAdd(ctx.Msg.ChannelID, ctx.Msg.ID, ctx.Emoji)
	failMessageSend(ctx.DG, ctx.Msg, ctx.FailureMessage)
	ct := database.CreateTurnFromContext(ctx.DB, ctx.DG, ctx.Msg, ctx.LastTurn, ctx.Guess, false)
	game.CreateNewGame(ctx.DB, ctx.DG, ctx.ChannelID, ctx.Msg.GuildID)
	go statistics.Collect(ctx.DB, ctx.DG, ctx.Msg, ctx.ChannelID, ct)
}

func failMessageSend(s *discordgo.Session, m *discordgo.MessageCreate, failMessage verbeage.ResponseParts) {
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
