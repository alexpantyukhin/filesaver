package main

import (
	"fmt"
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

	if update.CallbackQuery != nil {
		message := update.CallbackQuery.Message
		chat := message.Chat
		botAPI.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: chat.ID, MessageID: message.MessageID})

		documentMessage := update.CallbackQuery.Message.ReplyToMessage
		targetFolder := update.CallbackQuery.Data
		_, fileName, linkURL, err := getFileByMessage(documentMessage, botAPI)
		if err != nil {
			// Log
		}

		_, err = storage.DownloadFileIntoSubFolder(linkURL, fileName, targetFolder)
		if err != nil {
			// Log
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "saved into "+targetFolder+"!")
		_, err = botAPI.Send(msg)

		return
	}

	fileID, fileName, linkURL, err := getFileByMessage(update.Message, botAPI)
	if err != nil {
		// Log
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
			makrup := getKeyboardFromNames(folders)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Select folder.")
			msg.ReplyMarkup = makrup
			msg.ReplyToMessageID = update.Message.MessageID

			botAPI.Send(msg)
		}
	}
}

func getFileByMessage(message *tgbotapi.Message, botAPI *tgbotapi.BotAPI) (fileID string, fileName string, linkURL string, err error) {
	document := message.Document
	photos := message.Photo
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

	return fileID, fileName, linkURL, err
}

func getKeyboardFromNames(names []string) tgbotapi.InlineKeyboardMarkup {
	bottonsPerRows := 3
	var buttonsRows [][]tgbotapi.InlineKeyboardButton

	var buttonsRowsBuffer []tgbotapi.InlineKeyboardButton

	for _, name := range names {
		buttonsRowsBuffer = append(buttonsRowsBuffer, tgbotapi.NewInlineKeyboardButtonData(name, name))

		if len(buttonsRowsBuffer) >= bottonsPerRows {
			buttonsRows = append(buttonsRows, buttonsRowsBuffer)
			buttonsRowsBuffer = make([]tgbotapi.InlineKeyboardButton, 0)
		}

		fmt.Println(name)
	}

	if len(buttonsRowsBuffer) > 0 {
		buttonsRows = append(buttonsRows, buttonsRowsBuffer)
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttonsRows...)
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
		// if update.Message == nil {
		// 	continue
		// }

		handleMessage(update, bot, storage)
	}
}
