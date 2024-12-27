package handler

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
)

const (
	MaintenanceModeSubcommand = "maintenance-mode"
)

var (
	isMaintenanceModeEnabled = false
)

func StartupCheckMaintenanceMode(db *gorm.DB) {
	serviceCfg := database.GetServiceConfig(db)
	isMaintenanceModeEnabled = serviceCfg.MaintenanceMode
	slog.Info("Maintenance mode is", "enabled", isMaintenanceModeEnabled)
}

func IsMaintenanceModeEnabled() bool {
	return isMaintenanceModeEnabled
}

func maintenanceModeSlashCommandHandler(db *gorm.DB, s *discordgo.Session, i *discordgo.InteractionCreate, adminChannelId string) {
	if adminChannelId != "" && i.ChannelID != adminChannelId {
		replyToInteraction(s, i, "You do not have permission to change maintenance mode")
		return
	}

	if adminChannelId == "" && (i.GuildID == "" || i.Member.Permissions&discordgo.PermissionManageWebhooks == 0) {
		replyToInteraction(s, i, "You do not have permission to change maintenance mode")
		return
	}

	enable := i.ApplicationCommandData().Options[0].Options[0].BoolValue()

	database.SetMaintenanceMode(db, enable)
	isMaintenanceModeEnabled = enable

	replyToInteraction(s, i, fmt.Sprintf("Maintenance mode is now %t", enable))
}

func maintenanceModeSlashCommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:         MaintenanceModeSubcommand,
		Description:  "Enable or disable maintenance mode",
		Type:         discordgo.ApplicationCommandOptionSubCommand,
		ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeDM},
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionBoolean,
				Name:         "enable",
				Description:  "Enable maintenance mode",
				ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeDM},
				Required:     true,
			},
		},
	}
}
