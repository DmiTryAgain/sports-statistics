package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"sports-statistics/internal/config"
	"sports-statistics/internal/service/update_handler"
)

var Config = new(config.Config).Construct()

func Run() {
	updates, err := getListenTelegramBotUpdates()

	if err != nil {
		panic(err)
	}

	update_handler.Handle(&updates)
}

func getListenTelegramBotUpdates() (tgbotapi.UpdatesChannel, error) {
	bot := createBot(Config.GetBotSecret())
	/** TODO: Разобрать, что будет, если включить. Вынести в энвы, если надо */
	//bot.Debug = true

	/** TODO: Разобрать, что за NewUpdate() и почему 0 */
	u := tgbotapi.NewUpdate(0)
	u.Timeout = Config.GetBotUpdatesTimeout()

	return bot.GetUpdatesChan(u)
}

func createBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		panic(err)
	}

	return bot
}
