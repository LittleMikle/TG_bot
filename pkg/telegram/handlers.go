package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const (
	commandStart       = "start"
	replyStartTemplate = "Привет! Для сохранения ссылок в своем Pocket необходимо авторизоваться. Для этого переходи по ссылке:\n%s"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)

	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	b.bot.Send(msg)
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	authLink, err := b.generateAuthorizationLink(message.Chat.ID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(replyStartTemplate, authLink))
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я НЕ ЗНАЮ ТАКОЙ КОМАНДЫ :(")
	_, err := b.bot.Send(msg)
	return err
}
