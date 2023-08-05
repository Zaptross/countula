package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"
)

type DiscordConfig struct {
	Token string
}

func main() {
	var config DiscordConfig
	envconfig.Process("discord", &config)

	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		panic(err)
	}

	// dg.handler

	err = dg.Open()
	if err != nil {
		panic(err)
	}

	<-make(chan struct{})
}
