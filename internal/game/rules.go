package game

import (
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/rules"
	"gorm.io/gorm"
)

func getRulesForNewGame(db *gorm.DB, guildID string) int {
	ruleSettings := database.GetRuleSettingsForGuild(db, guildID)
	allRules := rules.AllRules

	allRules = rules.ApplyWeightsToRules(allRules, ruleSettings)

	pv := rules.GetRandomPreValidateRule(allRules)
	c := rules.GetRandomCountRule(allRules)
	v := rules.GetRandomValidateRule(allRules)

	return pv.Id() | c.Id() | v.Id()
}
