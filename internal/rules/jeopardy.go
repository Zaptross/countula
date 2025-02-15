package rules

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
)

type JeopardyRule struct {
	RuleWeight
	id       int
	ruleType string
}

func (jr JeopardyRule) Id() int {
	return jr.id
}
func (jr JeopardyRule) Name() string {
	return "Jeopardy"
}
func (jr JeopardyRule) Description() string {
	return "You **must** guess in the form: `What is X + Y?`, where X and Y are numbers that add up to the next number."
}
func (jr JeopardyRule) Weight() int {
	return jr.Current
}
func (jr JeopardyRule) WithWeight(weight int) Rule {
	return JeopardyRule{
		id:         jr.id,
		ruleType:   jr.ruleType,
		RuleWeight: SetupWeight(weight),
	}
}
func (jr JeopardyRule) Type() string {
	return jr.ruleType
}
func (jr JeopardyRule) OnNewGame(_ *gorm.DB, _ *discordgo.Session, _ database.Turn, _ string) {}
func (jr JeopardyRule) OnFailure(fc *FailureContext) *FailureContext {
	return fc
}

var (
	jeopardyRegex = regexp.MustCompile(`[Ww]hat is (\d+) (?:\+|plus) (\d+)\??`)
)

func (jr JeopardyRule) PreValidate(db *gorm.DB, dg *discordgo.Session, msg discordgo.Message) (int, error) {
	matches := jeopardyRegex.FindStringSubmatch(msg.Content)

	if len(matches) == 0 {
		return 0, errors.New("no matches found")
	}

	guess := 0
	for _, match := range matches[1:] {
		num, err := strconv.Atoi(match)
		if err != nil {
			return 0, err
		}
		guess += num
	}

	return guess, nil
}

var (
	Jeopardy = (func() Rule {
		jr := JeopardyRule{
			id:         JeopardyRuleId,
			RuleWeight: SetupWeight(JeopardyRuleWeight),
			ruleType:   PreValidateType,
		}

		registerRule(jr)

		return jr
	})()
)
