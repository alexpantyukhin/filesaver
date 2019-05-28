package main

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
)

var data = `
folder: ###
token: ###
`

type BotConfig struct {
	Folder string
	Token  string
}

func main() {
	botConfig := BotConfig{}

	err := yaml.Unmarshal([]byte(data), &botConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(botConfig.Token)
	storage := Storage{Config{botConfig.Folder}}
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		handleMessage(update, bot, storage)
	}
}
