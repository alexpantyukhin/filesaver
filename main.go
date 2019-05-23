package main

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func handleMessage(update tgbotapi.Update, botAPI *tgbotapi.BotAPI, storage Storage) {
	//userName := update.Message.From.UserName
	//text := update.Message.Text
	//chatID := update.Message.Chat.ID
	//messageID := update.Message.MessageID

	// _, err := botAPI.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: chatID, MessageID: messageID})

	if (update.Message.Document != nil) {
		if len(storage.GetInnerFolders()) == 0 {
			linkURL, err := botAPI.GetFileDirectURL(update.Message.Document.FileID)
			if (err != nil){
				// Log
			}

			err = storage.DownloadFileIntoFolder(linkURL, update.Message.Document.FileName)
			if (err != nil){
				// Log
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			_, err = botAPI.Send(msg)			

			if err != nil {
				// Log
			}
		}
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI("###")
	storage := Storage {Config {""}}
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		handleMessage(update, bot, storage)
	}
}
