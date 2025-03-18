package rules

import (
	"fmt"
	"math/rand/v2"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/emoji"
	"gorm.io/gorm"
)

type DoublingSeasonRule struct {
	RuleWeight
	id       int
	ruleType string
}

func (dsr DoublingSeasonRule) Id() int {
	return dsr.id
}
func (dsr DoublingSeasonRule) Name() string {
	return "Doubling Season"
}
func (dsr DoublingSeasonRule) Description() string {
	return "First there was one, then there was two, how many doubles can you do?"
}
func (dsr DoublingSeasonRule) Weight() int {
	return dsr.Current
}
func (dsr DoublingSeasonRule) WithWeight(weight int) Rule {
	return DoublingSeasonRule{
		id:         dsr.id,
		ruleType:   dsr.ruleType,
		RuleWeight: SetupWeight(weight),
	}
}
func (dsr DoublingSeasonRule) Type() string {
	return dsr.ruleType
}
func (dsr DoublingSeasonRule) OnNewGame(db *gorm.DB, s *discordgo.Session, ng database.Turn, channelID string) {
	doublingSeasonTurn := database.Turn{
		ChannelID: channelID,
		UserID:    ng.UserID,
		Game:      ng.Game,
		Rules:     ng.Rules,
		Turn:      ng.Turn + 1,
		Guess:     rand.Int() % 6, // pick a random starting number between 1 and 5 for variety
		Correct:   true,
	}

	correctEmoji := emoji.CHECK

	if ng.Rules&KeepyUppiesRuleId == KeepyUppiesRuleId {
		correctEmoji = emoji.BALLOON
	}

	msg, err := s.ChannelMessageSend(channelID, fmt.Sprintf("Let's start with say...\n\n%d", doublingSeasonTurn.Guess))
	if err != nil {
		panic("Could not create new game: " + err.Error())
	}
	go s.MessageReactionAdd(channelID, msg.ID, correctEmoji)

	doublingSeasonTurn.MessageID = msg.ID
	db.Create(&doublingSeasonTurn)
}
func (dsr DoublingSeasonRule) OnFailure(fc *FailureContext) *FailureContext { return fc }

func (dsr DoublingSeasonRule) Validate(db *gorm.DB, lastTurn database.Turn, msg discordgo.Message, guess int) bool {
	return guess == lastTurn.Guess*2
}

var (
	DoublingSeason = (func() ValidateRule {
		dsr := DoublingSeasonRule{
			id:         DoublingSeasonRuleId,
			RuleWeight: SetupWeight(DoublingSeasonRuleWeight),
			ruleType:   CountType,
		}

		registerRule(dsr)

		return dsr
	})()
)
