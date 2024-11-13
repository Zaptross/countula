package rules

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
)

type GoodyTwoShoesRule struct {
	RuleWeight
	id       int
	ruleType string
}

func (gtsr GoodyTwoShoesRule) Id() int {
	return gtsr.id
}
func (gtsr GoodyTwoShoesRule) Name() string {
	return "Goody Two Shoes"
}
func (gtsr GoodyTwoShoesRule) Description() string {
	return "Take two steps forward, then one step back."
}
func (gtsr GoodyTwoShoesRule) Weight() int {
	return gtsr.Current
}
func (gtsr GoodyTwoShoesRule) WithWeight(weight int) Rule {
	return GoodyTwoShoesRule{
		id:         gtsr.id,
		ruleType:   gtsr.ruleType,
		RuleWeight: SetupWeight(weight),
	}
}
func (gtsr GoodyTwoShoesRule) Type() string {
	return gtsr.ruleType
}
func (gtsr GoodyTwoShoesRule) OnNewGame(db *gorm.DB, s *discordgo.Session, ng database.Turn, channelID string) {
}

func (gtsr GoodyTwoShoesRule) Validate(db *gorm.DB, lastTurn database.Turn, msg discordgo.Message, guess int) bool {
	if lastTurn.Turn%2 == 1 {
		return guess == lastTurn.Guess-1
	}
	return guess == lastTurn.Guess+2
}

var (
	GoodyTwoShoes = (func() ValidateRule {
		gtsr := GoodyTwoShoesRule{
			id:         GoodyTwoShoesRuleId,
			RuleWeight: SetupWeight(GoodyTwoShoesRuleWeight),
			ruleType:   CountType,
		}

		registerRule(gtsr)

		return gtsr
	})()
)
