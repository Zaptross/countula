package commands

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	dhv "github.com/zaptross/godohuver"
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

	version := string(dat)
	versionMessage := fmt.Sprintf("Current: %s", version)

	semver, err := dhv.ExtractSemver(version)
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
