package handler

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func GetOnInteractionHandler(db *gorm.DB, adminChannelId string) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		}
	}
}

func GetSlashCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "count",
		Description: "Configure the counting channel",
		Options: []*discordgo.ApplicationCommandOption{
			settingsSlashCommand(),
			maintenanceModeSlashCommand(),
		},
	}
}

func replyToInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}
