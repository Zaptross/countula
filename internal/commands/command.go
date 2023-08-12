package commands

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type Command interface {
	Execute(*gorm.DB, *discordgo.Session, *discordgo.MessageCreate)
}

var allCommands = map[string]Command{
	RulesCommandName: RulesCommand{},
}

func GetCommand(commandName string) Command {
	return allCommands[commandName]
}