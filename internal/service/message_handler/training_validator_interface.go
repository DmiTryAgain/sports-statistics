package message_handler

import tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api"

type TrainingValidatorInterface interface {
	CheckBotCall(entities *[]tgBotApi.MessageEntity, firstWord string, botName string) bool
	IsEmptyMessage(words []string) bool
	CheckIsOnlyRussianText(str string) (bool, error)
	CheckIsOnlyInt(count string) (bool, error)
	CheckMinCorrectLen(words []string) bool
	CheckMinCorrectLenForPeriods(words []string) bool
}
