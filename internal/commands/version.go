package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/utils"
	dhv "github.com/zaptross/godohuver"
	"gorm.io/gorm"
)

type VersionCommand struct{}

const (
	VersionCommandName = "!version"
)

func (vc VersionCommand) Execute(db *gorm.DB, session *discordgo.Session, message *discordgo.MessageCreate) {
	versionMessage := fmt.Sprintf("Current: %s", utils.GetVersion())

	semver, err := dhv.ExtractSemver(utils.GetVersion())
	if err == nil {
		latest, err := dhv.GetLatestImage("zaptross/countula")

		if err == nil {
			if semver != latest.Tag {
				versionMessage = fmt.Sprintf("%s\nLatest: %s", versionMessage, latest.Tag)
			}
		}
	}

	session.ChannelMessageSend(message.ChannelID, versionMessage)
}

func (vc VersionCommand) Describe() string {
	return "If you must know, I can check my version."
}
