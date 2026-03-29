package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/DmiTryAgain/sports-statistics/pkg/app"
	"github.com/DmiTryAgain/sports-statistics/pkg/db"
	"github.com/DmiTryAgain/sports-statistics/pkg/tg"

	"github.com/BurntSushi/toml"
	"github.com/go-pg/pg/v10"
	"github.com/namsral/flag"
	"github.com/vmkteam/embedlog"
)

var (
	fs           = flag.NewFlagSetWithEnvPrefix(os.Args[0], "SPORTSTAT", 0)
	flConfigPath = fs.String("config", "config/local.toml", "Path to config file")
	flVerbose    = fs.Bool("verbose", false, "enable debug output")
	flJSONLogs   = fs.Bool("json", false, "enable json output")
	flDev        = fs.Bool("dev", false, "enable dev mode")
	cfg          tg.Config
)

func main() {
	flag.DefaultConfigFlagname = "config.flag"
	exitOnError(fs.Parse(os.Args[1:]))

	// read config
	if _, err := toml.DecodeFile(*flConfigPath, &cfg); err != nil {
		exitOnError(err)
	}

	// setup logger
	sl := embedlog.NewLogger(*flVerbose, *flJSONLogs)
	if *flDev {
		sl = embedlog.NewDevLogger()
	}
	slog.SetDefault(sl.Log()) // set default logger

	// check db connection
	dbconn := pg.Connect(cfg.Database)
	dbc := db.New(dbconn)
	v, err := dbc.Version()
	exitOnError(err)
	sl.Printf("connected to db, version=%s", v)

	// log all sql queries
	if *flDev {
		dbc.AddQueryHook(db.NewQueryLogger(sl))
	}

	a, err := app.New(sl, dbc, dbconn, cfg)
	if err != nil {
		exitOnError(err)
	}

	sl.Printf("start application")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go a.Run()

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
