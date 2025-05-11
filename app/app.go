package app

import (
	"github.com/DmiTryAgain/sports-statistics/config"
	"github.com/DmiTryAgain/sports-statistics/pkg/db"
	"github.com/DmiTryAgain/sports-statistics/pkg/embedlog"
	"github.com/DmiTryAgain/sports-statistics/pkg/tg"
	"github.com/pkg/errors"

	"github.com/go-pg/pg/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	embedlog.Logger

	cfg     config.Config
	db      db.DB
	dbc     *pg.DB
	tgBot   *tgbotapi.BotAPI
	handler *tg.MessageHandler
}

func New(verbose bool, db db.DB, dbc *pg.DB, cfg config.Config) (*App, error) {
	// create tg bot
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tgbot")
	}

	bot.Debug = cfg.Bot.Debug
	a := &App{
		cfg: cfg,
		db:  db,
		dbc: dbc,
	}

	a.SetStdLoggers(verbose)
	a.handler = tg.New(a.Logger, a.db, a.dbc, bot, a.cfg.Bot)

	return a, nil
}

func (a *App) Run() {
	a.handler.ListenAndHandle()
}

func (a *App) Shutdown() {
	a.Printf("shutting down ...")
	a.tgBot.StopReceivingUpdates()

	if err := a.dbc.Close(); err != nil {
		a.Errorf("failed to close database: %v", err)
	}
}
