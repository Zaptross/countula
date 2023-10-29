package handler

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/game"
	"github.com/zaptross/countula/internal/verbeage"
	"gorm.io/gorm"
)

const (
	ConfigureCommand = "!configure-countula"
)

func HandleConfigure(db *gorm.DB, s *discordgo.Session, m *discordgo.MessageCreate) {
	userPerms, err := s.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		log.Default().Println("Could not get user channel permissions", err)
		return
	}

	// Check user has permission to add webhooks and bots
	if userPerms&discordgo.PermissionManageWebhooks != discordgo.PermissionManageWebhooks {
		s.ChannelMessageSendReply(
			m.ChannelID,
			"You do not have permission to configure Countula. Please contact your server admin.\n(missing permission: `Manage Webhooks`)",
			m.Message.Reference(),
		)
		return
	}

	t, err := verbeage.GetRandomOnConfigureMessage().Message(verbeage.TemplateFields{})
	if err != nil {
		log.Default().Println("Could not get random on configure message", err)
		return
	}

	_, err = s.ChannelMessageSendReply(m.ChannelID, t, m.Message.Reference())

	if err != nil {
		log.Default().Println("Could not send on configure message", err)
		return
	}

	database.ConfigureFromMessage(db, m)
	game.CreateNewGame(db, s, m.ChannelID)
}
