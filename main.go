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

const (
	salt = "rheufhurhiuien"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env: %s", err.Error())
	}

	db, err := repo.NewPostgresDB(repo.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("База данных успешно подключена")

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	log.Println("Бот запущен")

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
			msg.Text = startMessage
		case "start":
			msg.Text = startMessage
		case "reg":
			hashId, err := repo.CreateId(update.Message.From.ID, generateIdHash(update.Message.From.ID), update.Message.Chat.ID, db)
			if err != nil {
				log.Println(err)
				msg.Text = errorMessage
				msg.ParseMode = "MarkdownV2"

				goto End
			}

			if hashId == "Уже есть" {
				msg.Text = "У вас уже есть ссылка"
				hash, timeM, err := repo.GetCurentHash(update.Message.From.ID, db)
				if err != nil {
					msg.Text = errorMessage
					msg.ParseMode = "MarkdownV2"

					goto End
				} else if timeM != "" {
					msg.Text = fmt.Sprintf("Ссылка на аутентификация тг: %s %s", generateLink(hash), timeM)
				} else {
					if err := repo.DeleteLink(update.Message.From.ID, db); err != nil {
						msg.Text = errorMessage
						msg.ParseMode = "MarkdownV2"

						goto End
					}

					msg.Text = "Ваша ссылка уже не дейстивтельна, вым надо получить новую через команду /reg"
				}
			} else {
				msg.Text = fmt.Sprintf("Ссылка на аутентификация тг: %s (действует 12 часов)", generateLink(hashId))
			}
		case "cancel":
			if err := repo.DeleteLink(update.Message.From.ID, db); err != nil {
				msg.Text = errorMessage
				msg.ParseMode = "MarkdownV2"

				goto End
			}
			msg.Text = "Ссылка удалена"
		default:
			msg.Text = "I don't know that command"
		}
	End:
		logs(update.Message.From.UserName, update.Message.Text, msg.Text)

		if _, err := bot.Send(msg); err != nil {
			log.Println(msg)
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
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func generateLink(hashId string) string {
	return fmt.Sprintf("https://upfollow.ru/Auth-tg/%s", hashId)
}
