package rules

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
)

type GoodyThreeShoesRule struct {
	RuleWeight
	id       int
	ruleType string
}

func (gtsr GoodyThreeShoesRule) Id() int {
	return gtsr.id
}
func (gtsr GoodyThreeShoesRule) Name() string {
	return "Goody... Three Shoes?"
}
func (gtsr GoodyThreeShoesRule) Description() string {
	return "Take three steps forward, then one step back, then one step back."
}
func (gtsr GoodyThreeShoesRule) Weight() int {
	return gtsr.Current
}
func (gtsr GoodyThreeShoesRule) WithWeight(weight int) Rule {
	return GoodyThreeShoesRule{
		id:         gtsr.id,
		ruleType:   gtsr.ruleType,
		RuleWeight: SetupWeight(weight),
	}
}
func (gtsr GoodyThreeShoesRule) Type() string {
	return gtsr.ruleType
}
func (gtsr GoodyThreeShoesRule) OnNewGame(db *gorm.DB, s *discordgo.Session, ng database.Turn, channelID string) {
}
func (gtsr GoodyThreeShoesRule) OnFailure(fc *FailureContext) *FailureContext {
	return fc
}

func (gtsr GoodyThreeShoesRule) Validate(db *gorm.DB, lastTurn database.Turn, msg discordgo.Message, guess int) bool {
	if lastTurn.Turn%3 > 0 {
		return guess == lastTurn.Guess-1
	}
	return guess == lastTurn.Guess+3
}

var (
	GoodyThreeShoes = (func() ValidateRule {
		gtsr := GoodyThreeShoesRule{
			id:         GoodyThreeShoesRuleId,
			RuleWeight: SetupWeight(GoodyThreeShoesRuleWeight),
			ruleType:   CountType,
		}

		registerRule(gtsr)

		return gtsr
	})()
)
