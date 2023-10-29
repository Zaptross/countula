package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
)

func GetMessageHandler(db *gorm.DB) func(*discordgo.Session, *discordgo.MessageCreate) {
	var serverConfigs []database.ServerConfig
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Ignore all messages that are not in a guild
		if m.GuildID == "" {
			return
		}

		if serverConfigs == nil {
			serverConfigs = database.GetAllServerConfigs(db)
		}

		// Ignore all messages that are not in a configured counting channel
		cfg := findConfigForServer(serverConfigs, m.ChannelID)
		if cfg == nil {

			// Ensure configure command is handled only once
			if m.Content == ConfigureCommand {
				HandleConfigure(db, s, m)

				// Invalidate the cache of server configs
				serverConfigs = nil
			}

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
		handleGuess(db, s, m, cfg)
	}
}

func findConfigForServer(serverConfigs []database.ServerConfig, channelID string) *database.ServerConfig {
	for _, cfg := range serverConfigs {
		if cfg.CountingChannelID == channelID {
			return &cfg
		}
	}
	return nil
}
