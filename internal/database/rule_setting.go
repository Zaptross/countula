package database

import (
	"log"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type RuleSetting struct {
	GuildID string `gorm:"primaryKey"`
	RuleID  int    `gorm:"primaryKey"`
	Weight  int
}

type Rule interface {
	Id() int
	SetWeight(int)
}

func GetRuleSettingsForGuild(db *gorm.DB, guildID string) []RuleSetting {
	var settings []RuleSetting
	db.Where("guild_id = ?", guildID).Find(&settings)
	return settings
}

func ApplyWeightsToRules(rules []Rule, settings []RuleSetting) {
	lo.ForEach(settings, func(setting RuleSetting, _ int) {
		rule, ok := lo.Find(rules, func(w Rule) bool { return w.Id() == setting.RuleID })

		if !ok {
			log.Printf("Could not find rule for rule id: %d\n", setting.RuleID)
		}

		rule.SetWeight(setting.Weight)
	})
}
