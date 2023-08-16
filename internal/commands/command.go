package commands

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type Command interface {
	Execute(*gorm.DB, *discordgo.Session, *discordgo.MessageCreate)
	Describe() string
}

var allCommands = map[string]Command{
	HelpCommandName:    HelpCommand{},
	ListCommandName:    ListCommand{},
	RulesCommandName:   RulesCommand{},
	StateCommandName:   StateCommand{},
	StatsCommandName:   StatsCommand{},
	VersionCommandName: VersionCommand{},
}

func GetCommand(commandName string) Command {
	return allCommands[commandName]
}
