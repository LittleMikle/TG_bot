package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
	"net/url"
)

const (
	commandStart       = "start"
	replyStartTemplate = "Привет! Для сохранения ссылок в своем Pocket необходимо авторизоваться. Для этого переходи по ссылке:\n%s"
	replyAlreadyAuth   = "Ты уже авторизирован, присылай ссылку, я сохраню!"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)

	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) initAuthProcess(message *tgbotapi.Message) error {
	authLink, err := b.generateAuthorizationLink(message.Chat.ID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(replyStartTemplate, authLink))
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Ссылка успешно сохранена!")

	_, err := url.ParseRequestURI(message.Text)
	if err != nil {
		msg.Text = "Плохая ссылка!"
		_, err = b.bot.Send(msg)
		return err
	}

	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		msg.Text = "Ты не авторизирован! Используй комманду /start"
		_, err = b.bot.Send(msg)
		return err
	}

	err = b.pocketClient.Add(context.Background(), pocket.AddInput{
		AccessToken: accessToken,
		URL:         message.Text,
	})
	if err != nil {
		msg.Text = "Не удалось сохранить ссылку..."
		_, err = b.bot.Send(msg)
		return err
	}

	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAuthProcess(message)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, replyAlreadyAuth)
	_, err = b.bot.Send(msg)
	if err != nil {
		return err
	}
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я НЕ ЗНАЮ ТАКОЙ КОМАНДЫ :(")
	_, err := b.bot.Send(msg)
	return err
}
