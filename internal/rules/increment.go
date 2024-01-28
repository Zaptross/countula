package rules

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
)

type IncrementRule struct {
	RuleWeight
	increment int
	id        int
	ruleType  string
}

var (
	IncrementOne   = createIncrementRule(1, IncrementRule1Id, IncrementRule1Weight)
	IncrementTwo   = createIncrementRule(2, IncrementRule2Id, IncrementRule2Weight)
	IncrementThree = createIncrementRule(3, IncrementRule3Id, IncrementRule3Weight)
	IncrementSeven = createIncrementRule(7, IncrementRule7Id, IncrementRule7Weight)
)

func createIncrementRule(increment int, id int, weight int) Rule {
	r := IncrementRule{
		increment:  increment,
		id:         id,
		RuleWeight: SetupWeight(weight),
		ruleType:   CountType,
	}

	registerRule(r)
	return r
}

func (ir IncrementRule) Id() int {
	return ir.id
}
func (ir IncrementRule) Name() string {
	return fmt.Sprintf("Increment by %d", ir.increment)
}
func (ir IncrementRule) Description() string {
	return fmt.Sprintf("Count up in increments of %d", ir.increment)
}
func (ir IncrementRule) Weight() int {
	return ir.Current
}
func (ir IncrementRule) WithWeight(weight int) Rule {
	return IncrementRule{
		id:         ir.id,
		increment:  ir.increment,
		ruleType:   ir.ruleType,
		RuleWeight: SetupWeight(weight),
	}
}
func (ir IncrementRule) Type() string {
	return ir.ruleType
}
func (ir IncrementRule) OnNewGame(_ *gorm.DB, _ *discordgo.Session, _ database.Turn, _ string) {}

func (ir IncrementRule) Validate(db *gorm.DB, lastTurn database.Turn, msg discordgo.Message, guess int) bool {
	return lastTurn.Guess+ir.increment == guess
}
