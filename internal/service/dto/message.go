package dto

import (
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api"
	tgBotRequest "sports-statistics/internal/app/request"
)

type Dto struct {
	upd *tgBotApi.Update
}

func (d *Dto) Construct(request *tgBotRequest.Request) *Dto {
	d.upd = request.GetUpdate()

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

func (d *Dto) GetTextEntities() *[]tgBotApi.MessageEntity {
	return d.upd.Message.Entities
}

func (d *Dto) GetMessageId() int {
	return d.upd.Message.MessageID
}
