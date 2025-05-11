package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DmiTryAgain/sports-statistics/app"
	"github.com/DmiTryAgain/sports-statistics/config"
	"github.com/DmiTryAgain/sports-statistics/pkg/db"

	"github.com/BurntSushi/toml"
	"github.com/go-pg/pg/v10"
)

var (
	flConfigPath = flag.String("config", "config/local.toml", "Path to config file")
	flVerboseSQL = flag.Bool("verbose_sql", false, "enable all sql output")
	flVerbose    = flag.Bool("verbose", false, "enable debug output")
	cfg          config.Config
)

func main() {
	flag.Parse()

	// read config
	if _, err := toml.DecodeFile(*flConfigPath, &cfg); err != nil {
		die(err)
	}

	// check db connection
	dbconn := pg.Connect(cfg.Database)
	dbc := db.New(dbconn)
	v, err := dbc.Version()
	die(err)
	log.Println(v)

	// log all sql queries
	if *flVerboseSQL {
		sqlLogger := log.New(os.Stdout, "Q", log.LstdFlags)
		dbc.AddQueryHook(db.NewQueryLogger(sqlLogger))
		if dbconn != nil {
			sqlLoggerR := log.New(os.Stdout, "QR", log.LstdFlags)
			dbconn.AddQueryHook(db.NewQueryLogger(sqlLoggerR))
		}
	}

	a, err := app.New(*flVerbose, dbc, dbconn, cfg, bot)
	if err != nil {
		die(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go a.Run()

	<-quit
	a.Shutdown()
}

// die calls log.Fatal if err wasn't nil.
func die(err error) {
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}
}
