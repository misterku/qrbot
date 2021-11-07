package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	qrcode "github.com/skip2/go-qrcode"
)

func handleResponse(rw http.ResponseWriter, err error) {
	rw.WriteHeader(200)
	if err != nil {
		log.Printf("Error: %s", err)
	}
}

func PrepareBot() (*tgbotapi.BotAPI, error) {
	return tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
}

func ReadBody(req *http.Request) ([]byte, error) {
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	return b, err
}

func DeserializeRequest(req *http.Request) (*tgbotapi.Update, error) {
	b, err := ReadBody(req)
	if err != nil {
		return nil, err
	}
	var update tgbotapi.Update
	err = json.Unmarshal(b, &update)
	if err != nil {
		return nil, err
	}
	return &update, nil
}

var allowedSchemas = []string{"https", "http", "ftp"}

func extractURL(uri string) (string, error) {
	url, err := url.ParseRequestURI(uri)

	if url.Host == "" {
		return "", fmt.Errorf("host is empty")
	}

	validSchema := false
	for _, scheme := range allowedSchemas {
		if strings.EqualFold(scheme, url.Scheme) {
			validSchema = true
		}
	}

	if !validSchema {
		return "", fmt.Errorf("schema is not allowed")
	}

	if err != nil {
		return "", err
	} else {
		return uri, nil
	}
}

func createTempFileWithQrCode(uri string) (*os.File, error) {
	tempFile, err := ioutil.TempFile("/tmp", "qrbot")
	if err != nil {
		return nil, err
	}

	if err := qrcode.WriteFile(uri, qrcode.Medium, 128, tempFile.Name()); err != nil {
		return nil, err
	}
	return tempFile, nil
}

func handleCommand(msg messageContext, rw http.ResponseWriter, bot *tgbotapi.BotAPI) error {
	if msg.Text == "/start" || msg.Text == "/help" {
		return sendText(msg, "To make QR core, just send URL to the chat", bot)
	} else {
		return sendText(msg, "Unfortunately, this message is not valid command", bot)
	}
}

func sendText(msg messageContext, text string, bot *tgbotapi.BotAPI) error {
	m := tgbotapi.NewMessage(msg.ChatID, text)
	m.ReplyToMessageID = msg.MessageID
	_, err := bot.Send(m)
	return err
}

func sendPhoto(msg messageContext, f *os.File, bot *tgbotapi.BotAPI) error {
	m := tgbotapi.NewPhotoUpload(msg.ChatID, f.Name())
	m.ReplyToMessageID = msg.MessageID
	_, err := bot.Send(m)
	return err
}

type messageContext struct {
	MessageID int
	ChatID    int64
	Text      string
}

func NewMessageContext(update *tgbotapi.Update) (*messageContext, error) {
	if update.Message != nil && update.Message.Chat != nil {
		return &messageContext{
			MessageID: update.Message.MessageID,
			ChatID:    update.Message.Chat.ID,
			Text:      update.Message.Text,
		}, nil
	} else {
		return nil, fmt.Errorf("couldn't get all necessary info from update")
	}

}

func Handler(rw http.ResponseWriter, req *http.Request) {
	bot, err := PrepareBot()
	if err != nil {
		handleResponse(rw, err)
		return
	}

	update, err := DeserializeRequest(req)
	if err != nil {
		handleResponse(rw, err)
		return
	}

	msgContext, err := NewMessageContext(update)
	if err != nil {
		handleResponse(rw, err)
		return
	}

	if strings.HasPrefix(msgContext.Text, "/") {
		err := handleCommand(*msgContext, rw, bot)
		handleResponse(rw, err)
	} else {
		uri, err := extractURL(msgContext.Text)
		if err != nil {
			err := sendText(*msgContext, "Unfortunately, this message is not valid URL", bot)
			handleResponse(rw, err)
			return
		}

		tempFile, err := createTempFileWithQrCode(uri)
		defer os.Remove(tempFile.Name())
		if err != nil {
			handleResponse(rw, err)
			return
		}

		err = sendPhoto(*msgContext, tempFile, bot)
		handleResponse(rw, err)
	}
}
