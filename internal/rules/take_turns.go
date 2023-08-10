package rules

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
)

type TakeTurnsRule struct {
	id       int
	weight   int
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
	return ttr.weight
}
func (ttr TakeTurnsRule) Type() string {
	return ttr.ruleType
}

func (ttr TakeTurnsRule) Validate(db *gorm.DB, dg *discordgo.Session, msg discordgo.Message, guess int) bool {
	var lastTurn database.Turn
	db.Last(&lastTurn)
	return lastTurn.UserID != msg.Author.ID
}

var (
	TakeTurns = (func() Rule {
		ttr := TakeTurnsRule{
			id:       TakeTurnsRuleId,
			weight:   100,
			ruleType: ValidateType,
		}

		registerRule(ttr)

		return ttr
	})()
)
