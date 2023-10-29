package handler

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/game"
	"github.com/zaptross/countula/internal/verbeage"
	"gorm.io/gorm"
)

const (
	ConfigureCommand = "!configure-countula"
)

func HandleConfigure(db *gorm.DB, s *discordgo.Session, m *discordgo.MessageCreate) {
	t, err := verbeage.GetRandomOnConfigureMessage().Message(verbeage.TemplateFields{})
	if err != nil {
		log.Default().Println("Could not get random on configure message", err)
		return
	}

	_, err = s.ChannelMessageSendReply(m.ChannelID, t, m.Message.Reference())

	if err != nil {
		log.Default().Println("Could not send on configure message", err)
		return
	}

	database.ConfigureFromMessage(db, m)
	game.CreateNewGame(db, s, m.ChannelID)
}
