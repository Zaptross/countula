package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/verbeage"
	"gorm.io/gorm"
)

type HelpCommand struct{}

const (
	HelpCommandName = "!help"
)

func (c HelpCommand) Execute(db *gorm.DB, s *discordgo.Session, m *discordgo.MessageCreate) {
	rp := verbeage.GetRandomHelpMessage()

	gm, err := s.GuildMember(m.GuildID, m.Author.ID)

	if err != nil {
		println(err.Error())
		return
	}

	reply, err := rp.Reply(verbeage.TemplateFields{
		Username: gm.Nick,
	})

	if err != nil {
		println(err.Error())
		return
	}

	_, err = s.ChannelMessageSendReply(m.ChannelID, reply, m.Reference())

	if err != nil {
		println(err.Error())
		return
	}
}

func (hc HelpCommand) Describe() string {
	return "I will offer my wisdom."
}
