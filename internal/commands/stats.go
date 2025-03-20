package commands

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type StatsCommand struct{}

const (
	StatsCommandName = "!stats"
)

func (c StatsCommand) Execute(db *gorm.DB, s *discordgo.Session, m *discordgo.MessageCreate) {
	go s.ChannelMessageSendReply(m.ChannelID, "Try the new slash command `/count stats`", m.Message.Reference())
}

func (c StatsCommand) Describe() string {
	return "I shall recount your... attempts."
}
