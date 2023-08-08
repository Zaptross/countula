package rules

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
)

type IncrementRule struct {
	increment int
	id        int
}

var (
	IncrementOne   = createIncrementRule(1, IncrementRule1Id)
	IncrementTwo   = createIncrementRule(2, IncrementRule2Id)
	IncrementThree = createIncrementRule(3, IncrementRule3Id)
	IncrementSeven = createIncrementRule(7, IncrementRule7Id)
)

func createIncrementRule(i int, id int) Rule {
	return IncrementRule{
		increment: i,
	}
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

func (ir IncrementRule) PreValidate(db *gorm.DB, dg *discordgo.Session, msg discordgo.Message) (int, error) {
	guess, err := strconv.ParseInt(msg.Content, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(guess), nil
}

func (ir IncrementRule) Validate(db *gorm.DB, dg *discordgo.Session, msg discordgo.Message, guess int) bool {
	var lastTurn database.Turn
	db.Last(&lastTurn)
	return lastTurn.Guess+ir.increment == guess
}
