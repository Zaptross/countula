package main

import (
	"log"
	"log/slog"
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
	Token          string
	AppID          string
	AdminChannelID string
}

func main() {
	var dbConfig database.DatabaseConfig
	envconfig.Process("database", &dbConfig)

	db := database.Connect(dbConfig)
	slog.Info("Connected to database")

	var botConfig DiscordConfig
	envconfig.Process("discord", &botConfig)

	dg, err := discordgo.New("Bot " + botConfig.Token)
	if err != nil {
		panic(err)
	}
	slog.Info("Connected to Discord")

	handler.StartupCheckMaintenanceMode(db)

	dg.AddHandler(handler.GetMessageHandler(db))
	slog.Info("Added message handler")

	dg.AddHandler(handler.GetOnInteractionHandler(db, botConfig.AdminChannelID))
	_, err = dg.ApplicationCommandCreate(botConfig.AppID, "", handler.GetSlashCommand())

	dg.AddHandler(handler.GetOnMessageDeletedHandler(db))

	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	slog.Info("Added slash commands handler")

	err = dg.Open()
	if err != nil {
		panic(err)
	}

	serverConfigs := database.GetAllServerConfigs(db)
	if handler.IsMaintenanceModeEnabled() {
		t, err := verbeage.GetRandomMaintenanceMessage().Reply(verbeage.TemplateFields{})

		if err != nil {
			t = ":eyes:"
		}

		for _, sc := range serverConfigs {
			dg.ChannelMessageSend(sc.CountingChannelID, t)
		}

		slog.Info("Maintenance mode is enabled")
	} else {
		if len(serverConfigs) != 0 {
			t, err := verbeage.GetRandomAwaken().Message(verbeage.TemplateFields{})

			if err != nil {
				t = ":eyes:"
			}

			for _, sc := range serverConfigs {
				dg.ChannelMessageSend(sc.CountingChannelID, t)
			}

			slog.Info("Sent awaken messages")
		}
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
