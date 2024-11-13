package game

import (
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/rules"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

func getRulesForNewGame(db *gorm.DB, guildID string) int {
	ruleSettings := database.GetRuleSettingsForGuild(db, guildID)
	allRules := rules.AllRules

	allRules = rules.ApplyWeightsToRules(allRules, ruleSettings)

	pv := rules.GetRandomPreValidateRule(allRules)
	c := rules.GetRandomCountRule(allRules)
	v := rules.GetRandomValidateRule(allRules)

	proposedRules := pv.Id() | c.Id() | v.Id()
	finalRules := rules.OverrideRuleSelections(proposedRules)

	if proposedRules != finalRules {
		slog.Info("Rules for new game altered by rule overrides", "proposed_rules", proposedRules, "final_rules", finalRules)
	}

	return finalRules
}
