package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/rules"
	"github.com/zaptross/countula/internal/verbeage"
	"gorm.io/gorm"
)

type RulesCommand struct{}

const (
	RulesCommandName = "!rules"
)

func (c RulesCommand) Execute(db *gorm.DB, s *discordgo.Session, m *discordgo.MessageCreate) {
	turn := database.GetCurrentTurn(db)

	ruleMessage := verbeage.GetRandomRuleMessage()

	rm, err := ruleMessage.Message(verbeage.TemplateFields{})

	if err != nil {
		panic("Could not create new game: " + err.Error())
	}

	rules := rules.GetRuleTextsForGame(turn)

	_, err = s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("%s\n%s", rm, strings.Join(rules, "\n")), m.Message.Reference())
	if err != nil {
		panic("Could not create new game: " + err.Error())
	}
}

func (c RulesCommand) Describe() string {
	return "I shall reiterate the rules."
}
