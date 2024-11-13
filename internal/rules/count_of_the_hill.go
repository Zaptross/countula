package rules

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/emoji"
	"gorm.io/gorm"
)

type CountOfTheHillRule struct {
	RuleWeight
	id       int
	ruleType string
}

func (cothr CountOfTheHillRule) Id() int {
	return cothr.id
}
func (cothr CountOfTheHillRule) Name() string {
	return "Count of the Hill"
}
func (cothr CountOfTheHillRule) Description() string {
	return "To become Count of the Hill, you must guess correctly after the previous player by counting up in 2's.\nWhile you are Count of the Hill, you must guess correctly after yourself by counting up in 4's."
}
func (cothr CountOfTheHillRule) Weight() int {
	return cothr.Current
}
func (cothr CountOfTheHillRule) WithWeight(weight int) Rule {
	return CountOfTheHillRule{
		id:         cothr.id,
		ruleType:   cothr.ruleType,
		RuleWeight: SetupWeight(weight),
	}
}
func (cothr CountOfTheHillRule) Type() string {
	return cothr.ruleType
}
func (cothr CountOfTheHillRule) OnNewGame(db *gorm.DB, s *discordgo.Session, ng database.Turn, channelID string) {
}
func (cothr CountOfTheHillRule) OnFailure(fc *FailureContext) *FailureContext {
	var gameTurns []database.Turn
	fc.DB.Find(&gameTurns, "game = ?", fc.LastTurn.Game)

	hillTurns := lo.Reduce(gameTurns, func(acc map[string]int, t database.Turn, i int) map[string]int {
		// If the current turn is correct and the previous turn was by the same user
		// then increment the count for that user.
		if i > 0 && t.Correct && gameTurns[i-1].UserID == t.UserID {
			acc[t.UserID]++
		}
		return acc
	}, map[string]int{})

	hillUser := ""
	hillPoints := 0

	for user, points := range hillTurns {
		if points > hillPoints {
			hillUser = user
			hillPoints = points
		}
	}

	if hillUser != "" {
		go fc.DG.ChannelMessageSend(fc.LastTurn.ChannelID, fmt.Sprintf("The Count of the Hill was <@%s> with %d hill points %s", hillUser, hillPoints, emoji.FIST))
	}

	return fc
}

func (cothr CountOfTheHillRule) Validate(db *gorm.DB, lastTurn database.Turn, msg discordgo.Message, guess int) bool {
	if lastTurn.UserID == msg.Author.ID {
		return guess == lastTurn.Guess+4
	}
	return guess == lastTurn.Guess+2
}

var (
	CountOfTheHill = (func() ValidateRule {
		cothr := CountOfTheHillRule{
			id:         CountOfTheHillRuleId,
			RuleWeight: SetupWeight(CountOfTheHillRuleWeight),
			ruleType:   CountType,
		}

		registerRule(cothr)

		return cothr
	})()
)
