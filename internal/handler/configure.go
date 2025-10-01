package handler

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/dgutil"
	"github.com/zaptross/countula/internal/game"
	"github.com/zaptross/countula/internal/verbeage"
	"gorm.io/gorm"
)

const (
	ConfigureCommand    = "!configure-countula"
	ConfigureSubcommand = "configure"

	configureMissingPermissions = "You do not have permission to configure Countula. Please contact your server admin.\n(missing permission: `Manage Webhooks`)"
	configureNonSlashCommand    = "Hey admin! It looks like you tried to configure Countula using a message. Countula now supports slash commands, which are easier to use and more secure. To configure Countula, please use the `/count configure` slash command in the channel you want to set up."
)

func configureSubCommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        ConfigureSubcommand,
		Description: "Configure Countula for this channel",
	}
}

func HandleConfigure(db *gorm.DB, s *discordgo.Session, m *discordgo.MessageCreate) {
	userPerms, err := s.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		log.Default().Println("Could not get user channel permissions", err)
		return
	}

	if !hasPermissionsToConfigure(userPerms) {
		_, err := s.ChannelMessageSendReply(
			m.ChannelID,
			configureMissingPermissions,
			m.Message.Reference(),
		)
		log.Default().Println("Could not send missing permissions message", err)
		return
	}

	_, err = s.ChannelMessageSendReply(m.ChannelID, configureNonSlashCommand, m.Message.Reference())
	if err != nil {
		log.Default().Println("Could not send non-slash command message", err)
		return
	}
}

func configureSlashCommandHandler(db *gorm.DB, s *discordgo.Session, i *discordgo.InteractionCreate) {
	userPerms, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
	if err != nil {
		log.Default().Println("Could not get user channel permissions", err)
		return
	}

	// Check user has permission to add webhooks and bots
	if !hasPermissionsToConfigure(userPerms) {
		_, err := dgutil.InteractionEdit(s, i, configureMissingPermissions)
		if err != nil {
			log.Default().Println("Could not send missing permissions message", err)
		}
		return
	}

	onConfigureMessage, err := verbeage.GetRandomOnConfigureMessage().Message(verbeage.TemplateFields{})
	if err != nil {
		log.Default().Println("Could not get random on configure message", err)
		return
	}

	_, err = s.ChannelMessageSend(i.ChannelID, onConfigureMessage)

	if err != nil {
		log.Default().Println("Could not send on configure message", err)
		return
	}

	database.ConfigureFromMessage(db, i)
	game.CreateNewGame(db, s, i.ChannelID, i.GuildID)
}

func hasPermissionsToConfigure(perms int64) bool {
	return perms&discordgo.PermissionManageWebhooks == discordgo.PermissionManageWebhooks
}
