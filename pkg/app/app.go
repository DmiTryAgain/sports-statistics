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

	cfg            tg.Config
	dbc            db.DB
	dbConn         *pg.DB
	tgBot          *tgbotapi.BotAPI
	handler        *tg.MessageHandler
	shutdownCancel context.CancelFunc
}

func New(lg embedlog.Logger, dbc db.DB, dbConn *pg.DB, cfg tg.Config) (*App, error) {
	// create tg bot
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		return nil, fmt.Errorf("create tgbot, err=%w", err)
	}

	bot.Debug = cfg.Bot.Debug
	a := &App{
		Logger: lg,
		cfg:    cfg,
		dbc:    dbc,
		dbConn: dbConn,
		tgBot:  bot,
	}

	return a, nil
}

func (a *App) Run() {
	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())
	a.shutdownCancel = shutdownCancel
	a.handler = tg.New(shutdownCtx, a.Logger, a.dbc, db.NewStatisticRepo(a.dbc), a.tgBot, a.cfg.Bot)
	a.handler.ListenAndHandle(shutdownCtx)
}

func (a *App) Shutdown() {
	a.Printf("shutting down ...")
	a.tgBot.StopReceivingUpdates()
	a.shutdownCancel()

	if err := a.dbc.Close(); err != nil {
		a.Errorf("failed to close database: %v", err)
	}
}
