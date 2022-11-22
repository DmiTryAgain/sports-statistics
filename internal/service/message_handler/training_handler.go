package message_handler

import (
	"fmt"
	"sports-statistics/internal/config"
	cr "sports-statistics/internal/repositiry/command_repository"
	sr "sports-statistics/internal/repositiry/db/statistic"
	tr "sports-statistics/internal/repositiry/db/training"
	"sports-statistics/internal/repositiry/periods_repository"
	"sports-statistics/internal/service/entity/statistic"
	"sports-statistics/internal/service/entity/user"
	"sports-statistics/internal/service/helpers"
	"strconv"
	"strings"
	"time"
)

const chatTypeGroup = "group"
const dbErrorMessage = "Произошла ошибка БД!"

type MessageHandler struct {
	validator   TrainingValidatorInterface
	sliceHelper *helpers.SliceHelper
}

func (m *MessageHandler) Construct() Handler {
	m.validator = new(TrainingValidator)
	m.sliceHelper = new(helpers.SliceHelper)

	return m
}

func (m *MessageHandler) HandleWithResponse(dto *Dto) (string, bool, error) {
	firstSliceIndex := m.sliceHelper.FirstSliceElemIndex()
	wordsFromMessText := m.sliceHelper.SplitStringToSlice(dto.GetText(), " ")
	lowerBotName := strings.ToLower(config.Configs.GetBotName())
	isBotCalled := m.validator.CheckBotCall(dto.GetTextEntities(), wordsFromMessText[firstSliceIndex], lowerBotName)

	if !isBotCalled && dto.GetChatType() == chatTypeGroup {
		return "", false, nil
	} else if isBotCalled {
		wordsFromMessText = m.sliceHelper.DeleteElemFromSlice(wordsFromMessText, firstSliceIndex)
	}

	if m.validator.IsEmptyMessage(wordsFromMessText) {
		//sendMessage(bot, update, )
		return "Чё?", true, nil
	}

	command := wordsFromMessText[firstSliceIndex]
	wordsFromMessText = m.sliceHelper.DeleteElemFromSlice(wordsFromMessText, firstSliceIndex)
	commandIsValid, err := m.validator.CheckIsOnlyRussianText(command)

	if err != nil {
		return "Произошла ошибка при проверке команды на валидность!", true, err
	}

	if !commandIsValid {
		return "Команда должна состоять только из русских букв.", true, nil
	}

	commandRepository := new(cr.CommandRepository).Construct()

	_, isAddCommand := commandRepository.GetAddCommands()[command]
	_, isShowCommand := commandRepository.GetShowCommands()[command]

	switch true {
	case isAddCommand:
		return m.handleAddCommand(dto, wordsFromMessText)
	case isShowCommand:
		return m.handleShowCommand(dto)
	}

	return "Упс! Произошло некорректное поведение! Обратитесь за помощью к разработчику!", true, nil
}

func (m *MessageHandler) handleAddCommand(dto *Dto, words []string) (string, bool, error) {
	if !m.validator.CheckMinCorrectLen(words) {
		return "Введи корректное наименование упражнения и число повторений.", true, nil
	}

	training := words[m.sliceHelper.FirstSliceElemIndex()]

	trainingIsValid, err := m.validator.CheckIsOnlyRussianText(training)

	if err != nil {
		return "Произошла ошибка при проверке упражнения на валидность!", true, err
	}

	if !trainingIsValid {
		return "Указанное упражнение должно состоять только из русских букв.", true, nil
	}

	count := words[m.sliceHelper.SecondSliceElemIndex()]

	isValidCount, err := m.validator.CheckIsOnlyInt(count)

	if err != nil {
		return "Произошла ошибка при проверке количества повторений на валидность!", true, err
	}

	if !isValidCount {
		return "Указанное количество повторений должно состоять только из цифр.", true, nil
	}

	countInt, err := strconv.Atoi(count)

	if err != nil {
		return "Произошла ошибка при проверке количества повторений на валидность!", true, err
	}

	trainingRepository := new(tr.Repository).Construct()
	statisticRepository := new(sr.Repository).Construct()

	defer statisticRepository.Destruct()
	defer trainingRepository.Destruct()

	trainingEntity := trainingRepository.GetTrainingByName(training)

	if err != nil {
		fmt.Println(err)
		return dbErrorMessage, true, err
	}

	statisticRepository.AddStatistic(
		new(statistic.Statistic).Construct(
			nil,
			trainingEntity,
			new(statistic.User).Construct(new(user.Id).Construct(dto.GetUserId())),
			nil,
			new(statistic.Count).Construct(countInt),
			nil,
		),
	)

	err = statisticRepository.GetError()
	fmt.Println(err)
	if err != nil {
		return dbErrorMessage, true, err
	}

	return fmt.Sprintf("Добавлено %s %d ", training, countInt), true, nil
}

func (m *MessageHandler) handleShowCommand(dto *Dto) (string, bool, error) {
	inputByPeriod := m.sliceHelper.SplitStringToSlice(dto.GetText(), " за ")

	if !m.validator.CheckMinCorrectLenForPeriods(inputByPeriod) {
		return "Введи корректно что именно показать и период.", true, nil
	}

	firstSliceIndex := m.sliceHelper.FirstSliceElemIndex()

	inputTrainings := m.sliceHelper.DeleteElemFromSlice(
		m.sliceHelper.SplitStringToSlice(
			inputByPeriod[firstSliceIndex],
			" ",
		),
		firstSliceIndex,
	)

	inputTrainingsAnyElems := m.sliceHelper.ConvertFromStringToAnyElems(inputTrainings)

	periods := m.sliceHelper.DeleteElemFromSlice(inputByPeriod, firstSliceIndex)
	statisticRepository := new(sr.Repository).Construct()

	defer statisticRepository.Destruct()

	correctPeriods, _, err := m.prepareCorrectAndInvalidPeriods(periods)

	if err != nil {
		return "Произошла ошибка при проверке периода на валидность!", true, err
	}

	stat := statisticRepository.GetByConditions(inputTrainingsAnyElems, correctPeriods, dto.GetUserId())

	if m.sliceHelper.IsEmptySliceStatisticEntity(stat) {
		return "К сожалению по вашему запросу результаты не найдены. " +
			"Попробуйте указать другие упражнения и период.", true, err
	} else {
		resultMessage := "Вы сделали:\n"

		for _, result := range stat {
			resultMessage += fmt.Sprintf(
				"%v в количестве %v раз, за %v подхода(ов).",
				result.GetTraining().GetName().GetValue(),
				result.GetCount().GetValue(),
				result.GetsSets().GetValue(),
			)
		}

		return resultMessage, true, nil
	}

}

func (m *MessageHandler) prepareDateInterval(interval string) ([]string, []string) {
	firstSliceIndex := m.sliceHelper.FirstSliceElemIndex()
	secondSliceIndex := m.sliceHelper.SecondSliceElemIndex()
	intervals := m.sliceHelper.SplitStringDatesToSlice(interval)
	var result []string
	var invalidPeriods []string

	if m.sliceHelper.CheckLenSlice(intervals, 1) {
		formattedInterval, err := m.getDateFromNums(intervals[firstSliceIndex])

		if err != nil {
			invalidPeriods = append(invalidPeriods, intervals[firstSliceIndex])
		}

		return append(result, formattedInterval), invalidPeriods
	}

	if m.sliceHelper.CheckLenSlice(intervals, 2) {
		formattedIntervalBegin, errBegin := m.getDateFromNums(intervals[firstSliceIndex])
		formattedIntervalEnd, errEnd := m.getDateFromNums(intervals[secondSliceIndex])

		if errBegin != nil {
			invalidPeriods = append(invalidPeriods, intervals[firstSliceIndex])
		}

		if errEnd != nil {
			invalidPeriods = append(invalidPeriods, intervals[secondSliceIndex])
		}

		if formattedIntervalBegin > formattedIntervalEnd {
			formattedIntervalBegin, formattedIntervalEnd = formattedIntervalEnd, formattedIntervalBegin
		}

		result = append(result, formattedIntervalBegin)
		return append(result, formattedIntervalEnd), invalidPeriods
	}

	invalidPeriods = append(invalidPeriods, "all")

	return result, invalidPeriods
}

func (m *MessageHandler) getDateFromNums(nums string) (string, error) {
	parse, err := time.Parse("02012006", nums)
	if err != nil {
		parse, err = time.Parse("020106", nums)
		if err != nil {
			return "", err
		}
	}

	return parse.Format("2006-01-02"), nil
}

func (m *MessageHandler) prepareCorrectAndInvalidPeriods(periods []string) ([]string, []string, error) {
	var correctPeriods []string
	var invalidPeriods []string
	periodsRepository := new(periods_repository.Repository).Construct()
	firstSliceIndex := m.sliceHelper.FirstSliceElemIndex()

	for _, period := range periods {
		isValidText, err := m.validator.CheckIsOnlyRussianText(period)

		if err != nil {
			return correctPeriods, invalidPeriods, err
		}

		if isValidText {
			val, ok := periodsRepository.GetConditionsByPeriod(period)
			if ok {
				correctPeriods = append(correctPeriods, val)
			} else {
				invalidPeriods = append(invalidPeriods, period)
			}
		} else {
			numsPeriod, invalidDatePeriods := m.prepareDateInterval(period)
			if !m.sliceHelper.IsEmptySlice(invalidDatePeriods) {
				for _, invalid := range invalidDatePeriods {
					invalidPeriods = append(invalidPeriods, invalid)
				}
			} else {
				if m.sliceHelper.CheckLenSlice(numsPeriod, 1) {
					correctPeriods = append(
						correctPeriods,
						periodsRepository.GetConditionsByDate(numsPeriod[firstSliceIndex]),
					)
				}

				if m.sliceHelper.CheckLenSlice(numsPeriod, 2) {
					correctPeriods = append(
						correctPeriods,
						periodsRepository.GetConditionsByDateInterval(
							numsPeriod[firstSliceIndex],
							numsPeriod[m.sliceHelper.SecondSliceElemIndex()],
						),
					)
				}
			}
		}
	}

	return correctPeriods, invalidPeriods, nil
}
