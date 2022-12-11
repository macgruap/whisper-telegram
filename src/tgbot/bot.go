package tgbot

import (
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/m1guelpf/chatgpt-telegram/src/markdown"
)

type Bot struct {
	Token        string
	Username     string
	api          *tgbotapi.BotAPI
	editInterval time.Duration
}

func New(token string, editInterval time.Duration) (*Bot, error) {
	var api *tgbotapi.BotAPI
	var err error
	apiEndpoint, exist := os.LookupEnv("TELEGRAM_API_ENDPOINT")
	if exist && apiEndpoint != "" {
		api, err = tgbotapi.NewBotAPIWithAPIEndpoint(token, apiEndpoint)
	} else {
		api, err = tgbotapi.NewBotAPI(token)
	}
	if err != nil {
		return nil, err
	}

	return &Bot{
		Token:        token,
		Username:     api.Self.UserName,
		api:          api,
		editInterval: editInterval,
	}, nil
}

func (b *Bot) GetUpdatesChan() tgbotapi.UpdatesChannel {
	cfg := tgbotapi.NewUpdate(0)
	cfg.Timeout = 30
	return b.api.GetUpdatesChan(cfg)
}

func (b *Bot) Stop() {
	b.api.StopReceivingUpdates()
}

func (b *Bot) Send(chatID int64, replyTo int, text string) (tgbotapi.Message, error) {
	text = markdown.EnsureFormatting(text)
	msg := tgbotapi.NewMessage(chatID, text)
	if replyTo != 0 {
		msg.ReplyToMessageID = replyTo
	}
	return b.api.Send(msg)
}

func (b *Bot) SendEdit(chatID int64, messageID int, text string) error {
	text = markdown.EnsureFormatting(text)
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	msg.ParseMode = "Markdown"
	if _, err := b.api.Send(msg); err != nil {
		if err.Error() == "Bad Request: message is not modified: specified new message content and reply markup are exactly the same as a current content and reply markup of the message" {
			return nil
		}
		return err
	}
	return nil
}

func (b *Bot) SendTyping(chatID int64) {
	if _, err := b.api.Request(tgbotapi.NewChatAction(chatID, "typing")); err != nil {
		log.Printf("Couldn't send typing action: %v", err)
	}
}

func (b *Bot) GetFileDirectURL(fileID string) (string, error) {
	file, err := b.api.GetFile(tgbotapi.FileConfig{fileID})

	if err != nil {
		return "", err
	}

	return file.Link(b.Token), nil
}
