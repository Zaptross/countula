package dgutil

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func InteractionRespond(s *discordgo.Session, r *discordgo.InteractionCreate, content string) error {
	// Create a response to the interaction
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral, // Make the response ephemeral
		},
	}

	// Respond to the interaction
	err := s.InteractionRespond(r.Interaction, response)
	if err != nil {
		slog.Error("failed to respond to interaction", "error", err)
	}

	return err
}

func InteractionEdit(s *discordgo.Session, r *discordgo.InteractionCreate, content string) (*discordgo.Message, error) {
	// Create a response to edit the interaction
	response := &discordgo.WebhookEdit{
		Content: &content,
	}

	// Edit the interaction response
	return s.InteractionResponseEdit(r.Interaction, response)
}

func InteractionEditWithButtons(s *discordgo.Session, r *discordgo.InteractionCreate, content string, buttons []*discordgo.Button) (*discordgo.Message, error) {
	var components []discordgo.MessageComponent
	for _, button := range buttons {
		components = append(components, button)
	}

	// Create a response to edit the interaction with buttons
	response := &discordgo.WebhookEdit{
		Content: &content,
		Components: &[]discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: components,
			},
		},
	}

	// Edit the interaction response
	msg, err := s.InteractionResponseEdit(r.Interaction, response)

	if err != nil {
		slog.Error("failed to edit interaction with buttons", "error", err)
	}

	return msg, err
}
