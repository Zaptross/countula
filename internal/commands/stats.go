package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/statistics"
	"gorm.io/gorm"
)

type StatsCommand struct{}

const (
	StatsCommandName = "!stats"
)

func (c StatsCommand) Execute(db *gorm.DB, s *discordgo.Session, m *discordgo.MessageCreate) {
	go statistics.Display(db, s, m)
}
