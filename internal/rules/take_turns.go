package rules

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
)

type TakeTurnsRule struct {
	RuleWeight
	id       int
	ruleType string
}

func (ttr TakeTurnsRule) Id() int {
	return ttr.id
}
func (ttr TakeTurnsRule) Name() string {
	return "Take Turns"
}
func (ttr TakeTurnsRule) Description() string {
	return "You **must** allow another player to take a turn after yours."
}
func (ttr TakeTurnsRule) Weight() int {
	return ttr.Current
}
func (ttr TakeTurnsRule) SetWeight(weight int) {
	ttr.Current = weight
}
func (ttr TakeTurnsRule) Type() string {
	return ttr.ruleType
}
func (ttr TakeTurnsRule) OnNewGame(_ *gorm.DB, _ *discordgo.Session, _ database.Turn, _ string) {}

func (ttr TakeTurnsRule) Validate(db *gorm.DB, lastTurn database.Turn, msg discordgo.Message, guess int) bool {
	return lastTurn.UserID != msg.Author.ID
}

var (
	TakeTurns = (func() Rule {
		ttr := TakeTurnsRule{
			id:         TakeTurnsRuleId,
			RuleWeight: Weights(TakeTurnsRuleWeight),
			ruleType:   ValidateType,
		}

		registerRule(ttr)

		return ttr
	})()
)
