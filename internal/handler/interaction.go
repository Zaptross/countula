package handler

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func GetOnInteractionHandler(db *gorm.DB) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: ":gear: Working...",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		if i.Member.Permissions&discordgo.PermissionManageWebhooks == 0 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You do not have permission to use this command.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		// switch on subcommand name
		switch i.ApplicationCommandData().Options[0].Name {
		case SettingsSubcommand:
			settingsSlashCommandHandler(db, s, i)
		}
	}
}

func GetSlashCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "count",
		Description: "Configure the counting channel",
		Options: []*discordgo.ApplicationCommandOption{
			settingsSlashCommand(),
		},
	}
}

func replyToInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}
