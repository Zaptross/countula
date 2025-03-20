package handler

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/emoji"
	"github.com/zaptross/countula/internal/statistics"
	"gorm.io/gorm"
)

const (
	StatsSubcommand = "stats"
	StatsButton     = "stats_show"
)

func statsSlashCommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        StatsSubcommand,
		Description: "Display your statistics",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "source",
				Description: "Global or channel statistics",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Global",
						Value: "global",
					},
					{
						Name:  "Channel",
						Value: "channel",
					},
				},
			},
		},
	}
}

func statsSlashCommandHandler(db *gorm.DB, s *discordgo.Session, i *discordgo.InteractionCreate) {
	source := i.ApplicationCommandData().Options[0].Options[0].StringValue()

	var statsMessage string
	if source == "global" {
		statsMessage = statistics.DisplayGlobal(db, i.Member.User.ID)
	} else {
		statsMessage = statistics.Display(db, i.Member.User.ID, i.ChannelID)
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &statsMessage,
		Components: &[]discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					&discordgo.Button{
						Emoji: discordgo.ComponentEmoji{
							Name: emoji.EYES,
						},
						Label:    "Show Channel",
						Style:    discordgo.PrimaryButton,
						CustomID: StatsButton + source,
					},
				},
			},
		},
	})
}

func statsButtonHandler(db *gorm.DB, s *discordgo.Session, i *discordgo.InteractionCreate, customID string) {
	res := i.Message.Content
	source := "globally"

	if customID == StatsButton+"channel" {
		source = "in this channel"
	}

	s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("%s's Stats %s:\n%s", i.Member.Mention(), source, res))

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "You can dismiss this message now.",
		},
	})
}
