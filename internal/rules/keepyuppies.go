package rules

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
)

type KeepyUppiesRule struct {
	RuleWeight
	id       int
	ruleType string
}

func (ku KeepyUppiesRule) Id() int {
	return ku.id
}
func (ku KeepyUppiesRule) Name() string {
	return "Keepy Uppies"
}
func (ku KeepyUppiesRule) Description() string {
	return "Every guess keeps the balloon up in the air for one more hour. If the balloon hits the ground, the game is over. My friend Bluey taught me this one.\n\n...You have 24 hours."
}
func (ku KeepyUppiesRule) Weight() int {
	return ku.Current
}
func (ku KeepyUppiesRule) WithWeight(weight int) Rule {
	return KeepyUppiesRule{
		id:         ku.id,
		ruleType:   ku.ruleType,
		RuleWeight: SetupWeight(weight),
	}
}
func (ku KeepyUppiesRule) Type() string {
	return ku.ruleType
}
func (ku KeepyUppiesRule) OnNewGame(db *gorm.DB, s *discordgo.Session, ng database.Turn, channelID string) {
}
func (ku KeepyUppiesRule) OnFailure(fc *FailureContext) *FailureContext {
	return fc
}

func (ku KeepyUppiesRule) Validate(db *gorm.DB, lastTurn database.Turn, msg discordgo.Message, guess int) bool {
	var game database.Turn
	db.First(&game, "game = ? AND turn = ?", lastTurn.Game, 0)
	return game.CreatedAt.Add(24*time.Hour + time.Duration(lastTurn.Turn*int(time.Hour.Nanoseconds()))).After(time.Now())
}

var (
	KeepyUppies = (func() ValidateRule {
		ku := KeepyUppiesRule{
			id:         KeepyUppiesRuleId,
			RuleWeight: SetupWeight(KeepyUppiesRuleWeight),
			ruleType:   ValidateType,
		}

		registerRule(ku)

		return ku
	})()
)
