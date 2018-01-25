package telegram

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"

	"os"
)

var (
	bot *tgbotapi.BotAPI
	chatId int64
	err error
)

func Authorization() {
	botToken := os.Getenv("telegram_bot_token")
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
}

func Listen(handler func(command string, text string)) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, _ := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		chatId = update.Message.Chat.ID

		handler(update.Message.Command(), update.Message.Text)
	}
}

func Send(text string) {
	msg := tgbotapi.NewMessage(chatId, text)
	bot.Send(msg)
}


