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
folder: c:\\tmp\\folder
`

type BotConfig struct {
	Folder string
}

func handleMessage(update tgbotapi.Update, botAPI *tgbotapi.BotAPI, storage Storage) {
	//userName := update.Message.From.UserName
	//text := update.Message.Text
	//chatID := update.Message.Chat.ID
	//messageID := update.Message.MessageID

	// _, err := botAPI.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: chatID, MessageID: messageID})
	fileID := ""
	fileName := ""
	document := update.Message.Document
	photos := update.Message.Photo
	detectImageType := false
	if document != nil {
		fileID = document.FileID
		fileName = document.FileName
	} else if (photos != nil) && len(*photos) > 0 {
		photo := (*photos)[len(*photos)-1]
		fileID = photo.FileID
		fileName = strconv.FormatInt(time.Now().Unix(), 10)
		detectImageType = true
	}

	if fileID != "" {
		if len(storage.GetInnerFolders()) == 0 {
			linkURL, err := botAPI.GetFileDirectURL(fileID)
			if err != nil {
				// Log
			}

			//fmt.Println(linkURL)
			if detectImageType {
				ext := filepath.Ext(linkURL)
				fileName = fileName + ext
			}

			_, err = storage.DownloadFileIntoFolder(linkURL, fileName)
			if err != nil {
				// Log
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "saved!")
			_, err = botAPI.Send(msg)

			if err != nil {
				// Log
			}
		}
	}
}

func main() {
	botConfig := BotConfig{}

	err := yaml.Unmarshal([]byte(data), &botConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI("###")
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
