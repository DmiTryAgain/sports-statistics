package tg

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/DmiTryAgain/sports-statistics/pkg/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vmkteam/embedlog"
)

var (
	errCantDetectLang  = errors.New("can't detect language")
	errCantRecognizeEx = errors.New("can't recognize the exercise")
)

var (
	enLangRe     = regexp.MustCompile(`^[a-zA-Z0-9\s.,!?'"@#$%^&*()\-_=+;:<>/\\|}{\[\]\p{So}]*$`)
	ruLangRe     = regexp.MustCompile(`^[а-яА-ЯёЁ0-9\s.,!?'"@#$%^&*()\-_=+;:<>/\\|}{\[\]\p{So}]*$`)
	valueUnitRe  = regexp.MustCompile(`^(\d+[.,]?\d*)\s*([a-zA-Zа-яА-ЯёЁ]+)$`)
	justNumberRe = regexp.MustCompile(`^(\d+[.,]?\d*)$`)
)

type MessageHandler struct {
	embedlog.Logger

	dbc      db.DB
	statRepo db.StatisticRepo
	tgBot    *tgbotapi.BotAPI
	cfg      Bot
}

func New(logger embedlog.Logger, db db.DB, statRepo db.StatisticRepo, tgBot *tgbotapi.BotAPI, cfg Bot) *MessageHandler {
	h := &MessageHandler{
		Logger:   logger,
		dbc:      db,
		cfg:      cfg,
		tgBot:    tgBot,
		statRepo: statRepo,
	}

	return h
}

func (m *MessageHandler) ListenAndHandle(ctx context.Context) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = int(m.cfg.Timeout.Duration.Seconds())

	// Get updates chan to listen to them
	updates := m.tgBot.GetUpdatesChan(updateConfig)

	limit := make(chan struct{}, 100)
	// Listen messages
	for upd := range updates {
		limit <- struct{}{} // Для каждого апдейта заполняем канал

		go func() {
			defer func() { <-limit }() // Освобождаем канал после успешного завершения горутины

			if upd.Message == nil {
				return
			}

			lowerText := strings.ToLower(upd.Message.Text)
			// Проверяем, что обращались вообще к нам
			hasMention := m.hasBotMention(lowerText)
			if !hasMention && upd.FromChat().IsGroup() {
				return // Скипаем, если к нам не обращались или не писали нам в личку
			}

			// Достаём пользователя
			userID := strconv.FormatInt(upd.Message.From.ID, 10)
			// Чистим текст от мусора
			msgText := m.clearRawMsg(lowerText)
			// Определяем язык
			lang, err := m.detectLang(msgText)
			if err != nil {
				m.Print(ctx, err.Error(), "msg", msgText, "userID", userID)
				m.sendMsg(upd, "Can't detect a language😶 Please, use the only one keyboard layout chars")
				return
			}

			// Если ничего не осталось, отправляем соответствующий ответ
			if msgText == "" {
				m.Print(ctx, "received empty message", "rawMsg", upd.Message.Text, "userID", userID)
				m.sendMsg(upd, messagesByLang[lang][emptyMessage])
				return
			}

			text, err := m.handle(ctx, msgText, userID, lang)
			if err != nil { // В случае ошибки сообщаем об этом
				text = messagesByLang[lang][errMsg]
				m.Error(ctx, "an error occurred", "rawMsg", upd.Message.Text, "userID", userID, "err", err.Error()) // И логируем её
			}

			m.sendMsg(upd, text)
		}()
	}
}

// sendMsg Отправляет сообщение
func (m *MessageHandler) sendMsg(upd tgbotapi.Update, text string) {
	if text == "" {
		return
	}

	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	msg.ReplyToMessageID = upd.Message.MessageID
	msg.ParseMode = m.cfg.ReplyFormat
	if _, err := m.tgBot.Send(msg); err != nil {
		// TODO: make retries
		m.Errorf("failed to send message: %v", err)
	}
}

// handle Обрабатывает сообщение. Определяет команду и обрабатывает остальной текст в соответствии с команлой
func (m *MessageHandler) handle(ctx context.Context, msgText, userID string, lang language) (string, error) {
	switch remainedText, c := m.detectCmd(msgText, lang); c {
	case addCmd:
		return m.handleAdd(ctx, remainedText, userID, lang)
	case showCmd:
		return m.handleShow(ctx, remainedText, userID, lang)
	case helpCmd:
		return m.handleHelp(ctx, remainedText, lang)
	default:
		m.Print(ctx, "received unknown command", "msg", msgText, "userID", userID)
		return fmt.Sprintf("%s. %s: %s", messagesByLang[lang][cantRecognizeCmd], messagesByLang[lang][listCmd], allCmdTextByLang(lang)), nil
	}
}

// hasBotMention Проверяет, был ли бот заменшенен
func (m *MessageHandler) hasBotMention(msgTxt string) bool {
	return strings.Contains(msgTxt, "@"+strings.ToLower(m.cfg.Name))
}

// detectLang Определяет язык по сообщению. В текущей реализации просто смотрит, на кириллице или латиннице был текст
func (m *MessageHandler) detectLang(msgTxt string) (language, error) {
	switch {
	case ruLangRe.MatchString(msgTxt):
		return langRU, nil
	case enLangRe.MatchString(msgTxt):
		return langEN, nil
	}

	return "", errCantDetectLang
}

// clearRawMsg Убирает из текста вызов бота, символы пункутации, переносы строк, пробелы по краям
func (m *MessageHandler) clearRawMsg(rawMsg string) string {
	// Убираем название бота
	withoutMention := strings.ReplaceAll(rawMsg, "@"+strings.ToLower(m.cfg.Name), "")

	const dashPlaceHolder = "DASHPLACEHOLDER"
	// Делаем специальный плейсхолдер с тире, чтобы не удалить лишние тире
	reHyphen := regexp.MustCompile(`(\d)\s*-\s*(\d)`)
	withPlaceHoder := reHyphen.ReplaceAllString(withoutMention, fmt.Sprintf("${1}%s${2}", dashPlaceHolder))

	const dotPlaceHolder = "DOTPLACEHOLDER"

	// Protect dots in dates (mark them temporary)
	reDateDot := regexp.MustCompile(`(\d{2})\.(\d{2})\.(\d{2}|\d{4})`)
	withPlaceHoder = reDateDot.ReplaceAllString(withPlaceHoder, fmt.Sprintf("${1}%s${2}%s${3}", dotPlaceHolder, dotPlaceHolder))

	// Replace dots in floats (exclude dates)
	reFloatDot := regexp.MustCompile(`(\d)\.(\d)`)
	withPlaceHoder = reFloatDot.ReplaceAllString(withPlaceHoder, fmt.Sprintf("${1}%s${2}", dotPlaceHolder))

	// Убираем символы пунктуации
	rePunct := regexp.MustCompile(`[[:punct:]]`)
	withoutPuncts := rePunct.ReplaceAllString(withPlaceHoder, "")

	// Заменяем все отступы и переносы строк на одиночный пробел
	reSpaces := regexp.MustCompile(`\s+`)
	withoutSpaces := reSpaces.ReplaceAllString(withoutPuncts, " ")

	// Теперь возвращаем тире обратно на место плейсхолдера
	withDashes := strings.ReplaceAll(withoutSpaces, dashPlaceHolder, "-")

	// Возвращаем точки обратно для дробных значений
	withDots := strings.ReplaceAll(withDashes, dotPlaceHolder, ".")

	// Убираем пробелы по краям и возвращаем
	return strings.TrimSpace(withDots)
}

// detectCmd Рассчитывает, какого типа команда, строку без названия команды и саму команду
func (m *MessageHandler) detectCmd(rawMsg string, lang language) (cleaned string, cmd cmd) {
	// Берём первое слово, чтобы понять, что за команда
	words := strings.SplitN(rawMsg, " ", 2)
	if len(words) == 0 {
		return rawMsg, unknownCmd
	}

	if len(words) > 1 {
		cleaned = words[1]
	}

	return cleaned, cmdByLang[lang][strings.ToLower(words[0])]
}

// parseValueWithUnit пытается извлечь числовое значение и единицу измерения из слова.
// Поддерживает форматы: "80кг" (слитно).
// Если слово — только число без суффикса, возвращает значение с hasUnit=false.
func parseValueWithUnit(word string, lang language) (value float64, unit UnitDef, hasUnit bool, err error) {
	// Пробуем слитный формат: число + суффикс
	if matches := valueUnitRe.FindStringSubmatch(word); len(matches) == 3 {
		numStr := strings.ReplaceAll(matches[1], ",", ".")
		value, err = strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0, UnitDef{}, false, fmt.Errorf("parse number %q: %w", matches[1], err)
		}
		suffix := strings.ToLower(matches[2])
		if u, ok := unitSuffixByLang[lang][suffix]; ok {
			return value, u, true, nil
		}
		return 0, UnitDef{}, false, fmt.Errorf("unknown unit suffix %q", suffix)
	}

	// Пробуем голое число
	if matches := justNumberRe.FindStringSubmatch(word); len(matches) == 2 {
		numStr := strings.ReplaceAll(matches[1], ",", ".")
		value, err = strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0, UnitDef{}, false, fmt.Errorf("parse number %q: %w", matches[1], err)
		}
		return value, UnitDef{}, false, nil
	}

	return 0, UnitDef{}, false, fmt.Errorf("can't parse %q as value", word)
}

var (
	errCountRequired    = errors.New("count required")
	errWeightRequired   = errors.New("weight required")
	errDistanceRequired = errors.New("distance required")
	errDurationRequired = errors.New("duration required")
)

func (pp *ParsedParams) addParam(pt ParamType, val float64) {
	switch pt {
	case ParamCount:
		if pp.Count == nil {
			pp.Count = ptrFloat64(val)
		} else {
			*pp.Count += val
		}
	case ParamWeight:
		if pp.WeightKg == nil {
			pp.WeightKg = ptrFloat64(val)
		} else {
			*pp.WeightKg += val
		}
	case ParamDistance:
		if pp.DistanceM == nil {
			pp.DistanceM = ptrFloat64(val)
		} else {
			*pp.DistanceM += val
		}
	case ParamDuration:
		if pp.DurationSec == nil {
			pp.DurationSec = ptrFloat64(val)
		} else {
			*pp.DurationSec += val
		}
	}
}

func (pp *ParsedParams) validate(category ExerciseCategory) error {
	for _, rp := range category.RequiredParams() {
		switch rp {
		case ParamCount:
			if pp.Count == nil {
				return errCountRequired
			}
		case ParamWeight:
			if pp.WeightKg == nil {
				return errWeightRequired
			}
		case ParamDistance:
			if pp.DistanceM == nil {
				return errDistanceRequired
			}
		case ParamDuration:
			if pp.DurationSec == nil {
				return errDurationRequired
			}
		}
	}
	return nil
}

// parseExerciseParams парсит слова после упражнения, извлекая параметры в соответствии с категорией.
func parseExerciseParams(words []string, category ExerciseCategory, lang language) (ParsedParams, error) {
	var pp ParsedParams

	for i := 0; i < len(words); i++ {
		value, unit, hasUnit, err := parseValueWithUnit(words[i], lang)
		if err != nil {
			continue
		}

		if hasUnit {
			pp.addParam(unit.ParamType, value*unit.Multiplier)
			continue
		}

		// Голое число — проверяем следующее слово как суффикс
		if i+1 < len(words) {
			nextWord := strings.ToLower(words[i+1])
			if u, ok := unitSuffixByLang[lang][nextWord]; ok {
				pp.addParam(u.ParamType, value*u.Multiplier)
				i++
				continue
			}
		}

		// Голое число без суффикса → count
		pp.addParam(ParamCount, value)
	}

	if err := pp.validate(category); err != nil {
		return pp, err
	}

	return pp, nil
}

func (m *MessageHandler) handleAdd(ctx context.Context, rawMsg, tgUserID string, lang language) (string, error) {
	lg := m.With("msg", rawMsg, "userID", tgUserID)
	if rawMsg == "" {
		lg.Print(ctx, "received empty message")
		return fmt.Sprintf("%s. %s: %s", messagesByLang[lang][emptyEx], messagesByLang[lang][listEx], allExTextByLang(lang)), nil
	}

	words := strings.Split(rawMsg, " ")

	// Достаём упражнение
	ex, found, position := m.extractExerciseAndItsPosition(words, lang)
	if !found {
		lg.Print(ctx, "received unknown exercise", "exercise", words[0])
		return fmt.Sprintf("%s: %s. %s: %s", messagesByLang[lang][cantRecognizeEx], words[0], messagesByLang[lang][listEx], allExTextByLang(lang)), nil
	}

	position++ // Нужно продолжить со следующего слова
	category := ex.Category()
	remainingWords := words[position:]

	if len(remainingWords) < 1 {
		lg.Print(ctx, "exercise requires params", "exercise", ex)
		return m.missingParamMessage(category, lang), nil
	}

	parsedParams, err := parseExerciseParams(remainingWords, category, lang)
	if err != nil {
		lg.Print(ctx, "param parsing error", "exercise", ex, "err", err)
		return m.missingParamMessage(category, lang), nil
	}

	cnt := parsedParams.CountOrDefault()
	if cnt < 1 {
		lg.Print(ctx, "invalid exercise count", "count", cnt)
		return messagesByLang[lang][cntGE], nil
	}

	_, err = m.statRepo.AddStatistic(ctx, &db.Statistic{
		TgUserID: tgUserID,
		Exercise: ex.String(),
		Count:    cnt,
		Params:   parsedParams.ToDBParams(),
		StatusID: 1,
	})
	if err != nil {
		return "", err
	}

	return messagesByLang[lang][exAdded], nil
}

func (m *MessageHandler) missingParamMessage(category ExerciseCategory, lang language) string {
	switch category {
	case CategoryRepsWeight:
		return messagesByLang[lang][weightRequired]
	case CategoryDistTime:
		return messagesByLang[lang][distanceRequired]
	case CategoryDuration:
		return messagesByLang[lang][durationRequired]
	default:
		return messagesByLang[lang][cntRequired]
	}
}

func (m *MessageHandler) handleShow(ctx context.Context, rawMsg, tgUserID string, lang language) (res string, err error) {
	if rawMsg == "" {
		m.Print(ctx, "received empty message", "msg", rawMsg, "userID", tgUserID)
		return fmt.Sprintf("%s. %s: %s", messagesByLang[lang][emptyEx], messagesByLang[lang][listEx], allExTextByLang(lang)), nil
	}

	// Парсим текст, находим упражнения и период
	exercises, periodsFilter, invPeriods, err := m.parseRawMsgAsExercisesAndPeriods(ctx, rawMsg, lang)
	if err != nil {
		if errors.Is(err, errCantRecognizeEx) {
			return fmt.Sprintf("%s. %s: %s", messagesByLang[lang][cantRecognizeEx], messagesByLang[lang][listEx], allExTextByLang(lang)), nil
		}

		return "", fmt.Errorf("parse raw message as exercises and periods, rawMsg=%s, err=%w", rawMsg, err)
	}

	// Сразу добавляем в результат нераспознаные периоды
	if len(invPeriods) != 0 {
		res += fmt.Sprintf("%s: %s\n", messagesByLang[lang][periodsInvalid], strings.Join(invPeriods, ", "))
	}

	// Теперь идём за статистикой
	s := db.GroupedStatisticSearch{
		StatisticSearch: db.StatisticSearch{
			TgUserID:  &tgUserID,
			Exercises: exercises.StringSlice(),
		},
		Periods: periodsFilter.ToDB(),
	}
	stats, err := m.statRepo.GroupedStatisticByFilters(ctx, s)
	if err != nil {
		return "", fmt.Errorf("fetch statistic, err=%w", err)
	}

	// Если ничего нет, выходим
	if len(stats) == 0 {
		return res + messagesByLang[lang][nothingFound], nil
	}

	table, err := m.buildTableByStat(ctx, stats, lang)
	if err != nil {
		return "", fmt.Errorf("build table by stat, err=%w", err)
	}

	res += table

	return res, nil
}

// parseRawMsgAsExercisesAndPeriods Парсит воходное сообщение без знаков пунктуации.
// Разбивает на слова, находит упражнение, затем парсит период
func (m *MessageHandler) parseRawMsgAsExercisesAndPeriods(ctx context.Context, rawMsg string, lang language) (exercises Exercises, periods periods, invalidPeriods []string, err error) {
	// Разбиваем по пробелам
	words := strings.Split(rawMsg, " ")

	var currentWord int // Здесь запомним, на каком элементе выйдем из цикла

	// Идём по каждому слову и ищем упражнения, который надо достать до первого фейла
	for currentWord < len(words) {
		if textContainsAllExerciseWords(words[currentWord], lang) {
			m.Print(ctx, "the message contains all exercises", "msg", rawMsg, "all exercises word", words[currentWord])
			currentWord++ // Пропускаем это слово, фильтр будет пустой, значит вытащим и так всё
			break
		}

		// Вытаскиваем упражнение учитывая, что оно моглобыть из нескольких слов
		ex, found, wordsLen := m.extractExerciseAndItsPosition(words[currentWord:], lang)
		if !found { // Упражнение состоит из одного слова, продолжаем перебирать слова
			currentWord++
			break
		}
		exercises = append(exercises, ex)
		currentWord += wordsLen + 1 // Сдвигаем на то кол-во слов, которое занимает это упражнение
	}

	// Проверяем, если вышли, и не нашли ни одного упражнения
	if len(exercises) == 0 && currentWord == 0 {
		return nil, nil, nil, errCantRecognizeEx
	}

	// Смотрим, есть ли кусок фразы про весь период в тексте.
	// Если нет, то парсим каждый период.
	// Если да, или если не задан, то считаем, что нужно взять за всё время.
	periodWords := words[currentWord:]
	if len(periodWords) > 0 {
		// Слепим оставшуюся подстроку под период
		periodLeftPart := strings.Join(periodWords, " ")
		// Если в ней нет спец фразы для всех упражнений
		if !textContainsAllPeriodWords(periodLeftPart, lang) {
			// То идём парсить каждый элеент
			periods, invalidPeriods = m.prepareCorrectAndInvalidPeriods(ctx, periodWords, lang)
		}
	}

	return
}

// extractExerciseAndItsPosition Достаёт из набора слов упражнение с учётом того, что оно может состоять:
// - Из одного слова.
// - Из двух и более слов, первое из которых уже является корректным упражнением.
// - Из двух и более слов, первое из которых не является корректным упражнением.
// Принимает набор слов и язык.
// Возвращает первое найденное упражнение и индекс его последнего слова из набора слов.
func (m *MessageHandler) extractExerciseAndItsPosition(words []string, lang language) (exercise Exercise, ok bool, exIdx int) {
	// Когда слов нет
	if len(words) == 0 {
		return
	}

	multiwordExName := words[exIdx]

	// Пробуем достать упражнение по первому слову
	exercise, ok = exerciseByLang[lang][multiwordExName]

	// Если оно было одно, его и вернём
	if len(words) == 1 {
		return
	}

	// Когда больше одного, сдвинемся до конца всех слов текущего упражнения
	for len(words) > exIdx+1 {
		multiwordExName = fmt.Sprintf("%s %s", multiwordExName, words[exIdx+1])
		multiwordEx, exists := exerciseByLang[lang][multiwordExName]
		if !exists && !ok { // Если не распознано, пробуем со следующим словом, но только если мы ещё не находили упражнение
			exIdx++
			continue
		}

		// Если оно реально состоит из 2х и более слов, снова сдвигаем i на следующее слово
		if exercise.isZero() || exercise == multiwordEx {
			ok = true
			exercise = multiwordEx
			exIdx++
			continue
		}

		// Останавливаемся, если упражнения различаются или если не нашли.
		// Мы захватили уже следующее или не найдено ни одного упражнения.
		break
	}

	return
}

func (m *MessageHandler) prepareCorrectAndInvalidPeriods(ctx context.Context, periodWords []string, lang language) (res periods, invalid []string) {
	// Проходимся по каждому периоду
	for i := range periodWords {
		// Скипаем предлоги
		if _, isPreposition := prepositionByLang[lang][periodWords[i]]; isPreposition {
			continue
		}

		// Если он текстовый
		isText := m.langReByLang(lang).MatchString(periodWords[i])

		// То обработаем, попробуем взять интервалы из текста
		if isText {
			p, ok := m.periodByText(periodWords[i], time.Now(), lang)
			if ok { // Если получилось, добавляем в результат
				res = append(res, p)
				continue
			}

			// Иначе добавляем в невалидные
			m.Print(ctx, "captured invalid text period", "period", periodWords[i])
			invalid = append(invalid, periodWords[i])
			continue
		}

		// Иначе это должны быть даты, обработаем их
		p, inv := m.periodByTime(ctx, periodWords[i])
		invalid = append(invalid, inv...)
		if !p.IsZero() {
			res = append(res, p)
		}
	}

	return
}

// langReByLang Возвращает регулярку для проверки фразы, что она состоит только из букв в текущем языке
func (m *MessageHandler) langReByLang(lang language) *regexp.Regexp {
	switch lang {
	case langRU:
		return regexp.MustCompile(`^[а-яА-ЯёЁ]+$`)
	case langEN:
		return regexp.MustCompile(`^[a-zA-Z]+$`)
	}

	return nil
}

func (m *MessageHandler) periodByText(text string, now time.Time, lang language) (p period, ok bool) {
	switch periodByLang[lang][text] {
	case todayPeriod:
		p = period{
			from: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			to:   time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC),
		}
		ok = true
	case yesterdayPeriod:
		p = period{
			from: time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.UTC),
			to:   time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
		}
		ok = true
	case dayBeforeYesterdayPeriod:
		p = period{
			from: time.Date(now.Year(), now.Month(), now.Day()-2, 0, 0, 0, 0, time.UTC),
			to:   time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.UTC),
		}
		ok = true
	case weekPeriod:
		// Получаем текущий день недели
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Вс считаем послдним днём, а не первым
		}
		// Отнимаем от текущего момента кол-во дней равное индексу дня недели. +1 нужно, чтобы считать с понедельника
		monday := now.AddDate(0, 0, -weekday+1)
		p = period{
			from: time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.UTC),
			to:   time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC),
		}
		ok = true
	case monthPeriod:
		p = period{
			from: time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC),
			to:   time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC),
		}
		ok = true
	case yearPeriod:
		p = period{
			from: time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC),
			to:   time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC),
		}
		ok = true
	}

	return
}

func (m *MessageHandler) periodByTime(ctx context.Context, interval string) (p period, invalid []string) {
	// И разибваем части по оставшемуся тире между датами
	intervals := strings.Split(interval, "-")

	// Если дата только одна, тогда from и to одинаковы
	if len(intervals) == 1 {
		t, err := m.parseDate(intervals[0])
		if err != nil {
			m.Print(ctx, "captured invalid single number period", "period", intervals[0])
			return period{}, []string{intervals[0]}
		}

		return period{
			from: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC),   // С начала дня
			to:   time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, time.UTC), // До следующего дня (т.к. до - не включительно)
		}, nil
	}

	// Если даты две
	if len(intervals) == 2 {
		from, err := m.parseDate(intervals[0])
		if err != nil {
			m.Print(ctx, "captured invalid interval", "period", intervals[0], "intervals", intervals)
			invalid = append(invalid, intervals[0])
		}

		to, err := m.parseDate(intervals[1])
		if err != nil {
			m.Print(ctx, "captured invalid interval", "period", intervals[1], "intervals", intervals)
			invalid = append(invalid, intervals[1])
		}

		// Меняем местами, если from был позже
		if from.After(to) {
			from, to = to, from
		}

		return period{
			from: time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC), // С начала дня
			to:   time.Date(to.Year(), to.Month(), to.Day()+1, 0, 0, 0, 0, time.UTC),     // До следующего дня (т.к. до - не включительно)
		}, invalid
	}

	return period{}, []string{interval}
}

func (m *MessageHandler) parseDate(date string) (time.Time, error) {
	// Пробуем распарсить в формате с полным годом
	parsed, err := time.Parse("02.01.2006", date)
	if err != nil {
		// Если не получилось, пробуем распарсить в формате с коротким годом
		parsed, err = time.Parse("02.01.06", date)
		if err != nil {
			return time.Time{}, err
		}
	}

	return parsed, nil
}

func (m *MessageHandler) buildTableByStat(ctx context.Context, in []db.GroupedStatistic, lang language) (string, error) {
	if len(in) == 0 {
		m.Print(ctx, "captured empty statistic list")
		return "", nil
	}

	tgStats := NewGroupedStatisticList(in, lang)

	hasWeight := anyHasWeight(tgStats)
	hasDistance := anyHasDistance(tgStats)
	hasDuration := anyHasDuration(tgStats)

	// Build header
	header := messagesByLang[lang][tableExCol]
	if hasWeight {
		header += "\t" + messagesByLang[lang][tableWeightCol]
	}
	if hasDistance {
		header += "\t" + messagesByLang[lang][tableDistCol]
	}
	if hasDuration {
		header += "\t" + messagesByLang[lang][tableTimeCol]
	}
	header += "\t" + messagesByLang[lang][tableCntCol]
	header += "\t" + messagesByLang[lang][tableSetCol]

	// Build rows
	var b strings.Builder
	b.WriteString("```\n")

	wr := tabwriter.NewWriter(&b, 0, 1, 4, ' ', 0)
	fmt.Fprintln(wr, header)

	for _, s := range tgStats {
		row := s.TranslatedExercise
		if hasWeight {
			if s.WeightKg != nil {
				row += "\t" + formatWeight(*s.WeightKg, lang)
			} else {
				row += "\t-"
			}
		}
		if hasDistance {
			if s.DistanceM != nil {
				row += "\t" + formatDistance(*s.DistanceM, lang)
			} else {
				row += "\t-"
			}
		}
		if hasDuration {
			if s.SumDurationSec != nil {
				row += "\t" + formatDuration(*s.SumDurationSec, lang)
			} else {
				row += "\t-"
			}
		}
		row += fmt.Sprintf("\t%g\t%d", s.SumCount, s.Sets)
		fmt.Fprintln(wr, row)
	}

	if err := wr.Flush(); err != nil {
		return "", err
	}

	b.WriteString("```")
	return b.String(), nil
}

func (m *MessageHandler) handleHelp(ctx context.Context, rawMsg string, lang language) (string, error) {
	switch _, c := m.detectCmd(rawMsg, lang); c {
	case unknownCmd:
		return fmt.Sprintf(messagesByLang[lang][commonHelpMsg], m.cfg.Name), nil
	case addCmd:
		return fmt.Sprintf(messagesByLang[lang][addHelpMsg], m.cfg.Name) +
			fmt.Sprintf("%s: %s", messagesByLang[lang][listEx], allExTextByLang(lang)), nil
	case showCmd:
		return fmt.Sprintf(messagesByLang[lang][showHelpMsg], m.cfg.Name) +
			fmt.Sprintf("%s: %s", messagesByLang[lang][listPeriod], allPeriodsByLang(lang)), nil
	case helpCmd:
		return messagesByLang[lang][helpHelpMsg], nil
	default:
		m.Print(ctx, "captured invalid command to show help", "msg", rawMsg, "command", c)
		return fmt.Sprintf("`%s` %s", rawMsg, messagesByLang[lang][cmdNotSupported]), nil
	}
}
