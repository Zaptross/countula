package handler

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func GetMessageHandler(db *gorm.DB) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Ignore all messages that are not in a guild
		if m.GuildID == "" {
			return
		}

		// Messages that start with ! are commands, and should be handled by the command handler
		if m.Content[0] == '!' {
			handleCommand(db, s, m)
			return
		}

		// Messages that are not commands are guesses, and should be handled by the guess handler
		handleGuess(db, s, m)
	}
}
