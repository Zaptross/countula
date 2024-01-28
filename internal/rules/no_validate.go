package rules

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
)

type NoValidateRule struct {
	RuleWeight
	id       int
	ruleType string
}

func (nvr NoValidateRule) Id() int {
	return nvr.id
}
func (nvr NoValidateRule) Name() string {
	return ""
}
func (nvr NoValidateRule) Description() string {
	return ""
}
func (nvr NoValidateRule) Weight() int {
	return nvr.Current
}
func (nvr NoValidateRule) SetWeight(weight int) {
	nvr.Current = weight
}
func (nvr NoValidateRule) Type() string {
	return nvr.ruleType
}
func (nvr NoValidateRule) OnNewGame(_ *gorm.DB, _ *discordgo.Session, _ database.Turn, _ string) {}

func (nvr NoValidateRule) Validate(_ *gorm.DB, lastTurn database.Turn, _ discordgo.Message, _ int) bool {
	return true
}

var (
	NoValidate = (func() ValidateRule {
		nvr := NoValidateRule{
			id:         NoValidateRuleId,
			RuleWeight: Weights(30),
			ruleType:   ValidateType,
		}

		registerRule(nvr)

		return nvr
	})()
)
