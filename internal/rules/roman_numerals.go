package rules

import (
	"strings"

	rom "github.com/brandenc40/romannumeral"
	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
)

type RomanNumeralRule struct {
	RuleWeight
	id       int
	ruleType string
}

func (rnr RomanNumeralRule) Id() int {
	return rnr.id
}
func (rnr RomanNumeralRule) Name() string {
	return "Roman Numeral"
}
func (rnr RomanNumeralRule) Description() string {
	return "You **must** guess in Roman Numerals."
}
func (rnr RomanNumeralRule) Weight() int {
	return rnr.Current
}
func (rnr RomanNumeralRule) SetWeight(weight int) {
	rnr.Current = weight
}
func (rnr RomanNumeralRule) Type() string {
	return rnr.ruleType
}
func (rnr RomanNumeralRule) OnNewGame(_ *gorm.DB, _ *discordgo.Session, _ database.Turn, _ string) {}

func (rnr RomanNumeralRule) PreValidate(db *gorm.DB, dg *discordgo.Session, msg discordgo.Message) (int, error) {
	digits := strings.Split(msg.Content, " ")[0]
	guess, err := rom.StringToInt(digits)

	if err != nil {
		return 0, err
	}

	return guess, nil
}

var (
	RomanNumeral = (func() Rule {
		rnr := RomanNumeralRule{
			id:         RomanNumeralRuleId,
			RuleWeight: Weights(RomanNumeralRuleWeight),
			ruleType:   PreValidateType,
		}

		registerRule(rnr)

		return rnr
	})()
)
