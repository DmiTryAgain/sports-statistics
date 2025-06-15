package app

import (
	"context"
	"fmt"

	"github.com/DmiTryAgain/sports-statistics/pkg/db"
	"github.com/DmiTryAgain/sports-statistics/pkg/tg"

	"github.com/go-pg/pg/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vmkteam/embedlog"
)

type App struct {
	embedlog.Logger

	cfg     tg.Config
	db      db.DB
	dbc     *pg.DB
	tgBot   *tgbotapi.BotAPI
	handler *tg.MessageHandler
}

func New(lg embedlog.Logger, db db.DB, dbc *pg.DB, cfg tg.Config) (*App, error) {
	// create tg bot
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		return nil, fmt.Errorf("create tgbot, err=%w", err)
	}

	bot.Debug = cfg.Bot.Debug
	a := &App{
		Logger: lg,
		cfg:    cfg,
		db:     db,
		dbc:    dbc,
		tgBot:  bot,
	}

	a.handler = tg.New(a.Logger, a.dbc, bot, a.cfg.Bot)

	return a, nil
}

func (a *App) Run(ctx context.Context) {
	a.handler.ListenAndHandle(ctx)
}

func (a *App) Shutdown() {
	a.Printf("shutting down ...")
	a.tgBot.StopReceivingUpdates()

	if err := a.dbc.Close(); err != nil {
		a.Errorf("failed to close database: %v", err)
	}
}
