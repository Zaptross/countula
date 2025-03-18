package rules

import (
	"errors"
	"regexp"

	"github.com/Pramod-Devireddy/go-exprtk"
	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
)

type MathsRule struct {
	RuleWeight
	id       int
	ruleType string
}

var (
	mathsRegex = regexp.MustCompile(`^(?:[ ]|[0-9]+|[(,)+\-*/%^]|sin|cos|tan|acos|asin|atan|atan2|cosh|cot|csc|sec|sinh|tanh|d2r|r2d|d2g|g2d|hyp|min|max|avg|sum|abs|ceil|floor|round|roundn|exp|log|log10|logn|pow|root|sqrt|clamp)+$`)
)

func (mr MathsRule) Id() int {
	return mr.id
}
func (mr MathsRule) Name() string {
	return "Maths"
}
func (mr MathsRule) Description() string {
	return "You **must** guess in the form of a mathematical equation, which equals the next number after ignoring all decimal places. For example: `1.1111 * 9` which equals `9.9999` would be treated as `9`."
}
func (mr MathsRule) Weight() int {
	return mr.Current
}
func (mr MathsRule) WithWeight(weight int) Rule {
	return MathsRule{
		id:         mr.id,
		ruleType:   mr.ruleType,
		RuleWeight: SetupWeight(weight),
	}
}
func (mr MathsRule) Type() string {
	return mr.ruleType
}
func (mr MathsRule) OnNewGame(_ *gorm.DB, _ *discordgo.Session, _ database.Turn, _ string) {}
func (mr MathsRule) OnFailure(fc *FailureContext) *FailureContext {
	return fc
}

func (mr MathsRule) PreValidate(db *gorm.DB, dg *discordgo.Session, msg discordgo.Message) (int, error) {
	if !mathsRegex.MatchString(msg.Content) {
		// record attempts that don't match the regex
		go db.Create(&database.AuditLog{
			ChannelID: msg.ChannelID,
			MessageID: msg.ID,
			UserID:    msg.Author.ID,
			Username:  msg.Author.Username,
			Action:    "MathsRulePreValidate",
			Data:      msg.Content,
		})
		go dg.ChannelMessageSendReply(msg.ChannelID, "I appreciate the attempt, but that doesn't look like a valid mathematical expression to me.", msg.Reference())
		return 0, errors.New("forbidden mathematical expression")
	}

	expr := exprtk.NewExprtk()
	expr.SetExpression(msg.Content)

	err := expr.CompileExpression()
	if err != nil {
		return 0, err
	}

	result := expr.GetEvaluatedValue()

	return int(result), nil
}

var (
	Maths = (func() Rule {
		mr := MathsRule{
			id:         MathsRuleId,
			RuleWeight: SetupWeight(MathsRuleWeight),
			ruleType:   PreValidateType,
		}

		registerRule(mr)

		return mr
	})()
)
