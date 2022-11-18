package message_handler

import (
	"fmt"
	"log"
	"sports-statistics/internal/app"
	cr "sports-statistics/internal/repositiry/command_repository"
	sr "sports-statistics/internal/repositiry/db/statistic"
	tr "sports-statistics/internal/repositiry/db/training"
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
	lowerBotName := strings.ToLower(app.Config.GetBotName())
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
		return m.handleAddCommand(wordsFromMessText, dto)
	case isShowCommand:
		return m.handleShowCommand(dto)
	}
}

func (m *MessageHandler) handleAddCommand(words []string, dto *Dto) (string, bool, error) {
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

	if trainingRepository.GetError() != nil {
		return dbErrorMessage, true, err
	}

	statisticRepository.AddStatistic(
		trainingEntity.GetId(),
		new(statistic.Count).Construct(countInt),
		new(user.Id).Construct(dto.GetUserId()),
	)

	if statisticRepository.GetError() != nil {
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

	stat := statisticRepository.GetByConditions(inputTrainingsAnyElems, periods)

	var wherePeriods []string
	var invalidPeriods []string

	for _, value := range periods {
		isValidText, err := m.validator.CheckIsOnlyRussianText(value)

		if err != nil {
			return "Произошла ошибка при проверке периода на валидность!", true, err
		}

		if isValidText {
			val, ok := conditionsByPeriod[value]
			if ok {
				wherePeriods = append(wherePeriods, val)
			} else {
				invalidPeriods = append(invalidPeriods, value)
			}
		} else {
			numsPeriod, invalidDatePeriods := m.prepareDateInterval(value)
			if len(invalidDatePeriods) > 0 {
				for _, invalid := range invalidDatePeriods {
					invalidPeriods = append(invalidPeriods, invalid)
				}
			} else {
				if m.sliceHelper.CheckLenSlice(numsPeriod, 1) {
					wherePeriods = append(wherePeriods, "DATE(`created`) = DATE('"+numsPeriod[firstSliceIndex]+"')")
				}

				if m.sliceHelper.CheckLenSlice(numsPeriod, 2) {
					wherePeriods = append(wherePeriods, "DATE(`created`) >= DATE('"+numsPeriod[firstSliceIndex]+"') AND DATE(`created`) <= DATE('"+numsPeriod[m.sliceHelper.SecondSliceElemIndex()]+"')")
				}
			}
		}
	}

	periodsSQL := strings.Join(wherePeriods, " OR ")
	log.Printf("Message info: %v \n", inputTrainings)
	log.Printf("Message info: %v \n", periodsSQL)
	query := `
				SELECT t.name as train, sum(count) as total, count(count) as sets 
				FROM statistic 
				JOIN training t on t.id = statistic.training_id 
				WHERE telegram_user_id = ? 
				AND ` + periodsSQL + ` 
 				AND t.name in (?` + strings.Repeat(",?", len(inputTrainings)-1) + `)
				GROUP BY training_id;
			`
	//fmt.Println(query)

	result, err := db.Query(
		query,
		append(append([]any{}, update.Message.From.ID), inputTrainingsAnyElems...)...,
	)

	log.Printf("Message info:  %v \n", result)
	if err != nil {
		sendMessage(bot, update, fmt.Sprintf("Произошла ошибка БД: %d ", err))
		continue
	}

	var results []ResultStatistic
	for result.Next() {
		var resultStruct ResultStatistic
		err = result.Scan(&resultStruct.train, &resultStruct.total, &resultStruct.sets)

		if err != nil {
			sendMessage(bot, update, fmt.Sprintf("Произошла ошибка БД: %d ", err))
		}

		results = append(results, resultStruct)
	}

	if m.sliceHelper.IsEmptySlice(results) {
		sendMessage(
			bot,
			update,
			"К сожалению по вашему запросу результаты не найдены. "+
				"Попробуйте указать другие упражнения и период.",
		)
	} else {
		resultMessage := "Вы сделали:\n"

		for _, result := range results {
			resultMessage += fmt.Sprintf("%v в количестве %v раз, за %v подхода(ов).", result.train, result.total, result.sets)
		}

		sendMessage(bot, update, resultMessage)
	}

	result.Close()
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
