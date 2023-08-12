package rules

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/utils"
	"gorm.io/gorm"
)

type Rule interface {
	Id() int
	Name() string
	Weight() int
	Description() string
	Type() string
}

// The purpose of a prevalidate rule is to take the message and extract the guess from it, and return the guess as an int
type PreValidateRule interface {
	Rule
	// PreValidate returns the guess as an int and an error if the guess is invalid
	PreValidate(db *gorm.DB, dg *discordgo.Session, msg discordgo.Message) (int, error)
}

// The purpose of a validate rule is to take the guess and validate that it meets the rule criteria, returning a bool
type ValidateRule interface {
	Rule
	// Validate returns true if the guess is correct
	Validate(db *gorm.DB, dg *discordgo.Session, msg discordgo.Message, guess int) bool
}

const (
	PreValidateType = "PreValidate"
	ValidateType    = "Validate"
)

var (
	AllRules         = []Rule{}
	PreValidateRules = []Rule{}
	ValidateRules    = []Rule{}
)

func registerRule(r Rule) {
	AllRules = append(AllRules, r)
	switch r.Type() {
	case PreValidateType:
		PreValidateRules = append(PreValidateRules, r)
	case ValidateType:
		ValidateRules = append(ValidateRules, r)
	}
}

type RulesForTurn struct {
	PreValidateRules []PreValidateRule
	ValidateRules    []ValidateRule
}

func GetRulesForTurn(g database.Turn) RulesForTurn {
	rules := RulesForTurn{
		PreValidateRules: []PreValidateRule{},
		ValidateRules:    []ValidateRule{},
	}

	for _, r := range AllRules {
		if g.Rules&r.Id() == r.Id() {
			switch r.Type() {
			case PreValidateType:
				rules.PreValidateRules = append(rules.PreValidateRules, r.(PreValidateRule))
			case ValidateType:
				rules.ValidateRules = append(rules.ValidateRules, r.(ValidateRule))
			}
		}
	}

	return rules
}

func GetRuleTextsForGame(g database.Turn) []string {
	rules := []string{}

	for _, r := range AllRules {
		if g.Rules&r.Id() == r.Id() {
			rules = append(rules, fmt.Sprintf("**%s**: %s\n", r.Name(), r.Description()))
		}
	}

	return rules
}

func GetRandomPreValidateRule() PreValidateRule {
	return utils.RandFrom(PreValidateRules).(PreValidateRule)
}

func GetRandomValidateRule() ValidateRule {
	return utils.RandFrom(ValidateRules).(ValidateRule)
}
