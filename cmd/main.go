package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/game"
	"github.com/zaptross/countula/internal/handler"
	"github.com/zaptross/countula/internal/verbeage"
)

type DiscordConfig struct {
	Token           string
	AdminRoleId     string
	CountingChannel string
}

func main() {
	var dbConfig database.DatabaseConfig
	envconfig.Process("database", &dbConfig)

	db := database.Connect(dbConfig)

	var botConfig DiscordConfig
	envconfig.Process("discord", &botConfig)

	dg, err := discordgo.New("Bot " + botConfig.Token)
	if err != nil {
		panic(err)
	}

	dg.AddHandler(handler.GetMessageHandler(db, handler.Config{AdminRoleId: botConfig.AdminRoleId, CountingChannel: botConfig.CountingChannel}))

	err = dg.Open()
	if err != nil {
		panic(err)
	}

	t, err := verbeage.GetRandomAwaken().Message(verbeage.TemplateFields{})

	if err != nil {
		t = ":eyes:"
	}

	dg.ChannelMessageSend(botConfig.CountingChannel, t)

	// if there's no game in progress, start a new one
	if database.GetCurrentTurn(db).Game == 0 {
		game.CreateNewGame(db, dg, botConfig.CountingChannel)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
