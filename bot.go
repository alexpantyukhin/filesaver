package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const buttonsPerRow = 3

type telegramFile struct {
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
			log.Panicf("ERROR: Can't get file from message. Details: %v", err)
			return
		}

		_, err = storage.DownloadFileIntoSubFolder(telegramFile.linkURL, telegramFile.fileName, targetFolder)
		if err != nil {
			log.Panicf("ERROR: Can't download a file from message. Details: %v", err)
			return
		}

		text := "Saved into "+targetFolder+"!"
		msgedit := tgbotapi.NewEditMessageText(chat.ID, message.MessageID, text)

		_, err = botAPI.Send(msgedit)
		if err != nil {
			log.Panicf("ERROR: Can't edit a message. Details: %v", err)
			return
		}

		return
	}

	telegramFile, err := getFileByMessage(update.Message, botAPI)
	if err != nil {
		log.Panicf("ERROR: Can't get file from message. Details: %v", err)
		return
	}

	if telegramFile.fileID != "" {
		folders := storage.GetInnerFolders()
		if len(folders) == 0 {
			_, err := storage.DownloadFileIntoFolder(telegramFile.linkURL, telegramFile.fileName)
			if err != nil {
				log.Panicf("ERROR: Can't download a file from message. Details: %v", err)
				return
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "saved!")
			_, err = botAPI.Send(msg)

			if err != nil {
				log.Panicf("ERROR: Can't send a message. Details: %v", err)
				return
			}
		} else {
			makrup := getKeyboardFromNames(folders, buttonsPerRow)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Select folder.")
			msg.ReplyMarkup = makrup
			msg.ReplyToMessageID = update.Message.MessageID

			_, err = botAPI.Send(msg)

			if err != nil {
				log.Panicf("ERROR: Can't send a message. Details: %v", err)
				return
			}
		}
	}
}

func getFileByMessage(message *tgbotapi.Message, botAPI *tgbotapi.BotAPI) (*telegramFile, error) {
	document := message.Document
	photos := message.Photo
	var fileID string
	var fileName string
	var linkURL string

	if document != nil {
		fileID = document.FileID
		fileName = document.FileName
		var err error
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

	return &telegramFile{fileID: fileID, fileName: fileName, linkURL: linkURL} , nil
}

func getKeyboardFromNames(names []string, btnsPerRows int) tgbotapi.InlineKeyboardMarkup {
	var buttonsRows [][]tgbotapi.InlineKeyboardButton

	var buttonsRowsBuffer []tgbotapi.InlineKeyboardButton

	for _, name := range names {
		buttonsRowsBuffer = append(buttonsRowsBuffer, tgbotapi.NewInlineKeyboardButtonData(name, name))

		if len(buttonsRowsBuffer) >= btnsPerRows {
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