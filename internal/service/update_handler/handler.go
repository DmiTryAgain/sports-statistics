package update_handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"sports-statistics/internal/service/dto"
	"sports-statistics/internal/service/logger"
	"sports-statistics/internal/service/message_handler"
)

func Handle(updates *tgbotapi.UpdatesChannel) {
	for update := range *updates {
		messageDto := createMessageDto(&update)
		go logUpdates(new(logger.ConsoleLogger), messageDto)
		go handleMessage(new(message_handler.MessageHandler).Construct(), messageDto)
	}
}

func logUpdates(l logger.Logger, messDto *dto.Dto) {
	l.Log(createLoggerDto(messDto))
}

func createMessageDto(upd *tgbotapi.Update) *dto.Dto {
	return new(dto.Dto).Construct(upd)
}

func createLoggerDto(messDto *dto.Dto) *logger.Dto {
	return new(logger.Dto).Construct(messDto)
}

func createHandlerDto(messDto *dto.Dto) *message_handler.Dto {
	return new(message_handler.Dto).Construct(messDto)
}

func handleMessage(handler message_handler.Handler, messDto *dto.Dto) {
	handler.HandleWithResponse(createHandlerDto(messDto))
}
