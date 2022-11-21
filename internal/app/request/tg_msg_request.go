package request

import tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api"

type Request struct {
	update *tgBotApi.Update
}

func (r *Request) Construct(update *tgBotApi.Update) *Request {
	r.update = update

	return r
}

func (r *Request) GetUpdate() *tgBotApi.Update {
	return r.update
}
