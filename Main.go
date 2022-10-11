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
	"time"
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

	conditionsByPeriod := map[string]string{
		"сегодня":   "YEAR(`created`) = YEAR(NOW()) AND WEEK(`created`, 1) = WEEK(NOW(), 1) AND DAY(`created`) = DAY(NOW())",
		"вчера":     "WHERE MONTH(`created`) = MONTH(DATE_ADD(NOW(), INTERVAL -1 DAY)) and YEAR(`created`) = YEAR(DATE_ADD(NOW(), INTERVAL -1 DAY))",
		"позавчера": "WHERE MONTH(`created`) = MONTH(DATE_ADD(NOW(), INTERVAL -2 DAY)) and YEAR(`created`) = YEAR(DATE_ADD(NOW(), INTERVAL -2 DAY))",
		"неделю":    "YEAR(`created`) = YEAR(NOW()) AND WEEK(`created`, 1) = WEEK(NOW(), 1)",
		"месяц":     "MONTH(`created`) = MONTH(NOW()) AND YEAR(`created`) = YEAR(NOW())",
		"год":       "YEAR(`created`) = YEAR(NOW())",
	}

	bot := createBot(token)

	//bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		input := prepareInput(update.Message.Text, " ")
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

		log.Printf("Message info: \n"+
			"UserName: %s\n"+
			"ChatId: %d\n"+
			"ChatType: %s\n"+
			"Message: %s\n"+
			"UserId: %d\n",
			update.Message.Chat.UserName,
			update.Message.Chat.ID,
			update.Message.Chat.Type,
			update.Message.Text,
			update.Message.From.ID,
		)

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

			insert.Close()

			if err != nil {
				panic(err)
			}

			findTrain.Close()

			sendMessage(bot, update, fmt.Sprintf("Добавлено %s %d ", training, count))

		case "покажи", "выведи":

			inputByPeriod := prepareInput(update.Message.Text, " за ")

			if len(inputByPeriod) < 2 {
				sendMessage(bot, update, "Введи корректно что именно показать и период.")
				continue
			}

			/** TODO */
			inputTrainings := deleteElemFromSlice(prepareInput(inputByPeriod[0], " "), 0)
			periods := deleteElemFromSlice(inputByPeriod, 0)
			var wheres []string
			var invalidPeriods []string

			for _, value := range periods {
				_, isValidText := checkIsText(value)

				if isValidText {
					if conditionsByPeriod[value] != "" {
						wheres = append(wheres, conditionsByPeriod[value])
					} else {
						invalidPeriods = append(invalidPeriods, value)
					}
				} else {
					numsPeriod, invalidDatePeriods := prepareDateInterval(value)
					if len(invalidDatePeriods) > 0 {
						for _, invalid := range invalidDatePeriods {
							invalidPeriods = append(invalidPeriods, invalid)
						}
					} else {
						if len(numsPeriod) == 1 {
							wheres = append(wheres, "DATE(`created`) = DATE('"+numsPeriod[0]+"')")
						}

						if len(numsPeriod) == 2 {
							wheres = append(wheres, "DATE(`created`) >= DATE('"+numsPeriod[0]+"' AND DATE(`created`) <= DATE('"+numsPeriod[1]+"')")
						}
					}
				}
			}

			/** TODO */
			findTrain, err := db.Query(fmt.Sprintf("SELECT * from `training` where `Name` = '%s' LIMIT 1", training))
			for _, where := range wheres {
				//Поиск упражнения в БД

				if err != nil {
					sendMessage(bot, update, fmt.Sprintf("Произошла ошибка БД: %d ", err))
					continue
				}
			}

			log.Printf("Message info: %v \n", wheres)
			log.Printf("Message info: %v \n", invalidPeriods)

			continue
			if len(input) < 3 {
				sendMessage(bot, update, "Введи корректно что именно показать и период.")
				continue
			}

			_, isValidTraining := checkIsText(input[1])

			//Проверка указанного упражнения в сообщении на валидность
			if !isValidTraining {
				sendMessage(bot, update, "Указанное упражнение содержит некорректные символы.")
				continue
			}

			if input[2] == "за" {
				input = deleteElemFromSlice(input, 2)
			}

			sendMessage(bot, update, fmt.Sprintf("Команда \"%s\" в разработке.", command))
		case "удали":
			sendMessage(bot, update, fmt.Sprintf("Команда \"%s\" в разработке.", command))
		case "help", "помоги", "помощь":
			sendMessage(
				bot,
				update,
				fmt.Sprintf(
					"Привет! Я - бот, который поможет вести статистику спортивных упражнений, которые "+
						"ты выполняешь. Ты же ведь занимаешься спортом, верно?🤔\n"+
						"Так вот, чтоб было удобно вести учёт и смотреть статистику, ты можешь это делать "+
						"с помощью команд ко мне.\n"+
						"Я слушаю команды, когда ко мне обращаются. обратись ко мне вот так: `@%s`\n"+
						"Исключением является личная переписка. Если ты напишешь мне в личку, я буду реагировать на "+
						"любые твои сообщения. Но и в личных сообщениях поддерживается обращение, "+
						"если уж сильно хочется)\n"+
						"После обращения через пробел нужно написать команду и передать к ней данные, "+
						"чтобы записать/показать результаты."+
						"Список поддерживаемых команд: \n"+
						"Чтобы записать результаты, воспользуйтесь командой ``, напишите название упражнения и "+
						"количество повторений, которое сделали. Все слова отделяйте пробелом.\n"+
						"Например, Вы сделали подход из 10 подтягиваний. Чтобы я всё корректно записал, напишите в "+
						"чат: `@%s сделал подтягивание 10`",
					botName,
					botName,
				),
			)
		default:
			sendMessage(bot, update, fmt.Sprintf("Команда \"%s\" не найдена.", command))
		}

		db.Close()
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
	msg.ParseMode = "markdown"
	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)
}

func prepareInput(inputText string, separator string) []string {
	needle := regexp.MustCompile(`[[:punct:]]`)
	replacePuncts := needle.ReplaceAllString(strings.ToLower(inputText), "")
	replaceSpaces := regexp.MustCompile("\\s+")
	replace := replaceSpaces.ReplaceAllString(strings.TrimSpace(replacePuncts), " ")

	return strings.Split(replace, separator)
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

func prepareDateInterval(interval string) ([]string, []string) {
	replaceSpaces := regexp.MustCompile("\\s+")
	replace := replaceSpaces.ReplaceAllString(strings.TrimSpace(interval), " ")
	intervals := strings.Split(replace, " ")
	var result []string
	var invalidPeriods []string

	if len(intervals) == 1 {
		formattedInterval, err := getDateFromNums(intervals[0])
		if err != nil {
			invalidPeriods = append(invalidPeriods, intervals[0])
		}

		return append(result, formattedInterval), invalidPeriods
	}

	if len(intervals) == 2 {
		formattedIntervalBegin, errBegin := getDateFromNums(intervals[0])
		formattedIntervalEnd, errEnd := getDateFromNums(intervals[1])

		if errBegin != nil {
			invalidPeriods = append(invalidPeriods, intervals[0])
		}

		if errEnd != nil {
			invalidPeriods = append(invalidPeriods, intervals[1])
		}

		if formattedIntervalBegin > formattedIntervalEnd {
			formattedIntervalBegin, formattedIntervalEnd = formattedIntervalEnd, formattedIntervalBegin
		}

		result = append(result, formattedIntervalBegin)
		return append(result, formattedIntervalEnd), invalidPeriods
	}

	invalidPeriods = append(invalidPeriods, "all")

	return result, invalidPeriods
}

func getDateFromNums(nums string) (string, error) {
	parse, err := time.Parse("02012006", nums)
	if err != nil {
		parse, err = time.Parse("020106", nums)
		if err != nil {
			return "", err
		}
	}

	return parse.Format("2006-01-02"), nil
}
