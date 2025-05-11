package tg

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/DmiTryAgain/sports-statistics/config"
	"github.com/DmiTryAgain/sports-statistics/pkg/db"
	"github.com/DmiTryAgain/sports-statistics/pkg/embedlog"
	"github.com/DmiTryAgain/sports-statistics/pkg/statistic"

	"github.com/go-pg/pg/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const chatTypeGroup = "group"
const dbErrorMessage = "Произошла ошибка БД!"

type MessageHandler struct {
	db          db.DB
	dbc         *pg.DB
	logger      embedlog.Logger
	tgBot       *tgbotapi.BotAPI
	statManager statistic.Manager
	cfg         config.Bot
}

func New(logger embedlog.Logger, db db.DB, dbc *pg.DB, tgBot *tgbotapi.BotAPI, cfg config.Bot) *MessageHandler {
	return &MessageHandler{
		db:          db,
		dbc:         dbc,
		cfg:         cfg,
		tgBot:       tgBot,
		logger:      logger,
		statManager: statistic.NewManager(db, logger),
	}
}

func (m *MessageHandler) ListenAndHandle() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = int(m.cfg.Timeout)

	// Get updates chan to listen to them
	updates := m.tgBot.GetUpdatesChan(updateConfig)

	// Listen messages
	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Now that we know we've gotten a new message, we can construct a
		// reply! We'll take the Chat ID and Text from the incoming message
		// and use it to create a new message.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// We'll also say that this message is a reply to the previous message.
		// For any other specifications than Chat ID or Text, you'll need to
		// set fields on the `MessageConfig`.
		msg.ReplyToMessageID = update.Message.MessageID

		// Okay, we're sending our message off! We don't care about the message
		// we just sent, so we'll discard it.
		if _, err := m.tgBot.Send(msg); err != nil {
			// Note that panics are a bad way to handle errors. Telegram can
			// have service outages or network errors, you should retry sending
			// messages or more gracefully handle failures.
			m.logger.Errorf("failed to send message: %v", err)
		}
	}
}

func (m *MessageHandler) Handle(upd tgbotapi.Update) (*Response, error) {
	// Проверяем, что образались вообще к нам
	hasMention := m.hasBotMention(upd.Message.Text)
	if !hasMention && upd.FromChat().IsGroup() {
		return nil, nil // Скипаем, если к нам не обращались или не писали нам в личку
	}

	msgText := m.clearRawMsg(upd.Message.Text)
	// Обрабатываем, если ничего не осталось
	if res := m.handleEmptyMessage(msgText); res != nil {
		return res, nil
	}

	switch m.evaluateCmd(msgText) {
	case addCmd:
		return m.handleAddCommand(msgText)
	case showCmd:
		return m.handleShowCommand()
	case helpCmd:
		return m.handleHelpCommand(msgText)
	default:
		return &Response{Message: "Не могу обработать введёную Вами команду"}, nil
	}
}

// hasBotMention Проверяет, был ли бот заменшенен
func (m *MessageHandler) hasBotMention(msgTxt string) bool {
	lowerBotName := strings.ToLower(m.cfg.Name)
	return strings.Contains(msgTxt, "@"+lowerBotName)
}

func (m *MessageHandler) handleEmptyMessage(msgTxt string) *Response {
	if msgTxt == "" {
		return &Response{Message: "Чё?"}
	}

	return nil
}

// clearRawMsg Убирает из текста вызов бота, если он есть, затем убирает переносы строк, и после - пробелы
func (m *MessageHandler) clearRawMsg(rawMsg string) string {
	clearMention := strings.ReplaceAll(rawMsg, strings.ToLower(m.cfg.Name), "")
	clearLines := strings.Trim(clearMention, string(filepath.Separator))
	return strings.Trim(clearLines, " ")
}

// evaluateCmd Рассчитывает, какого типа команда
func (m *MessageHandler) evaluateCmd(rawMsg string) cmd {
	// Берём первое слово, чтобы понять, что за команда
	words := strings.Split(rawMsg, " ")
	if len(words) < 2 {
		return unknownCmd
	}

	return cmdByWord[strings.ToLower(words[0])]
}

func (m *MessageHandler) handleAddCommand(msgText string) (*Response, error) {
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

func (m *MessageHandler) handleShowCommand(dto *message_handler.Dto) (string, bool, error) {
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
func (m *MessageHandler) handleHelpCommand(rawMsg string) (*Response, error) {
	message := fmt.Sprintf(
		"Привет! Я помогу тебе вести статистику твоих спортивных упражнений."+
			"Ты же ведь занимаешься спортом, верно?🤔\n"+
			"В группах обращайся ко мне вот так: `@%s`, чтобы я тебя слушал.\n"+
			"В личных сообщениях можно и без обращения, там я слушаю все сообщения."+
			"После обращения через пробел напиши команду."+
			"Список поддерживаемых команд: \n"+
			"На добавление: `Сделал` или `Добавь` \n"+
			"На показ статистики: `Покажи` \n"+
			"Чтобы посмотреть помощь по каждой комманде, отправь: `помощь` *название команды*\n"+
			"Например: `Помощь Добавь`",
		m.cfg.Name,
	)

	words := strings.Split(rawMsg, " ")
	if len(words) < 2 {
		return &Response{Message: message}, nil
	}

	switch m.evaluateCmd(words[0]) {
	case addCmd:
		message = fmt.Sprintf(
			"Чтобы записать результаты, отправь команду на добавление упражнения (`сделал`). Затем, через "+
				"пробел укажи упражнение, которое сделал. Далее через пробел укажи сделанное количество \n"+
				"Например, ты сделал подход из 10 подтягиваний. Чтобы я всё корректно записал, напиши мне "+
				"`@%s сделал подтягивание 10`\n"+
				"Список доступных упражнений: `%s`",
			m.cfg.Name,
			exercises().String(),
		)
	case showCmd:
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
	case helpCmd:
		message = "Помощь к команде помощи не предусмотрена. " +
			"Надо ж было додуматься попросить помощь команде помощи🤔"
	return &Response{Message: message}, nil
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
