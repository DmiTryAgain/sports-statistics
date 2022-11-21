package update_handler

import (
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"sports-statistics/internal/app"
	tgBotRequest "sports-statistics/internal/app/request"
	"sports-statistics/internal/service/dto"
	"sports-statistics/internal/service/logger"
	"sports-statistics/internal/service/message_handler"
)

func Handle(bot *tgBotApi.BotAPI, updates *tgBotApi.UpdatesChannel) {
	for update := range *updates {
		request := new(tgBotRequest.Request).Construct(&update)
		messageDto := createMessageDto(request)
		go logUpdates(new(logger.ConsoleLogger), messageDto)
		go handleMessage(bot, new(message_handler.MessageHandler).Construct(), messageDto)
	}
}

func logUpdates(l logger.Logger, messDto *dto.Dto) {
	l.Log(createLoggerDto(messDto))
}

func createMessageDto(request *tgBotRequest.Request) *dto.Dto {
	return new(dto.Dto).Construct(request)
}

func createLoggerDto(messDto *dto.Dto) *logger.Dto {
	return new(logger.Dto).Construct(messDto)
}

func createHandlerDto(messDto *dto.Dto) *message_handler.Dto {
	return new(message_handler.Dto).Construct(messDto)
}

func handleMessage(bot *tgBotApi.BotAPI, handler message_handler.Handler, messDto *dto.Dto) {
	msg, send, _ := handler.HandleWithResponse(createHandlerDto(messDto))

	if send {
		sendMessage(bot, messDto.GetChatId(), messDto.GetMessageId(), msg)
	}
}

func sendMessage(bot *tgBotApi.BotAPI, chatId int64, messageId int, message string) {
	msg := tgBotApi.NewMessage(chatId, message)
	msg.ParseMode = app.Config.GetReplyFormat()
	msg.ReplyToMessageID = messageId

	bot.Send(msg)
}
