package handler

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/utils"
	"gorm.io/gorm"
)

func GetSlashCommand() *discordgo.ApplicationCommand {
	optional := []*discordgo.ApplicationCommandOption{}
	if utils.IsDevelopment() {
		optional = append(optional, debugSetupSlashCommand())
	}

	return &discordgo.ApplicationCommand{
		Name:        "count",
		Description: "Configure the counting channel",
		Options: append([]*discordgo.ApplicationCommandOption{
			settingsSlashCommand(),
			maintenanceModeSlashCommand(),
			statsSlashCommand(),
			configureSubCommand(),
		}, optional...),
	}
}

func GetOnInteractionHandler(db *gorm.DB, adminChannelId string) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionMessageComponent:
			handleMessageComponentInteraction(db, s, i)
		case discordgo.InteractionApplicationCommand:
			handleApplicationCommandInteraction(db, s, i, adminChannelId)
		}
	}
}

func handleApplicationCommandInteraction(db *gorm.DB, s *discordgo.Session, i *discordgo.InteractionCreate, adminChannelId string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: ":gear: Working...",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	// switch on subcommand name
	switch i.ApplicationCommandData().Options[0].Name {
	case SettingsSubcommand:
		settingsSlashCommandHandler(db, s, i)
	case MaintenanceModeSubcommand:
		maintenanceModeSlashCommandHandler(db, s, i, adminChannelId)
	case StatsSubcommand:
		statsSlashCommandHandler(db, s, i)
	case ConfigureSubcommand:
		configureSlashCommandHandler(db, s, i)
	// Debug commands
	case debugSetupSubCommand:
		debugSetupSlashCommandHandler(db, s, i)
	}
}

// Handle button clicks, etc.
func handleMessageComponentInteraction(db *gorm.DB, s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	switch true {
	case strings.Contains(customID, StatsButton):
		statsButtonHandler(db, s, i, customID)
	}
}

func replyToInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}
