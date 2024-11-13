package rules

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/verbeage"
	"gorm.io/gorm"
)

type FailureContext struct {
	DB             *gorm.DB
	DG             *discordgo.Session
	Msg            *discordgo.MessageCreate
	FailureMessage verbeage.ResponseParts
	LastTurn       database.Turn
	Guess          int
	ChannelID      string
	Emoji          string
}
