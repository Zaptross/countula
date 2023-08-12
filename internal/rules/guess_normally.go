package rules

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type GuessNormallyRule struct {
	id       int
	weight   int
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
	return gnr.weight
}
func (gnr GuessNormallyRule) Type() string {
	return gnr.ruleType
}
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
			id:       GuessNormallyRuleId,
			weight:   100,
			ruleType: PreValidateType,
		}

		registerRule(gnr)

		return gnr
	})()
)
