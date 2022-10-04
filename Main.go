package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const TOKEN = "TELEGRAM_BOT_TOKEN"
const BOT_NAME = "TELEGRAM_BOT_NAME"

const DB_USERNAME = "DB_USERNAME"
const DB_PASSWORD = "DB_PASSWORD"
const DB_DATABASE = "DB_DATABASE"
const DB_CHARSET = "DB_CHARSET"
const DB_HOST = "DB_HOST"
const DB_DSN = "DB_DSN"
const DB_TYPE = "DB_TYPE"

const CHAT_TYPE_GROUP = "group"

type Training struct {
	Id          int
	Alias, Name string
}

// init is invoked before main()
func init() {
	if err := godotenv.Load(".env", "docker\\.env"); err != nil {
		log.Print("Создай .env (смотри .env-sample)")
		panic(err)
	}
}

func main() {
	var token, _ = os.LookupEnv(TOKEN)
	var botName, _ = os.LookupEnv(BOT_NAME)
	//var dbUser, _ = os.LookupEnv(DB_USERNAME)
	//var dbPass, _ = os.LookupEnv(DB_PASSWORD)
	//var dbName, _ = os.LookupEnv(DB_DATABASE)
	//var dbCharset, _ = os.LookupEnv(DB_CHARSET)
	//var dbHost, _ = os.LookupEnv(DB_HOST)
	var dbDsn, _ = os.LookupEnv(DB_DSN)
	var dbType, _ = os.LookupEnv(DB_TYPE)

	bot := createBot(token)

	//bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			//log.Printf("[%s] %s", update.Message.Chat.Type, update.Message.Text)
			input := prepareInput(update.Message.Text)
			lowerBotName := strings.ToLower(botName)

			if !checkBotCall(update, input[0], lowerBotName) && update.Message.Chat.Type == CHAT_TYPE_GROUP {
				continue
			} else if checkBotCall(update, input[0], lowerBotName) {
				input = deleteElemFromSlice(input, 0)

				if len(input) == 0 {
					sendMessage(bot, update, "Чё?")
					continue
				}
			}

			command, isValidCommand := checkIsText(input[0])

			//Проверка команды в сообщении на валидность
			if !isValidCommand {
				sendMessage(bot, update, "Команда содержит недопустимые символы.")
				continue
			}

			db, err := sql.Open(dbType, dbDsn)

			if err != nil {
				sendMessage(bot, update, fmt.Sprintf("Ошибка подключения к базе данных: %d ", err))
			}

			switch command {
			case "сделал":

				if len(input) < 3 {
					sendMessage(bot, update, "Введи корректное наименование упражнения и число повторений.")
					continue
				}

				training, isValidTraining := checkIsText(input[1])

				//Проверка указанного упражнения в сообщении на валидность
				if !isValidTraining {
					sendMessage(bot, update, "Указанное упражнение содержит некорректные символы.")
					continue
				}

				count, isValidCount := checkIsInt(input[2])

				//Проверка указанного количества в сообщении на валидность
				if !isValidCount {
					sendMessage(bot, update, "Указанное количество содержит некорректные символы.")
					continue
				}

				//Поиск упражнения в БД
				findTrain, err := db.Query(fmt.Sprintf("SELECT * from `training` where `Name` = '%s' LIMIT 1", training))

				if err != nil {
					sendMessage(bot, update, fmt.Sprintf("Произошла ошибка БД: %d ", err))
					continue
				}

				var trainings []Training
				var train Training

				for findTrain.Next() {
					err = findTrain.Scan(&train.Id, &train.Alias, &train.Name)
					if err != nil {
						sendMessage(bot, update, fmt.Sprintf("Произошла ошибка БД: %d ", err))
					}

					trainings = append(trainings, train)
				}

				if train.Id == 0 {
					sendMessage(bot, update, fmt.Sprintf("Упражнение \"%s\" не найдено.", training))
					err := findTrain.Close()

					if err != nil {
						sendMessage(bot, update, fmt.Sprintf("Произошла ошибка БД: %d ", err))
					}

					continue
				}

				insert, err := db.Query(
					fmt.Sprintf(
						"INSERT INTO `statistic` (`telegram_user_id`, `training_id`, `count`) VALUES('%d', '%d', '%d')",
						update.Message.From.ID,
						train.Id,
						count,
					),
				)

				defer insert.Close()

				if err != nil {
					panic(err)
				}

				defer findTrain.Close()

				sendMessage(bot, update, fmt.Sprintf("Добавлено %s %d ", training, count))

			case "удали":
			case "покажи":
			default:
				sendMessage(bot, update, "Команда введена некорректно.")
			}

			defer db.Close()
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

func checkIsText(command string) (string, bool) {
	var isText, err = regexp.MatchString(`^[а-яА-ЯёЁ]+$`, command)
	if err != nil {
		fmt.Println(err)
	}

	return command, isText
}

func checkIsInt(count string) (int, bool) {
	var isInt, err = regexp.MatchString(`^(\d+)$`, count)
	if err != nil {
		fmt.Println(err)
	}

	countInt, _ := strconv.Atoi(count)

	return countInt, isInt
}

func sendMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, message string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)
}

func prepareInput(inputText string) []string {
	needle := regexp.MustCompile(`[[:punct:]]`)
	return strings.Split(needle.ReplaceAllString(strings.ToLower(inputText), ""), " ")
}

func checkBotCall(update tgbotapi.Update, firstWord string, botName string) bool {
	if update.Message.Entities == nil {
		return false
	}

	for _, Entity := range *update.Message.Entities {
		if Entity.Type != "mention" || firstWord != botName {
			return false
		}
	}

	return true
}

// Удалить элемент по индексу из слайса.
func deleteElemFromSlice(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}
