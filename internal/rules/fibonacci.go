package rules

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/emoji"
	"gorm.io/gorm"
)

type FibonacciRule struct {
	RuleWeight
	id       int
	ruleType string
}

func (fr FibonacciRule) Id() int {
	return fr.id
}
func (fr FibonacciRule) Name() string {
	return "Fibonacci's Sequence"
}
func (fr FibonacciRule) Description() string {
	return "Count up by adding the previous two numbers together"
}
func (fr FibonacciRule) Weight() int {
	return fr.Current
}
func (fr FibonacciRule) WithWeight(weight int) Rule {
	return FibonacciRule{
		id:         fr.id,
		ruleType:   fr.ruleType,
		RuleWeight: SetupWeight(weight),
	}
}
func (fr FibonacciRule) Type() string {
	return fr.ruleType
}
func (fr FibonacciRule) OnNewGame(db *gorm.DB, s *discordgo.Session, ng database.Turn, channelID string) {
	fibonacciTurn := database.Turn{
		ChannelID: channelID,
		UserID:    ng.UserID,
		Game:      ng.Game,
		Rules:     ng.Rules,
		Turn:      ng.Turn + 1,
		Guess:     1,
		Correct:   true,
	}

	zeroMessage := "0"
	if fibonacciTurn.Rules&RomanNumeralRuleId == RomanNumeralRuleId {
		zeroMessage = "nulla"
	}

	msg, err := s.ChannelMessageSend(channelID, zeroMessage)
	if err != nil {
		panic("Could not create new game: " + err.Error())
	}
	go s.MessageReactionAdd(channelID, msg.ID, emoji.CHECK)

	oneMessage := "1"
	if fibonacciTurn.Rules&RomanNumeralRuleId == RomanNumeralRuleId {
		oneMessage = "I"
	}

	msg, err = s.ChannelMessageSend(channelID, oneMessage)
	if err != nil {
		panic("Could not create new game: " + err.Error())
	}
	go s.MessageReactionAdd(channelID, msg.ID, emoji.CHECK)

	fibonacciTurn.MessageID = msg.ID
	db.Create(&fibonacciTurn)
}

func (fr FibonacciRule) Validate(db *gorm.DB, lastTurn database.Turn, msg discordgo.Message, guess int) bool {
	var secondLastTurn database.Turn
	db.Where("game = ? AND turn = ? AND channel_id = ?", lastTurn.Game, lastTurn.Turn-1, msg.ChannelID).First(&secondLastTurn)
	return secondLastTurn.Guess+lastTurn.Guess == guess
}

var (
	Fibonacci = (func() ValidateRule {
		fr := FibonacciRule{
			id:         FibonacciRuleId,
			RuleWeight: SetupWeight(FibonacciRuleWeight),
			ruleType:   CountType,
		}

		registerRule(fr)

		return fr
	})()
)
