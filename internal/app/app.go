package app

import (
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"sports-statistics/internal/config"
	"sports-statistics/internal/service/update_handler"
)

func Run() {
	bot := createBot(config.Configs.GetBotSecret())
	updates, err := getListenTelegramBotUpdates(bot)

	if err != nil {
		panic(err)
	}

	handler := new(update_handler.UpdateHandler).Construct(bot, &updates)
	handler.Handle()
}

func getListenTelegramBotUpdates(bot *tgBotApi.BotAPI) (tgBotApi.UpdatesChannel, error) {
	/** TODO: Разобрать, что будет, если включить. Вынести в энвы, если надо */
	//bot.Debug = true

	/** TODO: Разобрать, что за NewUpdate() и почему 0 */
	u := tgBotApi.NewUpdate(0)
	u.Timeout = config.Configs.GetBotUpdatesTimeout()

	return bot.GetUpdatesChan(u)
}

func createBot(token string) *tgBotApi.BotAPI {
	bot, err := tgBotApi.NewBotAPI(token)

	if err != nil {
		panic(err)
	}

	return bot
}
