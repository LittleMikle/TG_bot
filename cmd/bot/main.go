package main

import (
	"github.com/LittleMikle/TG_bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("TELEGRAM TOKEN")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient("POKET TOKEN")
	if err != nil {
		log.Fatal(err)
	}

	telegramBot := telegram.NewBot(bot, pocketClient, "LittleMikle")
	if err := telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}
