package update_handler

import (
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api"
	tgBotRequest "sports-statistics/internal/app/request"
	"sports-statistics/internal/config"
	"sports-statistics/internal/service/dto"
	"sports-statistics/internal/service/logger"
	"sports-statistics/internal/service/message_handler"
)

type UpdateHandler struct {
	bot     *tgBotApi.BotAPI
	updates *tgBotApi.UpdatesChannel
}

func (h *UpdateHandler) Construct(bot *tgBotApi.BotAPI, updates *tgBotApi.UpdatesChannel) *UpdateHandler {
	h.bot = bot
	h.updates = updates

	return h
}

func (h *UpdateHandler) Handle() {
	for update := range *h.updates {
		request := new(tgBotRequest.Request).Construct(&update)
		messageDto := h.createMessageDto(request)
		go h.logUpdates(new(logger.ConsoleLogger), messageDto)
		go h.handleMessage(new(message_handler.MessageHandler).Construct(), messageDto)
	}
}

func (h *UpdateHandler) logUpdates(l logger.Logger, messDto *dto.Dto) {
	l.Log(h.createLoggerDto(messDto))
}

func (h *UpdateHandler) createMessageDto(request *tgBotRequest.Request) *dto.Dto {
	return new(dto.Dto).Construct(request)
}

func (h *UpdateHandler) createLoggerDto(messDto *dto.Dto) *logger.Dto {
	return new(logger.Dto).Construct(messDto)
}

func (h *UpdateHandler) createHandlerDto(messDto *dto.Dto) *message_handler.Dto {
	return new(message_handler.Dto).Construct(messDto)
}

func (h *UpdateHandler) handleMessage(handler message_handler.Handler, messDto *dto.Dto) {
	msg, send, _ := handler.HandleWithResponse(h.createHandlerDto(messDto))

	if send {
		h.sendMessage(messDto.GetChatId(), messDto.GetMessageId(), msg)
	}

	handler.Destruct()
}

func (h *UpdateHandler) sendMessage(chatId int64, messageId int, message string) {
	msg := tgBotApi.NewMessage(chatId, message)
	msg.ParseMode = config.Configs.GetReplyFormat()
	msg.ReplyToMessageID = messageId

	h.bot.Send(msg)
}
