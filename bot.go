package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramFile struct {
	fileID string
	fileName string
	linkURL string
}

func handleMessage(update tgbotapi.Update, botAPI *tgbotapi.BotAPI, storage Storage) {

	if update.CallbackQuery != nil {
		message := update.CallbackQuery.Message
		chat := message.Chat

		documentMessage := update.CallbackQuery.Message.ReplyToMessage
		targetFolder := update.CallbackQuery.Data
		telegramFile, err := getFileByMessage(documentMessage, botAPI)
		if err != nil {
			// Log
		}

		_, err = storage.DownloadFileIntoSubFolder(telegramFile.linkURL, telegramFile.fileName, targetFolder)
		if err != nil {
			// Log
		}

		text := "Saved into "+targetFolder+"!"
		msgedit := tgbotapi.NewEditMessageText(chat.ID, message.MessageID, text)

		_, err = botAPI.Send(msgedit)

		return
	}

	telegramFile, err := getFileByMessage(update.Message, botAPI)
	if err != nil {
		// Log
	}

	if telegramFile.fileID != "" {
		folders := storage.GetInnerFolders()
		if len(folders) == 0 {
			_, err := storage.DownloadFileIntoFolder(telegramFile.linkURL, telegramFile.fileName)
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

			_, err = botAPI.Send(msg)

			if err != nil {
				// Log
			}
		}
	}
}

func getFileByMessage(message *tgbotapi.Message, botAPI *tgbotapi.BotAPI) (telegramFile *TelegramFile, err error) {
	document := message.Document
	photos := message.Photo
	var fileID string
	var fileName string
	var linkURL string

	if document != nil {
		fileID = document.FileID
		fileName = document.FileName
		linkURL, err = botAPI.GetFileDirectURL(fileID)
		if err != nil {
			// Log

			return nil, err
		}
	} else if (photos != nil) && len(*photos) > 0 {
		photo := (*photos)[len(*photos)-1]
		fileID = photo.FileID
		fileName = strconv.FormatInt(time.Now().Unix(), 10)

		linkURL, err := botAPI.GetFileDirectURL(fileID)
		if err != nil {
			// Log

			return nil, err
		}

		ext := filepath.Ext(linkURL)
		fileName = fileName + ext
	}

	return &TelegramFile{fileID: fileID, fileName: fileName, linkURL: linkURL} , err
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