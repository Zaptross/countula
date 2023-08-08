package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"
	"github.com/zaptross/countula/internal/handler"
	"gorm.io/gorm"
	"gorm.io/gorm/driver/postgres"
)

type DiscordConfig struct {
	Token string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Ssl      string
	Timezone string
}

func main() {
	var dbConfig DatabaseConfig
	envconfig.Process("database", &dbConfig)

	db, err := gorm.Open(postgres.Open("host="+dbConfig.Host+" user="+dbConfig.User+" password="+dbConfig.Password+" dbname="+dbConfig.Database+" port="+dbConfig.Port+" sslmode="+dbConfig.Ssl+" TimeZone="+dbConfig.Timezone), &gorm.Config{})
	if err != nil {
		panic(err)
	}

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

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
