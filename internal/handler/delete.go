package handler

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/brandenc40/romannumeral"
	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/emoji"
	"github.com/zaptross/countula/internal/rules"
	"gorm.io/gorm"
)

func GetOnMessageDeletedHandler(db *gorm.DB) func(*discordgo.Session, *discordgo.MessageDelete) {
	slog.Info("Added message deleted handler")
	return func(s *discordgo.Session, del *discordgo.MessageDelete) {
		var delTurn database.Turn
		db.Find(&delTurn, "message_id = ?", del.ID)

		if delTurn.MessageID == "" {
			return // ignore deleted messages that aren't turns
		}

		guessMsg := strconv.Itoa(delTurn.Guess)
		msgEmoji := emoji.CHECK

		if delTurn.Rules&rules.RomanNumeral.Id() == rules.RomanNumeral.Id() {
			g, err := romannumeral.IntToString(delTurn.Guess)
			if err == nil {
				guessMsg = g
			}
		}
		if delTurn.Rules&rules.Jeopardy.Id() == rules.Jeopardy.Id() {
			guessMsg = fmt.Sprintf("What is %d plus 0?", delTurn.Guess)
		}

		// Call out the user who deleted their message, and show their guess
		s.ChannelMessageSend(del.ChannelID, "Trying to hide your guess, huh? :eyes:")
		msg, _ := s.ChannelMessageSend(del.ChannelID, fmt.Sprintf("%s (guessed by <@%s> at %s)", guessMsg, delTurn.UserID, delTurn.CreatedAt.Format(time.RFC3339)))
		db.Model(&delTurn).Update("message_id", msg.ID)

		if !delTurn.Correct {
			msgEmoji = emoji.CROSS
		}
		if delTurn.Rules&rules.KeepyUppies.Id() == rules.KeepyUppies.Id() && delTurn.Correct {
			msgEmoji = emoji.BALLOON
		}
		s.MessageReactionAdd(del.ChannelID, msg.ID, msgEmoji)
	}
}
