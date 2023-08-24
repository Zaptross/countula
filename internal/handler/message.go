package handler

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type Config struct {
	AdminRoleId     string
	CountingChannel string
}

func GetMessageHandler(db *gorm.DB, config Config) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Ignore all messages that are not in a guild
		if m.GuildID == "" {
			return
		}

		if len(m.Content) == 0 {
			println("Message has no content", m.Content)
			return
		}

		// Messages that start with ! are commands, and should be handled by the command handler
		if m.Content[0] == '!' {
			handleCommand(db, s, m)
			return
		}

		// Messages that are not commands are guesses, and should be handled by the guess handler
		handleGuess(db, s, m, config)
	}
}
