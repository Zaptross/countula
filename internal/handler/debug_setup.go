package handler

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/rules"
	"gorm.io/gorm"
)

const (
	debugSetupSubCommand = "debug-setup"
)

func debugSetupSlashCommand() *discordgo.ApplicationCommandOption {
	rules := getRulesOptions(4, 1)

	options := append([]*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "The user to set as the last counter",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "current",
			Description: "The current number in the game",
			Required:    false,
		},
	}, rules...)

	// Sort options so required ones are first
	slices.SortFunc(options, func(a, b *discordgo.ApplicationCommandOption) int {
		if a.Required && !b.Required {
			return -1
		}
		if !a.Required && b.Required {
			return 1
		}
		return 0
	})

	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        debugSetupSubCommand,
		Description: "Setup a debug game",
		Options:     options,
	}
}

func getRulesOptions(count, required int) []*discordgo.ApplicationCommandOption {
	type RuleChoice struct {
		Name string
		ID   int
	}

	choices := lo.FilterMap(rules.AllRules, func(rule rules.Rule, _ int) (RuleChoice, bool) {
		return RuleChoice{
			Name: rule.Name(),
			ID:   rule.Id(),
		}, rule.Id() != rules.MathsRuleId && rule.Id() != rules.NoValidateRuleId
	})

	options := []*discordgo.ApplicationCommandOption{}

	for i := 0; i < count; i++ {
		options = append(options, &discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        fmt.Sprintf("rule-%d", i+1),
			Description: "One rule to add to the debug setup",
			Required:    i < required,
			Choices: lo.Map(choices, func(choice RuleChoice, _ int) *discordgo.ApplicationCommandOptionChoice {
				return &discordgo.ApplicationCommandOptionChoice{
					Name:  choice.Name,
					Value: choice.ID,
				}
			}),
		})
	}

	return options
}

func debugSetupSlashCommandHandler(db *gorm.DB, s *discordgo.Session, i *discordgo.InteractionCreate) {
	settings := i.ApplicationCommandData().Options[0].Options

	rulesCombined := 0
	for _, option := range settings {
		if strings.HasPrefix(option.Name, "rule-") {
			if value := option.IntValue(); value != 0 {
				rulesCombined |= int(value)
			}
		}
	}

	countAsUserID := ""
	countAsUserName := "no one"
	userOption, ok := lo.Find(settings, func(option *discordgo.ApplicationCommandInteractionDataOption) bool {
		return option.Name == "user"
	})
	if ok && userOption != nil && userOption.UserValue(s) != nil {
		countAsUser := userOption.UserValue(s)
		countAsUserID = countAsUser.ID
		countAsUserName = countAsUser.Username
	}

	current := 0
	currentOption, ok := lo.Find(settings, func(option *discordgo.ApplicationCommandInteractionDataOption) bool {
		return option.Name == "current"
	})
	if ok && currentOption != nil {
		current = int(currentOption.IntValue())
	}

	currentTurn := database.GetCurrentTurn(db, i.ChannelID)

	debugTurn := &database.Turn{
		ChannelID: i.ChannelID,
		Correct:   true,
		Game:      currentTurn.Game + 1,
		Guess:     current,
		Rules:     rulesCombined,
		Turn:      1,
		UserID:    countAsUserID,
	}

	db.Create(debugTurn)

	_, err := s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Debug game set up with current: %d, user %s as last counter\n\n%s", current, countAsUserName, strings.Join(rules.GetRuleTextsForGame(*debugTurn), "\n")))
	if err != nil {
		panic("Could not create debug game: " + err.Error())
	}
}
