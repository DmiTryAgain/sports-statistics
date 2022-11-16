package dto

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Dto struct {
	upd *tgbotapi.Update
}

func (d *Dto) Construct(upd *tgbotapi.Update) *Dto {
	d.upd = upd

	return d
}

func (d *Dto) GetUserId() int {
	return d.upd.Message.From.ID
}

func (d *Dto) GetUserName() string {
	return d.upd.Message.Chat.UserName
}

func (d *Dto) GetChatType() string {
	return d.upd.Message.Chat.Type
}

func (d *Dto) GetChatId() int64 {
	return d.upd.Message.Chat.ID
}

func (d *Dto) GetText() string {
	return d.upd.Message.Text
}

func (d *Dto) GetTextEntities() *[]tgbotapi.MessageEntity {
	return d.upd.Message.Entities
}
