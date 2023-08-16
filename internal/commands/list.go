package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type ListCommand struct{}

const (
	ListCommandName = "!list"
)

func (c ListCommand) Execute(db *gorm.DB, s *discordgo.Session, m *discordgo.MessageCreate) {
	commands := allCommands

	reply := "Commands:\n"
	for name, command := range commands {
		reply += fmt.Sprintf("`%s`: %s\n", name, command.Describe())
	}

	s.ChannelMessageSendReply(m.ChannelID, reply, m.Reference())
}

func (c ListCommand) Describe() string {
	return "Don't make me tell you again."
}
