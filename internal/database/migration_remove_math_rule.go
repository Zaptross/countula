package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const mathRuleId = 8192
const guessNormallyRuleId = 32

type removeMathRule struct{}

// compile time check to ensure removeMathRule implements Migration
var _ Migration = (*removeMathRule)(nil)

func (m *removeMathRule) ID() int {
	return 20251202
}

func (m *removeMathRule) Up(db *gorm.DB) error {
	// update all RuleSettings with math rule to have weight 0
	var ruleSetting RuleSetting
	db.Model(&ruleSetting).Where("rule_id = ?", mathRuleId).Update("weight", 0)

	// update all current games to remove math rule, replacing it with Guess Normally rule
	var currentTurns []Turn
	db.Distinct("ON (channel_id) *").Where("rules & ? = ?", mathRuleId, mathRuleId).Order("channel_id, game desc, turn desc").Find(&currentTurns)

	for _, turn := range currentTurns {
		turn.Rules = (turn.Rules &^ mathRuleId) | guessNormallyRuleId
		db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "game"}, {Name: "turn"}},
			DoUpdates: clause.AssignmentColumns([]string{"rules"}),
		}).Save(&turn)
	}

	return nil
}

func (m *removeMathRule) Down(db *gorm.DB) error {
	// no-op
	return nil
}
