package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/rules"
	"gorm.io/gorm"
)

const (
	SettingsSubcommand    = "settings"
	SettingsSubcommandGet = "get"
)

func settingsSlashCommandHandler(db *gorm.DB, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := database.GetServerConfig(db, i.GuildID)

	if config == nil {
		replyToInteraction(s, i, "No configuration found for this server.")
		return
	}

	if len(i.ApplicationCommandData().Options[0].Options) == 0 {
		replyToInteraction(s, i, fmt.Sprintf("No settings provided.\n\n To configure the rules, use `/count %s %s` or visit the configurator at: https://configurator.countula.zaptross.com", SettingsSubcommand, SettingsSubcommandGet))
		return
	}

	actionString := i.ApplicationCommandData().Options[0].Options[0].StringValue()

	existingSettings := database.GetRuleSettingsForGuild(db, i.GuildID)
	if actionString == SettingsSubcommandGet {
		if len(existingSettings) == 0 {
			replyToInteraction(s, i, "The settings for this server are the default. To configure the rules, visit the configurator at: https://configurator.countula.zaptross.com")
			return
		}

		allSettings := rules.ApplyWeightsToRules(rules.AllRules, existingSettings)

		var settingStrings []string
		for _, setting := range allSettings {
			settingStrings = append(settingStrings, fmt.Sprintf("%d:%d", setting.Id(), setting.Weight()))
		}

		replyToInteraction(s, i, fmt.Sprintf("Load this server's settings into the configurator:\nhttps://configurator.countula.zaptross.com/?s=%s", strings.Join(settingStrings, ",")))
		return
	}

	settingsString := i.ApplicationCommandData().Options[0].Options[1].StringValue()
	settings := lo.FilterMap(strings.Split(settingsString, ","), func(s string, _ int) ([]string, bool) {
		rule := strings.Split(s, ":")
		return rule, len(rule) == 2
	})

	if len(settings) == 0 {
		replyToInteraction(s, i, "No settings provided.\n\n To configure the rules, visit the configurator at: https://configurator.countula.zaptross.com")
		return
	}

	updates := []string{}

	for _, setting := range settings {
		id, iErr := strconv.Atoi(setting[0])
		weight, wErr := strconv.Atoi(setting[1])

		if iErr != nil || wErr != nil {
			replyToInteraction(s, i, fmt.Sprintf("Invalid setting: %s:%s", setting[0], setting[1]))
			return
		}

		database.CreateRuleSettingForGuild(db, i.GuildID, id, weight)

		rule := rules.GetRuleByID(id)
		oldWeight := rule.Weight()
		existingSetting, ok := lo.Find(existingSettings, func(settingRow database.RuleSetting) bool {
			return settingRow.RuleID == id
		})

		if ok {
			oldWeight = existingSetting.Weight
		}

		if rule != nil && oldWeight != weight && rule.Weight() != weight {
			updates = append(updates, fmt.Sprintf("**%s**: %d -> %d", rule.Name(), oldWeight, weight))
		}
	}
	updateMessage := "No settings updated."
	if len(updates) > 0 {
		updateMessage = fmt.Sprintf("Settings updated:\n%s", strings.Join(updates, "\n"))
	}

	replyToInteraction(s, i, updateMessage)
}

func settingsSlashCommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        SettingsSubcommand,
		Description: "Configure the rules settings.",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "action",
				Description: "Get or Set the settings.",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "get", Value: "get"},
					{Name: "set", Value: "set"},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "settings",
				Description: "the rule settings to set",
				Required:    false,
			},
		},
	}
}
