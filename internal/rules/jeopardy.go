package rules

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

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
	return "You **must** guess in the form: `What is X + Y?`, where X and Y are numbers that equate the next number."
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
	jeopardyRegex  = regexp.MustCompile(`[Ww]hat is ((?:-?\d+)(?: ?(?:[+]|plus|[-]|minus|[*]|times|[/]|divided by) ?(?:-?\d+))+)\??`)
	whatIsDetector = regexp.MustCompile(`[Ww]hat is (?:-?\d+)\??`)
)

func (jr JeopardyRule) PreValidate(db *gorm.DB, dg *discordgo.Session, msg discordgo.Message) (int, error) {
	match := jeopardyRegex.FindStringSubmatch(msg.Content)
	if len(match) < 2 {

		if whatIsDetector.MatchString(msg.Content) {
			dg.ChannelMessageSendReply(msg.ChannelID, "Just entering the number isn't good enough buddy.", msg.Reference())
		}

		return 0, errors.New("no match found")
	}

	// Replace words with operators to make it easier to parse, eg. "What is 1 plus 2?" -> "1 + 2"
	input := replaceMany(match[1], map[string]string{
		"plus":       "+",
		"minus":      "-",
		"times":      "*",
		"divided by": "/",
		" ":          "",
	})

	operatorCount := strings.Count(input, "+") + strings.Count(input, "-") + strings.Count(input, "*") + strings.Count(input, "/")

	if operatorCount < 1 || operatorCount > 3 {
		dg.ChannelMessageSendReply(msg.ChannelID, "I appreciate the attempt, but let's keep it brief shall we? Say... three steps at most.", msg.Reference())
		return 0, errors.New("invalid number of arithmetic operators, must be at least 1 and at most 3")
	}

	guess := 0
	action := '+'
	parts := []string{"0", "+"}

	current := ""
	for _, char := range input {
		if strings.ContainsRune("+-*/", char) {
			num, err := strconv.Atoi(current)

			if err != nil {
				return 0, err
			}

			guess = applyAction(action, guess, num)

			action = char
			parts = append(parts, current, string(char))
			current = ""
		} else {
			current += string(char)
		}
	}

	num, err := strconv.Atoi(current)
	if err != nil {
		return 0, err
	}

	parts = append(parts, current)
	guess = applyAction(action, guess, num)

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

func applyAction(action rune, a, b int) int {
	switch action {
	case '+':
		return a + b
	case '-':
		return a - b
	case '*':
		return a * b
	case '/':
		return a / b
	}

	return 0
}

func replaceMany(s string, replacements map[string]string) string {
	for k, v := range replacements {
		s = strings.ReplaceAll(s, k, v)
	}

	return s
}
