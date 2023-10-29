package database

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type ServerConfig struct {
	GuildID           string `gorm:"primaryKey"`
	CountingChannelID string `gorm:"primaryKey"`
	CreatedAt         time.Time
}

func GetAllServerConfigs(db *gorm.DB) []ServerConfig {
	var configs []ServerConfig
	db.Find(&configs)
	return configs
}

func ConfigureFromMessage(db *gorm.DB, m *discordgo.MessageCreate) {
	db.Create(&ServerConfig{
		GuildID:           m.GuildID,
		CountingChannelID: m.ChannelID,
	})
}
