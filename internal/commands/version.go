package commands

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type VersionCommand struct {
	version string
}

const (
	VersionCommandName = "!version"
)

func (vc VersionCommand) Execute(db *gorm.DB, session *discordgo.Session, message *discordgo.MessageCreate) {
	if vc.version != "" {
		session.ChannelMessageSend(message.ChannelID, vc.version)
		return
	}

	dat, err := os.ReadFile("/etc/program-version")

	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Could not read version file")
		return
	}

	session.ChannelMessageSend(message.ChannelID, string(dat))
}

func (vc VersionCommand) Describe() string {
	return "If you must know, I can check my version."
}
