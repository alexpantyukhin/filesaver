package main

import (
	"log"
	"path/filepath"
	"strconv"
	"time"

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

func handleMessage(update tgbotapi.Update, botAPI *tgbotapi.BotAPI, storage Storage) {
	//userName := update.Message.From.UserName
	//text := update.Message.Text
	//chatID := update.Message.Chat.ID
	//messageID := update.Message.MessageID

	// _, err := botAPI.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: chatID, MessageID: messageID})
	fileID := ""
	fileName := ""

	var linkURL string
	var err error

	document := update.Message.Document
	photos := update.Message.Photo
	if document != nil {
		fileID = document.FileID
		fileName = document.FileName
		linkURL, err = botAPI.GetFileDirectURL(fileID)
		if err != nil {
			// Log
		}
	} else if (photos != nil) && len(*photos) > 0 {
		photo := (*photos)[len(*photos)-1]
		fileID = photo.FileID
		fileName = strconv.FormatInt(time.Now().Unix(), 10)

		linkURL, err := botAPI.GetFileDirectURL(fileID)
		if err != nil {
			// Log
		}

		ext := filepath.Ext(linkURL)
		fileName = fileName + ext
	}

	if fileID != "" {
		folders := storage.GetInnerFolders()
		if len(folders) == 0 {
			_, err := storage.DownloadFileIntoFolder(linkURL, fileName)
			if err != nil {
				// Log
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "saved!")
			_, err = botAPI.Send(msg)

			if err != nil {
				// Log
			}
		} else {
			// var keyboardButtons [][]tgbotapi.KeyboardButton

			// for _, folder := range folders {
			// 	keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButtonRow(tgbotapi.KeyboardButton{Text: folder, RequestContact: false, RequestLocation: false}))
			// }

			// markup := tgbotapi.NewReplyKeyboard(keyboardButtons...)
			var keyboardButtons [][]tgbotapi.InlineKeyboardButton

			for _, folder := range folders {
				//button := tgbotapi.InlineKeyboardButton{Text: folder, RequestContact: false, RequestLocation: false}
				button := tgbotapi.NewInlineKeyboardButtonData(folder, "")
				keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardRow(button))
			}

			markup := tgbotapi.NewInlineKeyboardMarkup(keyboardButtons...)
			msg := tgbotapi.NewEditMessageReplyMarkup(update.Message.Chat.ID, update.Message.MessageID, markup)
			_, err = botAPI.Send(markup)
		}
	}
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
		if update.Message == nil {
			continue
		}

		handleMessage(update, bot, storage)
	}
}
