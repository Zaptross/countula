package rules

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
)

type GuessNormallyRule struct {
	RuleWeight
	id       int
	ruleType string
}

func (gnr GuessNormallyRule) Id() int {
	return gnr.id
}
func (gnr GuessNormallyRule) Name() string {
	return "Guess Normally"
}
func (gnr GuessNormallyRule) Description() string {
	return "You **must** guess normally, i.e. your guess as a number at the start of your message."
}
func (gnr GuessNormallyRule) Weight() int {
	return gnr.Current
}
func (gnr GuessNormallyRule) WithWeight(weight int) Rule {
	return GuessNormallyRule{
		id:         gnr.id,
		ruleType:   gnr.ruleType,
		RuleWeight: SetupWeight(weight),
	}
}
func (gnr GuessNormallyRule) Type() string {
	return gnr.ruleType
}
func (gnr GuessNormallyRule) OnNewGame(_ *gorm.DB, _ *discordgo.Session, _ database.Turn, _ string) {}

func (gnr GuessNormallyRule) PreValidate(db *gorm.DB, dg *discordgo.Session, msg discordgo.Message) (int, error) {
	digits := strings.Split(msg.Content, " ")[0]
	guess, err := strconv.Atoi(digits)

	if err != nil {
		return 0, err
	}

	return guess, nil
}

var (
	GuessNormally = (func() Rule {
		gnr := GuessNormallyRule{
			id:         GuessNormallyRuleId,
			RuleWeight: SetupWeight(GuessNormallyRuleWeight),
			ruleType:   PreValidateType,
		}

		registerRule(gnr)

		return gnr
	})()
)
