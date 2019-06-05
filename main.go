package main

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
)

var data = `
folder: ###
token: ###
log: true
telegramlog: false
`

// BotConfig contains configuration for running filesaver.
type BotConfig struct {
	Folder string
	Token  string
	Log	bool
	Telegramlog bool
}

func main() {
	botConfig := BotConfig{}

	err := yaml.Unmarshal([]byte(data), &botConfig)
	if err != nil {
		log.Fatalf("ERROR: Can't read the config file. Details: %v", err)
		return
	}

	bot, err := tgbotapi.NewBotAPI(botConfig.Token)
	if err != nil {
		log.Panicf("ERROR: Can't connect to telegram bot. Details: %v", err)
		return
	}

	storage := Storage{Config{botConfig.Folder}}

	bot.Debug = botConfig.Telegramlog

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		handleMessage(update, bot, storage)
	}
}
