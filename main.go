package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	statuses := make(map[int64]string)
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env: %s", err.Error())
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand /sayhi and /status."
		case "sayhi":
			msg.Text = "Hi :)"
		case "status":
			msg.Text = "I'm ok."
		case "start":
			msg.Text = "Привет! Этот бот создан для того, чтобы опавезать тебя о новых ответах на твои вопросы и комментариях к статьям, чтобы прикрепить свой аккаунт к боту, напиши команду /reg"
		case "reg":
			msg.Text = "Введите вашу почту и пароль через пробел:"
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
			statuses[update.Message.From.ID] = "wait for reg"
			logs(update.Message.From.UserName, update.Message.Text, msg.Text)

			continue
		case "cancel":
			msg.Text = fmt.Sprintf("Действие %s отмененно", statuses[update.Message.From.ID])
			delete(statuses, update.Message.From.ID)
		default:
			msg.Text = "I don't know that command"
		}

		if statuses[update.Message.From.ID] == "wait for reg" {
			emailAndPassword := update.Message.Text
			emailAndPasswordParts := strings.Split(emailAndPassword, " ")

			fmt.Println(emailAndPassword)

			if len(emailAndPasswordParts) == 2 {
				email := emailAndPasswordParts[0]
				password := emailAndPasswordParts[1]
				// TODO: работа с почтой и паролем
				log.Printf("[%s] Reg: %s %s", update.Message.From.UserName, email, password)
				msg.Text = fmt.Sprintf("Пользователь прикреплен email: %s password: %s", email, password)
				delete(statuses, update.Message.From.ID)
			} else {
				msg.Text = "Некорректные данные"
			}
		}
		logs(update.Message.From.UserName, update.Message.Text, msg.Text)
		if _, err := bot.Send(msg); err != nil {
			log.Fatal(err)
		}
	}
}

func logs(UserName, Text, TextU string) {
	log.Printf("[%s] (text: %s) %s", UserName, Text, TextU)
}

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
