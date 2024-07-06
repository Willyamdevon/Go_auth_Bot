package main

import (
	"bot/repo"
	"crypto/sha1"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
)

func main() {
	statuses := make(map[int64]string)
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env: %s", err.Error())
	}

	db, err := repo.NewPostgresDB(repo.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(db)

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
			hashId, err := repo.CreateId(update.Message.From.ID, generateIdHash(update.Message.From.ID), db)
			if err != nil {
				log.Println(err)
			}

			msg.Text = fmt.Sprintf("Ссылка на аутентификация тг: %s (действует 12 часов)", generateLink(hashId))
		case "cancel":
			msg.Text = fmt.Sprintf("Действие %s отмененно", statuses[update.Message.From.ID])
			delete(statuses, update.Message.From.ID)
		default:
			msg.Text = "I don't know that command"
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

func generateIdHash(id int64) string {
	hash := sha1.New()
	hash.Write([]byte(strconv.FormatInt(id, 10)))
	return fmt.Sprintf("%x", hash.Sum([]byte("rheufhurhiuien")))
}

func generateLink(hashId string) string {
	return fmt.Sprintf("htpps://upfollow.ru/Auth-tg/%s", hashId)
}
