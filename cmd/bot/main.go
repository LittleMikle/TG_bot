package main

import (
	"github.com/LittleMikle/TG_bot/pkg/repository/boltdb"
	"github.com/LittleMikle/TG_bot/pkg/server"
	"github.com/LittleMikle/TG_bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient("")
	if err != nil {
		log.Fatal(err)
	}

	db, err := boltdb.ConnectToBoltDB()
	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := boltdb.NewTokenRepoBolt(db)

	telegramBot := telegram.NewBot(
		bot, pocketClient, tokenRepository, "http://localhost:8081/")

	authorizationServer := server.NewAuthorizationServer(
		pocketClient, tokenRepository, "https://t.me/littlemikle_pocket_api_bot")

	go func() {
		if err = telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	err = authorizationServer.Start()
	if err != nil {
		log.Fatal(err)
	}
}
