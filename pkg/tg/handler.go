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

const maxMessageLen = 500

var (
	errCantDetectLang  = errors.New("can't detect language")
	errCantRecognizeEx = errors.New("can't recognize the exercise")
)

var (
	enLangRe     = regexp.MustCompile(`^[a-zA-Z0-9\s.,!?'"@#$%^&*()\-_=+;:<>/\\|}{\[\]\p{So}]*$`)
	ruLangRe     = regexp.MustCompile(`^[а-яА-ЯёЁ0-9\s.,!?'"@#$%^&*()\-_=+;:<>/\\|}{\[\]\p{So}]*$`)
	valueUnitRe  = regexp.MustCompile(`^(\d+[.,]?\d*)\s*([a-zA-Zа-яА-ЯёЁ]+)$`)
	justNumberRe = regexp.MustCompile(`^(\d+[.,]?\d*)$`)

	reHyphen   = regexp.MustCompile(`(\d)\s*-\s*(\d)`)
	reDateDot  = regexp.MustCompile(`(\d{2})\.(\d{2})\.(\d{2}|\d{4})`)
	reFloatDot = regexp.MustCompile(`(\d)\.(\d)`)
	rePunct    = regexp.MustCompile(`[[:punct:]]`)
	reSpaces   = regexp.MustCompile(`\s+`)

	ruLettersOnlyRe = regexp.MustCompile(`^[а-яА-ЯёЁ]+$`)
	enLettersOnlyRe = regexp.MustCompile(`^[a-zA-Z]+$`)
)

type MessageHandler struct {
	embedlog.Logger

	dbc       db.DB
	statRepo  db.StatisticRepo
	tgBot     *tgbotapi.BotAPI
	cfg       Bot
	sessions  *SessionStore
	rateLimit *rateLimiter
}

func New(shutdownCtx context.Context, logger embedlog.Logger, db db.DB, statRepo db.StatisticRepo, tgBot *tgbotapi.BotAPI, cfg Bot) *MessageHandler {
	h := &MessageHandler{
		Logger:    logger,
		dbc:       db,
		cfg:       cfg,
		tgBot:     tgBot,
		statRepo:  statRepo,
		sessions:  NewSessionStore(shutdownCtx),
		rateLimit: newRateLimiter(shutdownCtx),
	}

	return h
}

// ListenAndHandle Слушает апдейты из канала и обрабатывает входящие сообщения.
// shutdownCtx пока не используется, но возможно приготся позже для gracefully завершения работы, закрытия каналов и т.д.
func (m *MessageHandler) ListenAndHandle(shutdownCtx context.Context) {
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
			m.handleUpdate(upd)
		}()
	}
}

// handleUpdate обрабатывает один входящий апдейт
func (m *MessageHandler) handleUpdate(upd tgbotapi.Update) {
	ctx := context.Background()

	var (
		rateLimitID  int64
		fromCallback bool
	)
	if upd.Message != nil && upd.Message.From != nil {
		rateLimitID = upd.Message.From.ID
	} else if upd.CallbackQuery != nil {
		rateLimitID = upd.CallbackQuery.From.ID
		fromCallback = true
	}
	if rateLimitID != 0 && !m.rateLimit.Allow(rateLimitID) {
		m.Print(ctx, "same user exceeds the limit", "userID", rateLimitID, "fromCallback", fromCallback)
		return
	}

	// Обработка inline-кнопок
	if upd.CallbackQuery != nil {
		m.handleCallback(ctx, upd)
		return
	}

	if upd.Message == nil {
		return
	}

	// Обработка /start
	if upd.Message.IsCommand() && upd.Message.Command() == "start" {
		m.handleStart(ctx, upd)
		return
	}

	if len(upd.Message.Text) > maxMessageLen {
		m.sendMsg(upd, messagesByLang[langRU][msgTooLong])
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

	// Проверяем, не нажал ли пользователь ReplyKeyboard кнопку
	if c, lang, ok := m.detectReplyButton(upd.Message.Text); ok {
		m.handleReplyButton(ctx, upd, userID, lang, c)
		return
	}

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

	// Проверяем, есть ли активная сессия
	if session := m.sessions.Get(userID); session != nil && session.State != StateIdle {
		m.handleSessionText(ctx, upd, session, userID, msgText)
		return
	}

	text, err := m.handle(ctx, msgText, userID, lang)
	if err != nil { // В случае ошибки сообщаем об этом
		text = messagesByLang[lang][errMsg]
		m.Error(ctx, "an error occurred", "rawMsg", upd.Message.Text, "userID", userID, "err", err.Error()) // И логируем её
	}

	m.sendMsg(upd, text)
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
	withPlaceHoder := reHyphen.ReplaceAllString(withoutMention, fmt.Sprintf("${1}%s${2}", dashPlaceHolder))

	const dotPlaceHolder = "DOTPLACEHOLDER"

	// Protect dots in dates (mark them temporary)
	withPlaceHoder = reDateDot.ReplaceAllString(withPlaceHoder, fmt.Sprintf("${1}%s${2}%s${3}", dotPlaceHolder, dotPlaceHolder))

	// Replace dots in floats (exclude dates)
	withPlaceHoder = reFloatDot.ReplaceAllString(withPlaceHoder, fmt.Sprintf("${1}%s${2}", dotPlaceHolder))

	// Убираем символы пунктуации
	withoutPuncts := rePunct.ReplaceAllString(withPlaceHoder, "")

	// Заменяем все отступы и переносы строк на одиночный пробел
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
func parseValueWithUnit(word string, lang language) (ParsedValue, error) {
	// Пробуем слитный формат: число + суффикс
	if matches := valueUnitRe.FindStringSubmatch(word); len(matches) == 3 {
		numStr := strings.ReplaceAll(matches[1], ",", ".")
		v, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return ParsedValue{}, fmt.Errorf("parse number %q: %w", matches[1], err)
		}
		suffix := strings.ToLower(matches[2])
		if u, ok := unitSuffixByLang[lang][suffix]; ok {
			return ParsedValue{Value: v * u.Multiplier, Unit: &u}, nil
		}
		return ParsedValue{}, fmt.Errorf("unknown unit suffix %q", suffix)
	}

	// Пробуем голое число
	if matches := justNumberRe.FindStringSubmatch(word); len(matches) == 2 {
		numStr := strings.ReplaceAll(matches[1], ",", ".")
		v, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return ParsedValue{}, fmt.Errorf("parse number %q: %w", matches[1], err)
		}
		return ParsedValue{Value: v}, nil
	}

	return ParsedValue{}, fmt.Errorf("can't parse %q as value", word)
}

var (
	errCountRequired      = errors.New("count required")
	errWeightRequired     = errors.New("weight required")
	errDistanceRequired   = errors.New("distance required")
	errDurationRequired   = errors.New("duration required")
	errDistOrTimeRequired = errors.New("distance or duration required")
)

func (pp *ParsedParams) addParam(pt ParamType, val float64) {
	if cur := pp.GetParam(pt); cur != nil {
		pp.SetParam(pt, *cur+val)
		return
	}

	pp.SetParam(pt, val)
}

// parseExerciseParams парсит слова после упражнения, извлекая параметры в соответствии с категорией.
func parseExerciseParams(words []string, category ExerciseCategory, lang language) (ParsedParams, error) {
	var pp ParsedParams

	for i := 0; i < len(words); i++ {
		pv, err := parseValueWithUnit(words[i], lang)
		if err != nil {
			continue
		}

		if pv.HasUnit() {
			pp.addParam(pv.Unit.ParamType, pv.Value)
			continue
		}

		// Голое число — проверяем следующее слово как суффикс
		if i+1 < len(words) {
			nextWord := strings.ToLower(words[i+1])
			if u, ok := unitSuffixByLang[lang][nextWord]; ok {
				pp.addParam(u.ParamType, pv.Value*u.Multiplier)
				i++
				continue
			}
		}

		// Голое число без суффикса → count
		pp.addParam(ParamCount, pv.Value)
	}

	if err := category.ValidateParams(&pp); err != nil {
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
		StatusID: db.StatusEnabled,
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
		return messagesByLang[lang][distOrTimeRequired]
	case CategoryDuration:
		return messagesByLang[lang][durationRequired]
	case CategoryDurationWeight:
		return messagesByLang[lang][weightAndDurationRequired]
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

		if exists {
			// Более длинное совпадение всегда приоритетнее короткого
			ok = true
			exercise = multiwordEx
			exIdx++
			continue
		}

		if !ok {
			// Ещё не находили упражнение — пробуем со следующим словом
			exIdx++
			continue
		}

		// Уже нашли упражнение ранее, а более длинное не найдено — останавливаемся
		break
	}

	return
}

func (m *MessageHandler) prepareCorrectAndInvalidPeriods(ctx context.Context, periodWords []string, lang language) (res periods, invalid []string) {
	for i := 0; i < len(periodWords); {
		word := periodWords[i]

		// Скипаем предлоги
		if _, isPreposition := prepositionByLang[lang][word]; isPreposition {
			i++
			continue
		}

		// Если текстовое — пробуем многословный lookup (жадный, как для упражнений)
		if m.langReByLang(lang).MatchString(word) {
			p, consumed, ok := m.multiWordPeriodByText(periodWords[i:], time.Now(), lang)
			if ok {
				res = append(res, p)
				i += consumed
				continue
			}
			m.Print(ctx, "captured invalid text period", "period", word)
			invalid = append(invalid, word)
			i++
			continue
		}

		// Иначе это должны быть даты, обработаем их
		p, inv := m.periodByTime(ctx, word)
		invalid = append(invalid, inv...)
		if !p.IsZero() {
			res = append(res, p)
		}
		i++
	}

	return
}

// multiWordPeriodByText ищет наибольший совпадающий текстовый период, начиная с words[0].
// Возвращает период, количество потреблённых слов и флаг успеха.
func (m *MessageHandler) multiWordPeriodByText(words []string, now time.Time, lang language) (p period, consumed int, ok bool) {
	re := m.langReByLang(lang)
	candidate := ""
	bestMatch := -1
	var bestPeriod period

	for i, w := range words {
		if _, isPrep := prepositionByLang[lang][w]; isPrep {
			break
		}
		if !re.MatchString(w) {
			break
		}

		candidate = candidate + " " + w
		if i == 0 {
			candidate = w
		}

		if p2, isText := m.periodByText(candidate, now, lang); isText {
			bestMatch = i
			bestPeriod = p2
		}
	}

	if bestMatch >= 0 {
		return bestPeriod, bestMatch + 1, true
	}
	return period{}, 0, false
}

// langReByLang Возвращает регулярку для проверки фразы, что она состоит только из букв в текущем языке
func (m *MessageHandler) langReByLang(lang language) *regexp.Regexp {
	switch lang {
	case langRU:
		return ruLettersOnlyRe
	case langEN:
		return enLettersOnlyRe
	}

	return nil
}

func (m *MessageHandler) periodByText(text string, now time.Time, lang language) (p period, ok bool) {
	tp := periodByLang[lang][text]

	return m.periodByTextPeriod(tp, now)
}

func (m *MessageHandler) periodByTextPeriod(tp textPeriod, now time.Time) (p period, ok bool) {
	switch tp {
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
		// Отнимаем от текущего момента кол-во дней равное индексу дня недели. +1 нужно, чтобы считать с понедельника
		monday := now.AddDate(0, 0, -isoWeekday(now.Weekday())+1)
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

	case lastWeekPeriod:
		thisMonday := now.AddDate(0, 0, -isoWeekday(now.Weekday())+1)
		lastMonday := thisMonday.AddDate(0, 0, -7)
		p = period{
			from: time.Date(lastMonday.Year(), lastMonday.Month(), lastMonday.Day(), 0, 0, 0, 0, time.UTC),
			to:   time.Date(thisMonday.Year(), thisMonday.Month(), thisMonday.Day(), 0, 0, 0, 0, time.UTC),
		}
		ok = true
	case weekBeforeLastPeriod:
		thisMonday := now.AddDate(0, 0, -isoWeekday(now.Weekday())+1)
		lastMonday := thisMonday.AddDate(0, 0, -7)
		twoWeeksAgoMonday := thisMonday.AddDate(0, 0, -14)
		p = period{
			from: time.Date(twoWeeksAgoMonday.Year(), twoWeeksAgoMonday.Month(), twoWeeksAgoMonday.Day(), 0, 0, 0, 0, time.UTC),
			to:   time.Date(lastMonday.Year(), lastMonday.Month(), lastMonday.Day(), 0, 0, 0, 0, time.UTC),
		}
		ok = true
	case lastMonthPeriod:
		firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		p = period{
			from: firstOfThisMonth.AddDate(0, -1, 0),
			to:   firstOfThisMonth,
		}
		ok = true
	case monthBeforeLastPeriod:
		firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		p = period{
			from: firstOfThisMonth.AddDate(0, -2, 0),
			to:   firstOfThisMonth.AddDate(0, -1, 0),
		}
		ok = true
	case lastYearPeriod:
		p = period{
			from: time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, time.UTC),
			to:   time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC),
		}
		ok = true
	case yearBeforeLastPeriod:
		p = period{
			from: time.Date(now.Year()-2, 1, 1, 0, 0, 0, 0, time.UTC),
			to:   time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		ok = true

	default:
		p, ok = periodForWeekday(tp, now)
	}

	return
}

// isoWeekday конвертирует time.Weekday в ISO-номер: Mon=1 ... Sat=6, Sun=7.
func isoWeekday(w time.Weekday) int {
	if w == time.Sunday {
		return 7
	}
	return int(w)
}

func periodForWeekday(tp textPeriod, now time.Time) (period, bool) {
	targetWeekday, ok := weekdayByPeriod[tp]
	if !ok {
		return period{}, false
	}

	targetISO := isoWeekday(targetWeekday)
	currentISO := isoWeekday(now.Weekday())

	var targetDay time.Time
	switch {
	case targetISO == currentISO:
		targetDay = now
	case targetISO < currentISO:
		targetDay = now.AddDate(0, 0, -(currentISO - targetISO))
	default:
		targetDay = now.AddDate(0, 0, -(7 - (targetISO - currentISO)))
	}

	from := time.Date(targetDay.Year(), targetDay.Month(), targetDay.Day(), 0, 0, 0, 0, time.UTC)

	to := from.AddDate(0, 0, 1)
	if targetISO == currentISO {
		to = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC)
	}
	return period{from: from, to: to}, true
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
			if s.DurationSec != nil {
				row += "\t" + formatDuration(*s.DurationSec, lang)
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

// detectReplyButton проверяет, совпадает ли текст с одной из ReplyKeyboard кнопок.
// Возвращает команду, язык и true, если совпадение найдено.
func (m *MessageHandler) detectReplyButton(text string) (cmd, language, bool) {
	for lang, buttons := range replyButtonCmd {
		if c, ok := buttons[text]; ok {
			return c, lang, true
		}
	}
	return unknownCmd, "", false
}

// detectLangFromTgUser определяет язык по настройкам Telegram-пользователя.
// Русский — для ru/uk/be/kk, английский — для всех остальных.
func detectLangFromTgUser(user *tgbotapi.User) language {
	if user == nil {
		return langRU
	}
	switch user.LanguageCode {
	case "ru", "uk", "be", "kk":
		return langRU
	default:
		return langEN
	}
}

// handleStart обрабатывает команду /start — отправляет приветствие и ReplyKeyboard
func (m *MessageHandler) handleStart(ctx context.Context, upd tgbotapi.Update) {
	lang := detectLangFromTgUser(upd.Message.From)

	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, messagesByLang[lang][welcomeMsg])
	kb := replyKeyboard(lang)
	kb.ResizeKeyboard = true
	msg.ReplyMarkup = kb
	msg.ParseMode = m.cfg.ReplyFormat
	if _, err := m.tgBot.Send(msg); err != nil {
		m.Error(ctx, "failed to send start message", "err", err)
	}
}

// handleReplyButton обрабатывает нажатие ReplyKeyboard кнопки
func (m *MessageHandler) handleReplyButton(ctx context.Context, upd tgbotapi.Update, userID string, lang language, c cmd) {
	switch c {
	case addCmd:
		m.showAddExerciseScreen(ctx, upd, userID, lang)
	case showCmd:
		m.showStatExerciseScreen(ctx, upd, userID, lang)
	case helpCmd:
		text := fmt.Sprintf(messagesByLang[lang][commonHelpMsg], m.cfg.Name)
		m.sendMsg(upd, text)
	}
}

// showAddExerciseScreen показывает экран выбора упражнения с Quick-Add кнопками
func (m *MessageHandler) showAddExerciseScreen(ctx context.Context, upd tgbotapi.Update, userID string, lang language) {
	// Получаем частые упражнения пользователя
	frequent, err := m.statRepo.FrequentExercisesByUser(ctx, userID, 5)
	if err != nil {
		m.Error(ctx, "fetch frequent exercises", "err", err)
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	// Quick-Add кнопки
	qaRows := quickAddInlineKeyboard(frequent, lang)
	if len(qaRows) > 0 {
		rows = append(rows, qaRows...)
	}

	// Кнопки упражнений (первая страница)
	exKb := exerciseInlineKeyboard(0, lang, "add")
	rows = append(rows, exKb.InlineKeyboard...)

	kb := tgbotapi.NewInlineKeyboardMarkup(rows...)

	text := messagesByLang[lang][chooseExerciseOrText]
	if len(qaRows) > 0 {
		text = messagesByLang[lang][yourFrequent] + "\n" + text
	}

	msgID := m.sendMsgWithKeyboard(upd, text, kb)

	session := &UserSession{
		State:            StateAwaitExercise,
		Lang:             lang,
		LastBotMessageID: msgID,
		ChatID:           upd.Message.Chat.ID,
	}
	m.sessions.Set(userID, session)
}

// showStatExerciseScreen показывает экран выбора упражнения для статистики
func (m *MessageHandler) showStatExerciseScreen(ctx context.Context, upd tgbotapi.Update, userID string, lang language) {
	exercises, err := m.statRepo.UniqueExercisesByUser(ctx, userID)
	if err != nil {
		m.Error(ctx, "fetch unique exercises", "err", err)
		m.sendMsg(upd, messagesByLang[lang][errMsg])
		return
	}

	if len(exercises) == 0 {
		m.sendMsg(upd, messagesByLang[lang][nothingFound])
		return
	}

	kb := showExerciseInlineKeyboard(exercises, lang)
	text := messagesByLang[lang][chooseExercise]
	msgID := m.sendMsgWithKeyboard(upd, text, kb)

	session := &UserSession{
		State:            StateShowAwaitExercise,
		Lang:             lang,
		LastBotMessageID: msgID,
		ChatID:           upd.Message.Chat.ID,
	}
	m.sessions.Set(userID, session)
}

// handleCallback обрабатывает нажатие inline-кнопки
func (m *MessageHandler) handleCallback(ctx context.Context, upd tgbotapi.Update) {
	cb := upd.CallbackQuery
	userID := strconv.FormatInt(cb.From.ID, 10)

	action, err := parseCallbackData(cb.Data)
	if err != nil {
		m.Error(ctx, "parse callback data", "data", cb.Data, "err", err)
		m.answerCallback(cb.ID, "")
		return
	}

	session := m.sessions.Get(userID)
	if session == nil {
		switch action.Type {
		case cbQuickAdd, cbCancel, cbExercisePage:
			session = &UserSession{Lang: detectLangFromTgUser(cb.From)}
		default:
			m.answerCallback(cb.ID, messagesByLang[detectLangFromTgUser(cb.From)][sessionExpired])
			return
		}
	}

	switch action.Type {
	case cbSelectExercise:
		m.handleCBSelectExercise(ctx, cb, session, userID, action)
	case cbSelectWeight:
		m.handleCBSelectWeight(ctx, cb, session, userID, action)
	case cbSelectCount:
		m.handleCBSelectCount(ctx, cb, session, userID, action)
	case cbSelectDistance:
		m.handleCBSelectDistance(ctx, cb, session, userID, action)
	case cbSelectDuration:
		m.handleCBSelectDuration(ctx, cb, session, userID, action)
	case cbCustomInput:
		m.handleCBCustomInput(ctx, cb, session, userID, action)
	case cbQuickAdd:
		m.handleCBQuickAdd(ctx, cb, session, userID, action)
	case cbShowExercise:
		m.handleCBShowExercise(ctx, cb, session, userID, action)
	case cbShowAll:
		m.handleCBShowAll(ctx, cb, session, userID)
	case cbShowPeriod:
		m.handleCBShowPeriod(ctx, cb, session, userID, action)
	case cbSkipParam:
		m.handleCBSkipParam(ctx, cb, session, userID, action)
	case cbCancel:
		m.handleCBCancel(ctx, cb, session, userID)
	case cbExercisePage:
		m.handleCBExercisePage(ctx, cb, session, userID, action)
	default:
		m.answerCallback(cb.ID, "")
	}
}

// handleCBSelectExercise обрабатывает выбор упражнения
func (m *MessageHandler) handleCBSelectExercise(ctx context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, userID string, action CallbackAction) {
	session.Exercise = action.Exercise
	m.sessions.Set(userID, session)
	m.answerCallback(cb.ID, "")
	m.advanceAddDialog(ctx, cb, session, userID)
}

// handleCBSelectWeight обрабатывает выбор веса
func (m *MessageHandler) handleCBSelectWeight(ctx context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, userID string, action CallbackAction) {
	session.Params.WeightKg = ptr(action.Value)
	m.sessions.Set(userID, session)
	m.answerCallback(cb.ID, "")
	m.advanceAddDialog(ctx, cb, session, userID)
}

// handleCBSelectCount обрабатывает выбор количества повторений
func (m *MessageHandler) handleCBSelectCount(ctx context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, userID string, action CallbackAction) {
	session.Params.Count = ptr(action.Value)
	m.sessions.Set(userID, session)
	m.answerCallback(cb.ID, "")
	m.advanceAddDialog(ctx, cb, session, userID)
}

// handleCBSelectDistance обрабатывает выбор дистанции
func (m *MessageHandler) handleCBSelectDistance(ctx context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, userID string, action CallbackAction) {
	session.Params.DistanceM = ptr(action.Value)
	m.sessions.Set(userID, session)
	m.answerCallback(cb.ID, "")
	m.advanceAddDialog(ctx, cb, session, userID)
}

// handleCBSelectDuration обрабатывает выбор длительности
func (m *MessageHandler) handleCBSelectDuration(ctx context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, userID string, action CallbackAction) {
	session.Params.DurationSec = ptr(action.Value)
	m.sessions.Set(userID, session)
	m.answerCallback(cb.ID, "")
	m.advanceAddDialog(ctx, cb, session, userID)
}

// handleCBCustomInput обрабатывает нажатие кнопки "Другой"
func (m *MessageHandler) handleCBCustomInput(_ context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, userID string, action CallbackAction) {
	session.State = StateAwaitCustomValue
	session.CustomTarget = action.CustomTarget
	m.sessions.Set(userID, session)

	var text string
	switch action.CustomTarget {
	case TargetWeight:
		text = messagesByLang[session.Lang][enterCustomWeight]
	case TargetCount:
		text = messagesByLang[session.Lang][enterCustomCount]
	case TargetDistance:
		text = messagesByLang[session.Lang][enterCustomDistance]
	case TargetDuration:
		text = messagesByLang[session.Lang][enterCustomDuration]
	}

	m.editMessageRemoveKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, text)
	m.answerCallback(cb.ID, "")
}

// handleCBQuickAdd обрабатывает Quick-Add кнопку — сразу сохраняем
func (m *MessageHandler) handleCBQuickAdd(ctx context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, userID string, action CallbackAction) {
	_, err := m.statRepo.AddStatistic(ctx, &db.Statistic{
		TgUserID: userID,
		Exercise: action.Exercise.String(),
		Count:    action.Count,
		Params:   action.Params,
		StatusID: db.StatusEnabled,
	})
	if err != nil {
		m.Error(ctx, "quick add failed", "err", err)
		m.answerCallback(cb.ID, messagesByLang[session.Lang][errMsg])
		return
	}

	confirmText := formatAddConfirmation(action.Exercise, action.Count, action.Params, session.Lang)
	m.editMessageRemoveKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, confirmText)
	m.answerCallback(cb.ID, messagesByLang[session.Lang][exAdded])
	m.sessions.Delete(userID)
}

// handleCBShowExercise обрабатывает выбор упражнения для статистики
func (m *MessageHandler) handleCBShowExercise(_ context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, userID string, action CallbackAction) {
	session.State = StateShowAwaitPeriod
	session.ShowExercises = Exercises{action.Exercise}
	m.sessions.Set(userID, session)

	exName := exTextByLang[session.Lang][action.Exercise]
	if exName == "" {
		exName = action.Exercise.String()
	}
	text := fmt.Sprintf(messagesByLang[session.Lang][choosePeriod], exName)
	kb := periodInlineKeyboard(session.Lang)
	m.editMessageWithKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, text, kb)
	m.answerCallback(cb.ID, "")
}

// handleCBShowAll обрабатывает кнопку "Всё" для статистики
func (m *MessageHandler) handleCBShowAll(_ context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, userID string) {
	session.State = StateShowAwaitPeriod
	session.ShowExercises = nil // nil = все упражнения
	m.sessions.Set(userID, session)

	text := fmt.Sprintf(messagesByLang[session.Lang][choosePeriod], messagesByLang[session.Lang][allExBtn])
	kb := periodInlineKeyboard(session.Lang)
	m.editMessageWithKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, text, kb)
	m.answerCallback(cb.ID, "")
}

// handleCBShowPeriod обрабатывает выбор периода для статистики
func (m *MessageHandler) handleCBShowPeriod(ctx context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, userID string, action CallbackAction) {
	p, _ := m.periodByTextPeriod(action.Period, time.Now())

	var periodsFilter periods
	if !p.IsZero() {
		periodsFilter = periods{p}
	}

	s := db.GroupedStatisticSearch{
		StatisticSearch: db.StatisticSearch{
			TgUserID:  &userID,
			Exercises: session.ShowExercises.StringSlice(),
		},
		Periods: periodsFilter.ToDB(),
	}

	stats, err := m.statRepo.GroupedStatisticByFilters(ctx, s)
	if err != nil {
		m.Error(ctx, "fetch statistic", "err", err)
		m.answerCallback(cb.ID, messagesByLang[session.Lang][errMsg])
		return
	}

	if len(stats) == 0 {
		m.editMessageRemoveKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, messagesByLang[session.Lang][nothingFound])
		m.answerCallback(cb.ID, "")
		m.sessions.Delete(userID)
		return
	}

	table, err := m.buildTableByStat(ctx, stats, session.Lang)
	if err != nil {
		m.Error(ctx, "build table", "err", err)
		m.answerCallback(cb.ID, messagesByLang[session.Lang][errMsg])
		return
	}

	m.editMessageRemoveKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, table)
	m.answerCallback(cb.ID, "")
	m.sessions.Delete(userID)
}

// handleCBCancel обрабатывает кнопку отмены
func (m *MessageHandler) handleCBCancel(_ context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, userID string) {
	m.editMessageRemoveKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, messagesByLang[session.Lang][cancelBtn])
	m.answerCallback(cb.ID, "")
	m.sessions.Delete(userID)
}

// handleCBSkipParam обрабатывает пропуск необязательного параметра
func (m *MessageHandler) handleCBSkipParam(ctx context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, userID string, action CallbackAction) {
	paramType := customTargetToParamType(action.CustomTarget)
	if session.SkippedParams == nil {
		session.SkippedParams = make(map[ParamType]struct{})
	}
	session.SkippedParams[paramType] = struct{}{}
	m.sessions.Set(userID, session)
	m.answerCallback(cb.ID, "")
	m.advanceAddDialog(ctx, cb, session, userID)
}

// handleCBExercisePage обрабатывает переключение страницы упражнений
func (m *MessageHandler) handleCBExercisePage(_ context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, _ string, action CallbackAction) {
	kb := exerciseInlineKeyboard(action.Page, session.Lang, action.Context)
	m.editMessageWithKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, messagesByLang[session.Lang][chooseExercise], kb)
	m.answerCallback(cb.ID, "")
}

// nextStep описывает следующий незаполненный параметр в диалоге добавления
type nextStep struct {
	Param    ParamType
	Optional bool // показывать кнопку «Пропустить»
}

// nextUnfilledParam определяет следующий незаполненный параметр (required → soft-required → optional).
// Возвращает nil, если все параметры заполнены.
func nextUnfilledParam(exercise Exercise, params *ParsedParams, skipped map[ParamType]struct{}) *nextStep {
	cat := exercise.Category()

	// 1. Обязательные
	for _, rp := range cat.RequiredParams() {
		if params.GetParam(rp) == nil {
			return &nextStep{Param: rp, Optional: false}
		}
	}

	// 2. Soft-required (хотя бы один из списка)
	for _, sp := range cat.SoftRequiredParams() {
		if _, ok := skipped[sp]; ok {
			continue
		}
		if params.GetParam(sp) == nil {
			optional := sp == ParamDuration && params.DistanceM != nil
			return &nextStep{Param: sp, Optional: optional}
		}
	}

	// 3. Необязательные
	for _, op := range exercise.OptionalParams() {
		if _, ok := skipped[op]; ok {
			continue
		}
		if params.GetParam(op) == nil {
			return &nextStep{Param: op, Optional: true}
		}
	}

	return nil
}

// advanceAddDialog проверяет, какие параметры ещё не заполнены, и показывает следующий шаг
func (m *MessageHandler) advanceAddDialog(ctx context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, userID string) {
	step := nextUnfilledParam(session.Exercise, &session.Params, session.SkippedParams)
	if step == nil {
		m.saveFromSession(ctx, cb, session, userID)
		return
	}
	m.showParamStepCB(ctx, cb, session, userID, step)
}

// showParamStepCB показывает шаг выбора параметра, редактируя callback-сообщение
func (m *MessageHandler) showParamStepCB(ctx context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, userID string, step *nextStep) {
	exName := exTextByLang[session.Lang][session.Exercise]
	lang := session.Lang
	cat := session.Exercise.Category()

	switch step.Param {
	case ParamWeight:
		session.State = StateAwaitWeight
		m.sessions.Set(userID, session)
		weights, err := m.statRepo.UniqueWeightsByExercise(ctx, userID, session.Exercise.String(), 5)
		if err != nil {
			m.Error(ctx, "fetch unique weights", "err", err)
		}
		if step.Optional {
			text := fmt.Sprintf(messagesByLang[lang][chooseOptionalWeight], exName) + "\n" + messagesByLang[lang][orWriteText]
			m.editMessageWithKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, text, optionalWeightInlineKeyboard(weights, lang))
		} else {
			text := fmt.Sprintf(messagesByLang[lang][chooseWeight], exName) + "\n" + messagesByLang[lang][orWriteText]
			m.editMessageWithKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, text, weightInlineKeyboard(weights, lang))
		}

	case ParamCount:
		session.State = StateAwaitCount
		m.sessions.Set(userID, session)
		var detail string
		if session.Params.WeightKg != nil {
			detail = " " + formatWeight(*session.Params.WeightKg, lang)
		}
		text := fmt.Sprintf(messagesByLang[lang][chooseCount], exName+detail) + "\n" + messagesByLang[lang][orWriteText]
		m.editMessageWithKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, text, countInlineKeyboard(lang))

	case ParamDistance:
		session.State = StateAwaitDistance
		m.sessions.Set(userID, session)
		if step.Optional {
			text := fmt.Sprintf(messagesByLang[lang][chooseOptionalDistance], exName) + "\n" + messagesByLang[lang][orWriteText]
			m.editMessageWithKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, text, optionalDistanceInlineKeyboard(lang))
		} else {
			text := fmt.Sprintf(messagesByLang[lang][chooseDistance], exName) + "\n" + messagesByLang[lang][orWriteText]
			m.editMessageWithKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, text, distanceInlineKeyboard(lang))
		}

	case ParamDuration:
		session.State = StateAwaitDuration
		m.sessions.Set(userID, session)
		var detail string
		if session.Params.DistanceM != nil {
			detail = " " + formatDistance(*session.Params.DistanceM, lang)
		}
		if step.Optional {
			text := fmt.Sprintf(messagesByLang[lang][chooseOptionalDuration], exName) + "\n" + messagesByLang[lang][orWriteText]
			m.editMessageWithKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, text, optionalDurationInlineKeyboard(cat, lang))
		} else {
			text := fmt.Sprintf(messagesByLang[lang][chooseDuration], exName+detail) + "\n" + messagesByLang[lang][orWriteText]
			m.editMessageWithKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, text, durationInlineKeyboard(cat, lang))
		}
	}
}

// saveFromSession сохраняет упражнение из сессии в БД
func (m *MessageHandler) saveFromSession(ctx context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, userID string) {
	cnt := session.Params.CountOrDefault()

	_, err := m.statRepo.AddStatistic(ctx, &db.Statistic{
		TgUserID: userID,
		Exercise: session.Exercise.String(),
		Count:    cnt,
		Params:   session.Params.ToDBParams(),
		StatusID: db.StatusEnabled,
	})
	if err != nil {
		m.Error(ctx, "save from session failed", "err", err)
		m.answerCallback(cb.ID, messagesByLang[session.Lang][errMsg])
		return
	}

	confirmText := formatAddConfirmation(session.Exercise, cnt, session.Params.ToDBParams(), session.Lang)
	m.editMessageRemoveKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, confirmText)
	m.answerCallback(cb.ID, messagesByLang[session.Lang][exAdded])
	m.sessions.Delete(userID)
}

// handleSessionText обрабатывает текстовый ввод в контексте активной сессии
func (m *MessageHandler) handleSessionText(ctx context.Context, upd tgbotapi.Update, session *UserSession, userID, text string) {
	switch session.State {
	case StateAwaitExercise:
		// Пытаемся распарсить как полное сообщение добавления
		result, err := m.handleAdd(ctx, text, userID, session.Lang)
		if err != nil {
			m.sendMsg(upd, messagesByLang[session.Lang][errMsg])
			m.Error(ctx, "session handleAdd", "err", err, "text", text, "user", userID, "session", session.String())
		} else {
			m.sendMsg(upd, result)
		}
		m.sessions.Delete(userID)

	case StateAwaitCustomValue:
		m.handleCustomValueInput(ctx, upd, session, userID, text)

	case StateAwaitWeight, StateAwaitCount, StateAwaitDistance, StateAwaitDuration:

		if err := m.handleRemainingParamsText(ctx, upd, session, userID, text); err != nil {
			m.Error(ctx, "session handleRemainingParamsText", "err", err, "text", text, "user", userID, "session", session.String())
			m.sendMsg(upd, messagesByLang[session.Lang][errMsg])
		}

	default:
		// Нет обработки — fallback на обычную обработку
		result, err := m.handle(ctx, text, userID, session.Lang)
		if err != nil {
			m.sendMsg(upd, messagesByLang[session.Lang][errMsg])
			m.Error(ctx, "session fallback handle", "err", err, "text", text, "user", userID, "session", session.String())
		} else if result != "" {
			m.sendMsg(upd, result)
		}
	}
}

// handleCustomValueInput обрабатывает ручной ввод значения для кнопки "Другой"
func (m *MessageHandler) handleCustomValueInput(ctx context.Context, upd tgbotapi.Update, session *UserSession, userID, text string) {
	pv, err := parseValueWithUnit(text, session.Lang)
	if err != nil {
		m.sendMsg(upd, fmt.Sprintf(messagesByLang[session.Lang][paramInvalid], text))
		return
	}

	if pv.HasUnit() {
		session.Params.addParam(pv.Unit.ParamType, pv.Value)
	} else {
		// Голое число — записываем в параметр, соответствующий текущему CustomTarget
		session.Params.SetParam(customTargetToParamType(session.CustomTarget), pv.Value)
	}

	m.sessions.Set(userID, session)

	// Показываем следующий шаг или сохраняем, если все параметры заполнены.
	// editMsgID=0: клавиатура уже убрана кнопкой «Другой», отправляем новое сообщение.
	if !m.showNextParamAsNewMessage(ctx, upd, session, userID, 0, 0) {
		m.saveFromSessionText(ctx, upd, session, userID)
	}
}

// handleRemainingParamsText обрабатывает ввод параметра текстом в контексте ожидания кнопки
func (m *MessageHandler) handleRemainingParamsText(ctx context.Context, upd tgbotapi.Update, session *UserSession, userID, text string) error {
	pv, err := parseValueWithUnit(text, session.Lang)
	if err != nil {
		m.sendMsg(upd, fmt.Sprintf(messagesByLang[session.Lang][paramInvalid], text))
		// Игнорируем эту ошибку, просто отправим пользователю, что некорректные единицы измерения
		//nolint:nilerr
		return nil
	}

	if pv.HasUnit() {
		session.Params.addParam(pv.Unit.ParamType, pv.Value)
	} else {
		// Голое число — применяем в зависимости от текущего состояния
		pt, ok := session.State.ParamType()
		if !ok {
			return fmt.Errorf("unknown state: %d", session.State)
		}

		session.Params.SetParam(pt, pv.Value)
	}

	m.sessions.Set(userID, session)

	// Редактируем старое сообщение с клавиатурой (если есть) вместо отправки нового
	chatID := upd.Message.Chat.ID
	editMsgID := session.LastBotMessageID

	if !m.showNextParamAsNewMessage(ctx, upd, session, userID, chatID, editMsgID) {
		// Все параметры заполнены — убираем старую клавиатуру и сохраняем
		if editMsgID > 0 {
			confirmText := formatAddConfirmation(session.Exercise, session.Params.CountOrDefault(), session.Params.ToDBParams(), session.Lang)
			m.editMessageRemoveKeyboard(chatID, editMsgID, confirmText)
		}
		m.saveFromSessionText(ctx, upd, session, userID)
	}

	return nil
}

// saveFromSessionText сохраняет из сессии при текстовом вводе
func (m *MessageHandler) saveFromSessionText(ctx context.Context, upd tgbotapi.Update, session *UserSession, userID string) {
	cnt := session.Params.CountOrDefault()

	_, err := m.statRepo.AddStatistic(ctx, &db.Statistic{
		TgUserID: userID,
		Exercise: session.Exercise.String(),
		Count:    cnt,
		Params:   session.Params.ToDBParams(),
		StatusID: db.StatusEnabled,
	})
	if err != nil {
		m.Error(ctx, "save from session text failed", "err", err)
		m.sendMsg(upd, messagesByLang[session.Lang][errMsg])
		return
	}

	confirmText := formatAddConfirmation(session.Exercise, cnt, session.Params.ToDBParams(), session.Lang)
	m.sendMsg(upd, confirmText)
	m.sessions.Delete(userID)
}

// showNextParamAsNewMessage показывает следующий необходимый параметр (required, soft-required, optional)
// в виде нового сообщения с кнопками. Возвращает true, если показал шаг; false — если все параметры заполнены.
// Если editMsgID > 0 — редактирует существующее сообщение вместо отправки нового.
func (m *MessageHandler) showNextParamAsNewMessage(ctx context.Context, upd tgbotapi.Update, session *UserSession, userID string, chatID int64, editMsgID int) bool {
	step := nextUnfilledParam(session.Exercise, &session.Params, session.SkippedParams)
	if step == nil {
		return false
	}

	exName := exTextByLang[session.Lang][session.Exercise]
	lang := session.Lang
	cat := session.Exercise.Category()

	var (
		text string
		kb   tgbotapi.InlineKeyboardMarkup
	)

	switch step.Param {
	case ParamWeight:
		session.State = StateAwaitWeight
		weights, _ := m.statRepo.UniqueWeightsByExercise(ctx, userID, session.Exercise.String(), 5)
		if step.Optional {
			text = fmt.Sprintf(messagesByLang[lang][chooseOptionalWeight], exName) + "\n" + messagesByLang[lang][orWriteText]
			kb = optionalWeightInlineKeyboard(weights, lang)
		} else {
			text = fmt.Sprintf(messagesByLang[lang][chooseWeight], exName) + "\n" + messagesByLang[lang][orWriteText]
			kb = weightInlineKeyboard(weights, lang)
		}

	case ParamCount:
		session.State = StateAwaitCount
		var detail string
		if session.Params.WeightKg != nil {
			detail = " " + formatWeight(*session.Params.WeightKg, lang)
		}
		text = fmt.Sprintf(messagesByLang[lang][chooseCount], exName+detail) + "\n" + messagesByLang[lang][orWriteText]
		kb = countInlineKeyboard(lang)

	case ParamDistance:
		session.State = StateAwaitDistance
		if step.Optional {
			text = fmt.Sprintf(messagesByLang[lang][chooseOptionalDistance], exName) + "\n" + messagesByLang[lang][orWriteText]
			kb = optionalDistanceInlineKeyboard(lang)
		} else {
			text = fmt.Sprintf(messagesByLang[lang][chooseDistance], exName) + "\n" + messagesByLang[lang][orWriteText]
			kb = distanceInlineKeyboard(lang)
		}

	case ParamDuration:
		session.State = StateAwaitDuration
		var detail string
		if session.Params.DistanceM != nil {
			detail = " " + formatDistance(*session.Params.DistanceM, lang)
		}
		if step.Optional {
			text = fmt.Sprintf(messagesByLang[lang][chooseOptionalDuration], exName) + "\n" + messagesByLang[lang][orWriteText]
			kb = optionalDurationInlineKeyboard(cat, lang)
		} else {
			text = fmt.Sprintf(messagesByLang[lang][chooseDuration], exName+detail) + "\n" + messagesByLang[lang][orWriteText]
			kb = durationInlineKeyboard(cat, lang)
		}
	}

	m.sessions.Set(userID, session)
	session.LastBotMessageID = m.sendOrEditKeyboard(upd, chatID, editMsgID, text, kb)
	m.sessions.Set(userID, session)
	return true
}

// sendOrEditKeyboard редактирует существующее сообщение (если editMsgID > 0) или отправляет новое.
// Возвращает ID сообщения с клавиатурой.
func (m *MessageHandler) sendOrEditKeyboard(upd tgbotapi.Update, chatID int64, editMsgID int, text string, kb tgbotapi.InlineKeyboardMarkup) int {
	if editMsgID > 0 {
		m.editMessageWithKeyboard(chatID, editMsgID, text, kb)
		return editMsgID
	}
	return m.sendMsgWithKeyboard(upd, text, kb)
}

// sendMsgWithKeyboard отправляет сообщение с inline-кнопками и возвращает ID отправленного сообщения
func (m *MessageHandler) sendMsgWithKeyboard(upd tgbotapi.Update, text string, keyboard tgbotapi.InlineKeyboardMarkup) int {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboard
	msg.ParseMode = m.cfg.ReplyFormat
	sent, err := m.tgBot.Send(msg)
	if err != nil {
		m.Errorf("failed to send message with keyboard: %v", err)
		return 0
	}
	return sent.MessageID
}

// answerCallback отвечает на callback query (убирает "часики" на кнопке)
func (m *MessageHandler) answerCallback(callbackID, text string) {
	callback := tgbotapi.NewCallback(callbackID, text)
	if _, err := m.tgBot.Request(callback); err != nil {
		m.Errorf("failed to answer callback: %v", err)
	}
}

// editMessageWithKeyboard редактирует сообщение бота, обновляя текст и кнопки
func (m *MessageHandler) editMessageWithKeyboard(chatID int64, messageID int, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, text, keyboard)
	edit.ParseMode = m.cfg.ReplyFormat
	if _, err := m.tgBot.Send(edit); err != nil {
		m.Errorf("failed to edit message with keyboard: %v", err)
	}
}

// editMessageRemoveKeyboard редактирует сообщение, убирая кнопки
func (m *MessageHandler) editMessageRemoveKeyboard(chatID int64, messageID int, text string) {
	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	edit.ParseMode = m.cfg.ReplyFormat
	if _, err := m.tgBot.Send(edit); err != nil {
		m.Errorf("failed to edit message: %v", err)
	}
}

func ptr[T any](t T) *T { return &t }
