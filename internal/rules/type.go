package rules

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type Rule interface {
	Id() int
	Name() string
	// Description returns a description of the rule
	Description() string

	// PreValidate returns the guess as an int and an error if the guess is invalid
	PreValidate(db *gorm.DB, dg *discordgo.Session, msg discordgo.Message) (int, error)

	// Validate returns true if the guess is correct
	Validate(db *gorm.DB, dg *discordgo.Session, msg discordgo.Message, guess int) bool
}
