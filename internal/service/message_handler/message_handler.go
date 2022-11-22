package message_handler

import (
	"bytes"
	"fmt"
	"sports-statistics/internal/config"
	cr "sports-statistics/internal/repositiry/command_repository"
	sr "sports-statistics/internal/repositiry/db/statistic"
	tr "sports-statistics/internal/repositiry/db/training"
	"sports-statistics/internal/repositiry/periods_repository"
	"sports-statistics/internal/service/entity/statistic"
	"sports-statistics/internal/service/entity/user"
	"sports-statistics/internal/service/helpers"
	"sports-statistics/internal/service/repository/command"
	"sports-statistics/internal/service/repository/periods"
	sri "sports-statistics/internal/service/repository/statistic"
	tri "sports-statistics/internal/service/repository/training"
	"strconv"
	"strings"
	"time"
)

const chatTypeGroup = "group"
const dbErrorMessage = "Произошла ошибка БД!"

type MessageHandler struct {
	validator           TrainingValidatorInterface
	sliceHelper         *helpers.SliceHelper
	stringHelper        *helpers.StringHelper
	commandRepository   command.RepositoryInterface
	trainingRepository  tri.RepositoryInterface
	statisticRepository sri.RepositoryInterface
	periodsRepository   periods.RepositoryInterface
}

func (m *MessageHandler) Construct() Handler {
	m.validator = new(TrainingValidator)
	m.sliceHelper = new(helpers.SliceHelper)
	m.stringHelper = new(helpers.StringHelper)
	m.commandRepository = new(cr.CommandRepository).Construct()
	m.trainingRepository = new(tr.Repository).Construct()
	m.statisticRepository = new(sr.Repository).Construct()
	m.periodsRepository = new(periods_repository.Repository).Construct()

	return m
}

func (m *MessageHandler) Destruct() {
	defer m.statisticRepository.Destruct()
	defer m.trainingRepository.Destruct()
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

	inputCommand := wordsFromMessText[firstSliceIndex]
	wordsFromMessText = m.sliceHelper.DeleteElemFromSlice(wordsFromMessText, firstSliceIndex)
	commandIsValid, err := m.validator.CheckIsOnlyRussianText(inputCommand)

	if err != nil {
		return "Произошла ошибка при проверке команды на валидность!", true, err
	}

	if !commandIsValid {
		return "Команда должна состоять только из русских букв.", true, nil
	}

	_, isAddCommand := m.commandRepository.GetAddCommands()[inputCommand]
	_, isShowCommand := m.commandRepository.GetShowCommands()[inputCommand]
	_, isHelpCommand := m.commandRepository.GetHelpCommands()[inputCommand]

	switch true {
	case isAddCommand:
		return m.handleAddCommand(dto, wordsFromMessText)
	case isShowCommand:
		return m.handleShowCommand(dto)
	case isHelpCommand:
		return m.handleHelpCommand(wordsFromMessText)
	}

	return "Не могу обработать введёную Вами команду.", true, nil
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

	trainingEntity := m.trainingRepository.GetTrainingByName(training)

	if err != nil {
		fmt.Println(err)
		return dbErrorMessage, true, err
	}

	m.statisticRepository.AddStatistic(
		new(statistic.Statistic).Construct(
			nil,
			trainingEntity,
			new(statistic.User).Construct(new(user.Id).Construct(dto.GetUserId())),
			nil,
			new(statistic.Count).Construct(countInt),
			nil,
		),
	)

	err = m.statisticRepository.GetError()
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

	rawPeriods := m.sliceHelper.DeleteElemFromSlice(inputByPeriod, firstSliceIndex)

	correctPeriods, _, err := m.prepareCorrectAndInvalidPeriods(rawPeriods)

	if err != nil {
		return "Произошла ошибка при проверке периода на валидность!", true, err
	}

	stat := m.statisticRepository.GetByConditions(inputTrainingsAnyElems, correctPeriods, dto.GetUserId())

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
func (m *MessageHandler) handleHelpCommand(words []string) (string, bool, error) {
	botName := config.Configs.GetBotName()
	addCommands := m.stringHelper.KeyMapToString(m.commandRepository.GetAddCommands(), "`, `")
	showCommands := m.stringHelper.KeyMapToString(m.commandRepository.GetShowCommands(), "`, `")

	message := fmt.Sprintf(
		"Привет! Я - бот, который поможет вести статистику спортивных упражнений, которые "+
			"ты выполняешь. Ты же ведь занимаешься спортом, верно?🤔\n"+
			"Я слушаю команды, когда ко мне обращаются. обратись ко мне вот так: `@%s`\n"+
			"Исключением является личная переписка. Если ты напишешь мне в личку, я буду реагировать на "+
			"любые твои сообщения. Но и в личных сообщениях поддерживается обращение, "+
			"если уж сильно хочется)\n"+
			"После обращения через пробел нужно написать команду и передать к ней данные, "+
			"Список поддерживаемых команд: \n"+
			"На добавление: `%s` \n"+
			"На показ статистики: `%s` \n"+
			"Чтобы посмотреть помощь по каждой комманде, отправь: `помощь` *название команды*\n"+
			"Например: `помощь добавить`",
		botName,
		addCommands,
		showCommands,
	)

	if !m.sliceHelper.IsEmptySlice(words) {
		commandToHelpMessage := words[m.sliceHelper.FirstSliceElemIndex()]

		_, isAddCommand := m.commandRepository.GetAddCommands()[commandToHelpMessage]
		_, isShowCommand := m.commandRepository.GetShowCommands()[commandToHelpMessage]
		_, isHelpCommand := m.commandRepository.GetHelpCommands()[commandToHelpMessage]

		switch true {
		case isAddCommand:
			allowTrainings := m.trainingRepository.GetTrainingNames()
			trainingBuf := bytes.Buffer{}

			for _, training := range allowTrainings {
				trainingBuf.WriteString(training.GetName().GetValue() + "`, `")
			}

			message = fmt.Sprintf(
				"Чтобы записать результаты, отправь команду на добавление упражнения (`%s`). Затем через "+
					"пробел укажи название упражнения, которое сделал. Далее через пробел укажи количество "+
					"повторений, которое сделал.\n"+
					"Например, ты сделал подход из 10 подтягиваний. Чтобы я всё корректно записал, напиши мне "+
					"`@%s сделал подтягивание 10`\n"+
					"Список доступных упражнений: `%s`",
				addCommands,
				botName,
				trainingBuf.String(),
			)
		case isShowCommand:
			periodsStr := strings.Join(m.periodsRepository.GetAllowTextPeriods(), "`, `")
			message = fmt.Sprintf(
				"Чтобы показать статистику, отправь команду `%s`. Затем укажи название упражнений, статистику "+
					"которых ты хочешь посмотреть. *Можно ввести несколько, разделив упражнения запятой,* "+
					"например, `подтягивание, отжимание`.\n"+
					"Далее укажи период, за который ты хочешь посмотреть статистику. Период будет корректно "+
					"распознан, если после указанных упражнений последует предлог *за*. Периодов можно указывать "+
					"несколько через запятую. Для каждого периода нужно так же нужен предлог *за*.\n"+
					"Например, нужно вывести статистику по подтягиваниям за сегодня, за 15.10.2022, "+
					"за период с 01.10.2022 по 10.10.2022. Чтобы периоды обработались корректно, введи периоды"+
					"следующим образом:\n"+
					"`за сегодня, за 15.10.2022, за 01.10.2022-10.10.2022`\n"+
					"Если период будет указан некорректно, результат будет без учёта некорректного периода. Если при "+
					"вводе интервала дата *от* окажется больше даты *до*, они поменяются местами и результат за этот "+
					"период будет найден корректно.\n"+
					"В итоге корректная команда будет выглядеть следующим образом: \n"+
					"`@%s покажи подтягивание, отжимание за сегодня, за 15.10.2022, за 01.10.2022-10.10.2022`\n"+
					"Список поддерживаемых текстовых периодов: `%s`",
				showCommands,
				botName,
				periodsStr,
			)
		case isHelpCommand:
			message = "Помощь к команде помощи не предусмотрена. " +
				"Надо ж было додуматься попросить помощь команде помощи🤔"
		}
	}

	return message, true, nil
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
	firstSliceIndex := m.sliceHelper.FirstSliceElemIndex()

	for _, period := range periods {
		isValidText, err := m.validator.CheckIsOnlyRussianText(period)

		if err != nil {
			return correctPeriods, invalidPeriods, err
		}

		if isValidText {
			val, ok := m.periodsRepository.GetConditionsByPeriod(period)
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
						m.periodsRepository.GetConditionsByDate(numsPeriod[firstSliceIndex]),
					)
				}

				if m.sliceHelper.CheckLenSlice(numsPeriod, 2) {
					correctPeriods = append(
						correctPeriods,
						m.periodsRepository.GetConditionsByDateInterval(
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
