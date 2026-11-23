package tg

import (
	"time"

	"github.com/go-pg/pg/v10"
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
	// TimeZone Часовой пояс для границ периодов и дат в статистике, например "Europe/Moscow". Если не задан, берётся пояс сервера
	TimeZone Location `toml:"TimeZone"`
}

// Duration is a parsable from toml time.Duration.
type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

// Location is a parsable from toml time.Location.
type Location struct {
	*time.Location
}

func (l *Location) UnmarshalText(text []byte) (err error) {
	l.Location, err = time.LoadLocation(string(text))
	return err
}
