package handler

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/commands"
	"gorm.io/gorm"
)

func handleCommand(db *gorm.DB, s *discordgo.Session, m *discordgo.MessageCreate) {
	components := strings.Split(m.Content, " ")

	command := commands.GetCommand(components[0])

	if command != nil {
		command.Execute(db, s, m)
	}
}
