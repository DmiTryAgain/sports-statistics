package config

import (
	"github.com/joho/godotenv"
	"os"
)

const botName = "TELEGRAM_BOT_NAME"
const token = "TELEGRAM_BOT_TOKEN"
const replyFormat = "TELEGRAM_REPLY_FORMAT"
const updsTimeout = 60

const dbDsn = "DB_DSN"
const dbType = "DB_TYPE"

var Configs = new(Config).Construct()

type Config struct {
	name, secret, dsn, typeDb, replyFormat string
	updsTimeout                            int
}

func (b *Config) Construct() *Config {
	if err := godotenv.Load(".env", "docker\\.env"); err != nil {
		panic(err)
	}

	var token, _ = os.LookupEnv(token)
	var name, _ = os.LookupEnv(botName)
	var replyFormat, _ = os.LookupEnv(replyFormat)

	var dbDsn, _ = os.LookupEnv(dbDsn)
	var dbType, _ = os.LookupEnv(dbType)

	b.name = name
	b.secret = token
	b.replyFormat = replyFormat
	b.updsTimeout = updsTimeout

	b.dsn = dbDsn
	b.typeDb = dbType

	return b
}

func (b *Config) GetBotName() string {
	return b.name
}

func (b *Config) GetBotSecret() string {
	return b.secret
}

func (b *Config) GetReplyFormat() string {
	return b.replyFormat
}

func (b *Config) GetBotUpdatesTimeout() int {
	return b.updsTimeout
}

func (b *Config) GetDbDsn() string {
	return b.dsn
}

func (b *Config) GetDbType() string {
	return b.typeDb
}
