package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/handler"
	"github.com/zaptross/countula/internal/verbeage"
)

type DiscordConfig struct {
	Token string
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

	dg.AddHandler(handler.GetMessageHandler(db))

	err = dg.Open()
	if err != nil {
		panic(err)
	}

	serverConfigs := database.GetAllServerConfigs(db)
	if len(serverConfigs) != 0 {
		t, err := verbeage.GetRandomAwaken().Message(verbeage.TemplateFields{})

		if err != nil {
			t = ":eyes:"
		}

		for _, sc := range serverConfigs {
			dg.ChannelMessageSend(sc.CountingChannelID, t)
		}
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
