package config

import (
	"github.com/go-pg/pg/v10"
	"time"
)

type Config struct {
	Database *pg.Options
	Bot      Bot
}

type Bot struct {
	Token       string
	Name        string
	ReplyFormat string
	Debug       bool
	Timeout     Duration `toml:"Timeout"`
}

// Duration is a parsable from toml time.Duration.
type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
