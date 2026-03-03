package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"
	"github.com/zaptross/countula/internal/database"
	"github.com/zaptross/countula/internal/handler"
	"github.com/zaptross/countula/internal/utils"
)

type DiscordConfig struct {
	Token          string
	AppID          string
	AdminChannelID string
}

func main() {
	var logLevel slog.Level
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	slog.SetLogLoggerLevel(logLevel)
	slog.Info("configured logger", "level", logLevel)

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

	if botConfig.AdminChannelID != "" {
		msgs, err := dg.ChannelMessages(botConfig.AdminChannelID, 10, "", "", "")
		if err != nil {
			slog.Error("Failed to fetch messages from admin channel", "error", err)
		} else {
			latestMessageTimestamp := time.Time{}
			latestMessage := &discordgo.Message{}
			latestStartupVersion := ""

			// Find the most recent message sent by the bot in the admin channel
			for _, msg := range msgs {
				if msg.Timestamp.After(latestMessageTimestamp) && msg.Author.ID == dg.State.User.ID {
					latestMessageTimestamp = msg.Timestamp
					latestMessage = msg

					// Check if the message content contains a version number
					if strings.Contains(msg.Content, "version:") {
						parts := strings.Split(msg.Content, "version:")
						if len(parts) > 1 {
							latestStartupVersion = strings.TrimSpace(parts[1])
						}
					}
				}
			}

			// If the most recent message is older than 5 minutes or it's a different version, send a startup message
			if time.Now().After(latestMessageTimestamp.Add(time.Minute*5)) || latestStartupVersion != utils.GetVersion() {
				slog.Debug("alerting admin to start", "latestMessageID", latestMessage.ID, "latestMessageTimestamp", latestMessage.Timestamp, "latestStartupVersion", latestStartupVersion)
				_, err := dg.ChannelMessageSend(botConfig.AdminChannelID, fmt.Sprintf("Countula started on version: %s", utils.GetVersion()))
				if err != nil {
					slog.Error("Failed to send startup message to admin channel", "error", err)
				}
			} else {
				slog.Debug("skipping startup message, recent message found", "messageID", latestMessage.ID, "timestamp", latestMessage.Timestamp, "versionInMessage", latestStartupVersion)
			}
		}
	} else {
		slog.Info("Admin channel ID not set, skipping startup message")
	}

	slog.Info("countula started successfully", "version", utils.GetVersion())
	dg.UpdateCustomStatus(fmt.Sprintf("Version: %s", utils.GetVersion()))

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
