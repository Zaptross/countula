package database

import (
	"gorm.io/gorm"
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
