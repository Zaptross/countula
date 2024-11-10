package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RuleSetting struct {
	GuildID string `gorm:"primaryKey"`
	RuleID  int    `gorm:"primaryKey"`
	Weight  int
}

func GetRuleSettingsForGuild(db *gorm.DB, guildID string) []RuleSetting {
	var settings []RuleSetting
	db.Where("guild_id = ?", guildID).Find(&settings)
	return settings
}

func CreateRuleSettingForGuild(db *gorm.DB, guildID string, ruleID int, weight int) {
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "guild_id"}, {Name: "rule_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"weight"}),
	}).Create(&RuleSetting{
		GuildID: guildID,
		RuleID:  ruleID,
		Weight:  weight,
	})
}
