package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/DmiTryAgain/sports-statistics/app"
	"github.com/DmiTryAgain/sports-statistics/config"
	"github.com/DmiTryAgain/sports-statistics/pkg/db"

	"github.com/BurntSushi/toml"
	"github.com/go-pg/pg/v10"
	"github.com/namsral/flag"
	"github.com/vmkteam/embedlog"
)

var (
	fs           = flag.NewFlagSetWithEnvPrefix(os.Args[0], "SELLERSRV", 0)
	flConfigPath = fs.String("config", "config/local.toml", "Path to config file")
	flVerbose    = fs.Bool("verbose", false, "enable debug output")
	flJSONLogs   = fs.Bool("json", false, "enable json output")
	flDev        = fs.Bool("dev", false, "enable dev mode")
	cfg          config.Config
)

func main() {
	flag.Parse()

	// read config
	if _, err := toml.DecodeFile(*flConfigPath, &cfg); err != nil {
		exitOnError(err)
	}

	// setup logger
	sl, ctx := embedlog.NewLogger(*flVerbose, *flJSONLogs), context.Background()
	if *flDev {
		sl = embedlog.NewDevLogger()
	}
	slog.SetDefault(sl.Log()) // set default logger

	// check db connection
	dbconn := pg.Connect(cfg.Database)
	dbc := db.New(dbconn)
	v, err := dbc.Version()
	exitOnError(err)
	sl.Print(ctx, "connected to db", "version", v)

	// log all sql queries
	if *flDev {
		dbc.AddQueryHook(db.NewQueryLogger(sl))
	}

	a, err := app.New(sl, dbc, dbconn, cfg)
	if err != nil {
		exitOnError(err)
	}

	sl.Print(ctx, "start application")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go a.Run(ctx)

	<-quit
	a.Shutdown()
}

// exitOnError calls os.Exit if err wasn't nil.
func exitOnError(err error) {
	if err != nil {
		//nolint:sloglint
		slog.Error(err.Error())
		os.Exit(1)
	}
}
