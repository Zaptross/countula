package main

import (
	"log"
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
	AppID string
}

func main() {
	var dbConfig database.DatabaseConfig
	envconfig.Process("database", &dbConfig)

	db := database.Connect(dbConfig)
	log.Default().Println("Connected to database")

	var botConfig DiscordConfig
	envconfig.Process("discord", &botConfig)

	dg, err := discordgo.New("Bot " + botConfig.Token)
	if err != nil {
		panic(err)
	}
	log.Default().Println("Connected to Discord")

	dg.AddHandler(handler.GetMessageHandler(db))
	log.Default().Println("Added message handler")

	dg.AddHandler(handler.GetOnInteractionHandler(db))
	_, err = dg.ApplicationCommandCreate(botConfig.AppID, "", handler.GetSlashCommand())

	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	log.Default().Println("Added slash commands handler")

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

		log.Default().Println("Sent awaken messages")
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
