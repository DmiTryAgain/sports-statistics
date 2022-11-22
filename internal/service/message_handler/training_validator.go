package message_handler

import (
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"regexp"
	"sports-statistics/internal/service/helpers"
)

const textEntityMention = "mention"
const minCorrectLenWordsSlice = 2
const minCorrectLenForPeriodsSlice = 2

type TrainingValidator struct {
}

func (v TrainingValidator) CheckBotCall(entities *[]tgBotApi.MessageEntity, firstWord string, botName string) bool {
	if entities == nil {
		return false
	}

	for _, entity := range *entities {
		if entity.Type != textEntityMention || firstWord != botName {
			return false
		}
	}

	return true
}

func (v TrainingValidator) IsEmptyMessage(words []string) bool {
	return helpers.SliceHelper{}.IsEmptySlice(words)
}
func (v TrainingValidator) CheckIsOnlyRussianText(str string) (bool, error) {
	return regexp.MatchString(`^[а-яА-ЯёЁ]+$`, str)
}
func (v TrainingValidator) CheckIsOnlyInt(count string) (bool, error) {
	return regexp.MatchString(`^(\d+)$`, count)
}

func (v TrainingValidator) CheckMinCorrectLen(words []string) bool {
	return helpers.SliceHelper{}.CheckMinLenSlice(words, minCorrectLenWordsSlice)
}

func (v TrainingValidator) CheckMinCorrectLenForPeriods(words []string) bool {
	return helpers.SliceHelper{}.CheckMinLenSlice(words, minCorrectLenForPeriodsSlice)
}
