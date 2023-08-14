package rules

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
)

type IncrementRule struct {
	increment int
	id        int
	weight    int
	ruleType  string
}

var (
	IncrementOne   = createIncrementRule(1, IncrementRule1Id, 30)
	IncrementTwo   = createIncrementRule(2, IncrementRule2Id, 30)
	IncrementThree = createIncrementRule(3, IncrementRule3Id, 20)
	IncrementSeven = createIncrementRule(7, IncrementRule7Id, 20)
)

func createIncrementRule(increment int, id int, weight int) Rule {
	r := IncrementRule{
		increment: increment,
		id:        id,
		weight:    weight,
		ruleType:  CountType,
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
	return ir.weight
}
func (ir IncrementRule) Type() string {
	return ir.ruleType
}

func (ir IncrementRule) Validate(db *gorm.DB, lastTurn database.Turn, msg discordgo.Message, guess int) bool {
	return lastTurn.Guess+ir.increment == guess
}
