package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"
	"strings"
)

const TOKEN = "TELEGRAM_BOT_TOKEN"

// init is invoked before main()
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("Создай .env (смотри .env-sample)")
	}
}

func main() {
	var token, _ = os.LookupEnv(TOKEN)

	bot := createBot(token)

	//bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			input := strings.Split(update.Message.Text, " ")
			command, isValidCommand := handleInputCommand(input[0])

			//Проверка команды в сообщении на валидность
			if !isValidCommand {
				sendMessage(bot, update, "Я не понял, чё ты хочешь.")
				continue
			}

			// TODO: Handle commands
			sendMessage(bot, update, command)
		}
	}
}

func createBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	return bot
}

func handleInputCommand(command string) (string, bool) {
	var isText, err = regexp.MatchString(`^[а-яА-ЯёЁ]+$`, command)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(isText, err)

	return command, isText
}

func sendMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, message string) {

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)
}
