package handler

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/commands"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
)

func handleCommand(db *gorm.DB, s *discordgo.Session, m *discordgo.MessageCreate) {
	components := strings.Split(m.Content, " ")

	command := commands.GetCommand(components[0])

	if command != nil {
		command.Execute(db, s, m)

		commandData := ""

		if len(components) > 1 {
			commandData = strings.Join(components[1:], " ")
		}

		go db.Create(&database.AuditLog{
			UserID:    m.Author.ID,
			Username:  m.Author.Username,
			MessageID: m.ID,
			Action:    components[0],
			Data:      commandData,
		})
	}
}
