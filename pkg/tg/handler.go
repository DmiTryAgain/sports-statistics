package tg

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/DmiTryAgain/sports-statistics/pkg/db"

	"github.com/go-pg/pg/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vmkteam/embedlog"
)

var (
	errCantDetectLang = errors.New("can't detect language")
)

var (
	enLangRe = regexp.MustCompile(`^[a-zA-Z0-9\s.,!?'"@#$%^&*()\-_=+;:<>/\\|}{\[\]\p{So}]*$`)
	ruLangRe = regexp.MustCompile(`^[а-яА-ЯёЁ0-9\s.,!?'"@#$%^&*()\-_=+;:<>/\\|}{\[\]\p{So}]*$`)
)

type MessageHandler struct {
	embedlog.Logger

	dbc      *pg.DB
	statRepo db.StatisticRepo
	tgBot    *tgbotapi.BotAPI
	cfg      Bot
}

func New(logger embedlog.Logger, dbc *pg.DB, tgBot *tgbotapi.BotAPI, cfg Bot) *MessageHandler {
	h := &MessageHandler{
		Logger:   logger,
		dbc:      dbc,
		cfg:      cfg,
		tgBot:    tgBot,
		statRepo: db.NewStatisticRepo(dbc),
	}

	return h
}

func (m *MessageHandler) ListenAndHandle(ctx context.Context) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = int(m.cfg.Timeout.Duration.Seconds())

	// Get updates chan to listen to them
	updates := m.tgBot.GetUpdatesChan(updateConfig)

	// Listen messages
	for upd := range updates {
		if upd.Message == nil {
			continue
		}

		lowerText := strings.ToLower(upd.Message.Text)
		// Проверяем, что обращались вообще к нам
		hasMention := m.hasBotMention(lowerText)
		if !hasMention && upd.FromChat().IsGroup() {
			continue // Скипаем, если к нам не обращались или не писали нам в личку
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
			continue
		}

		// Если ничего не осталось, отправляем соответствующий ответ
		if msgText == "" {
			m.Print(ctx, "received empty message", "rawMsg", upd.Message.Text, "userID", userID)
			m.sendMsg(upd, messagesByLang[lang][emptyMessage])
			continue
		}

		text, err := m.handle(ctx, msgText, userID, lang)
		if err != nil { // В случае ошибки сообщаем об этом
			text = messagesByLang[lang][errMsg]
			m.Error(ctx, "an error occurred", "rawMsg", upd.Message.Text, "userID", userID, "err", err.Error()) // И логируем её
		}

		m.sendMsg(upd, text)
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
	withPlacehoder := reHyphen.ReplaceAllString(withoutMention, fmt.Sprintf("${1}%s${2}", dashPlaceHolder))

	// Убираем символы пунктуации
	rePunct := regexp.MustCompile(`[[:punct:]]`)
	withoutPuncts := rePunct.ReplaceAllString(withPlacehoder, "")

	// Заменяем все отступы и переносы строк на одиночный пробел
	reSpaces := regexp.MustCompile(`\s+`)
	withoutSpaces := reSpaces.ReplaceAllString(withoutPuncts, " ")

	// Теперь возвращаем тире обратно на место плейсхолдера
	withDashes := strings.ReplaceAll(withoutSpaces, dashPlaceHolder, "-")

	// Убираем пробелы по краям и возвращаем
	return strings.TrimSpace(withDashes)
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

func (m *MessageHandler) handleAdd(ctx context.Context, rawMsg, tgUserID string, lang language) (string, error) {
	if rawMsg == "" {
		m.Print(ctx, "received empty message", "msg", rawMsg, "userID", tgUserID)
		return fmt.Sprintf("%s. %s: %s", messagesByLang[lang][emptyEx], messagesByLang[lang][listEx], allExTextByLang(lang)), nil
	}

	var i int // Нужно для того, чтобы понять, где заканчиваются слова упражнений
	words := strings.Split(rawMsg, " ")
	ex, ok := exerciseByLang[lang][words[i]]
	if !ok {
		m.Print(ctx, "received unknown exercise", "msg", rawMsg, "userID", tgUserID, "exercise", words[0])
		return fmt.Sprintf("%s: %s. %s: %s", messagesByLang[lang][cantRecognizeEx], words[0], messagesByLang[lang][listEx], allExTextByLang(lang)), nil
	}

	// Проверим, что тут могло быть упражнение из нескольких слов
	multiwordExName := words[i]
	for len(words) > i+1 {
		multiwordExName = fmt.Sprintf("%s %s", multiwordExName, words[i+1])
		multiwordEx, exists := exerciseByLang[lang][multiwordExName]
		if !exists { // Если не распознано, мы наверное дошли до интервала, остановимся
			break
		}

		// Если оно реально состоит из 2х слов, снова сдвигаем i на следующее слово
		if ex == multiwordEx {
			i++
			continue
		}

		// Останавливаемся, если упражнения различаются. Мы захватили уже следующее
		break
	}
	i++ // Нам нужно начать со следующего слова после упражнения

	// Если в упражнении должно быть задано количество
	if ex.mustHaveCnt() {
		if len(words[i:]) < 1 { // Проверяем, что слова вообще остались после упражнений
			m.Print(ctx, "exercise must have count", "msg", rawMsg, "userID", tgUserID, "exercise", ex)
			return messagesByLang[lang][cntRequired], nil
		}

		cnt, err := strconv.ParseFloat(words[i], 64)
		if err != nil {
			return fmt.Sprintf("%s: %s", messagesByLang[lang][cntInvalid], words[i]), nil //nolint:nilerr
		} else if cnt < 1 {
			m.Print(ctx, "invalid exercise count", "msg", rawMsg, "userID", tgUserID, "count", cnt)
			return messagesByLang[lang][cntGE], nil
		}

		_, err = m.statRepo.AddStatistic(ctx, &db.Statistic{
			TgUserID: tgUserID,
			Exercise: ex.String(),
			Count:    cnt,
			Params:   nil,
			StatusID: 1,
		})
		if err != nil {
			return "", err
		}
	}

	return messagesByLang[lang][exAdded], nil
}

func (m *MessageHandler) handleShow(ctx context.Context, rawMsg, tgUserID string, lang language) (res string, err error) {
	if rawMsg == "" {
		m.Print(ctx, "received empty message", "msg", rawMsg, "userID", tgUserID)
		return fmt.Sprintf("%s. %s: %s", messagesByLang[lang][emptyEx], messagesByLang[lang][listEx], allExTextByLang(lang)), nil
	}

	// Удаляем ненужный предлог
	rawMsg = m.clearFromTo(rawMsg, lang)
	// Разбиваем по пробелам
	words := strings.Split(rawMsg, " ")

	var (
		exrs Exercises // Сюда запишем список упражнений, по которым надо будет фильтрануть
		i    int       // Здесь запомним, на каком элементе выйдем из цикла
	)

	// Идём по каждому слову и ищем упражнения, который надо достать до первого фейла
	for i < len(words) {
		if textContainsAllExerciseWords(words[i], lang) {
			m.Print(ctx, "the message contains all exercises", "msg", rawMsg, "userID", tgUserID, "all exercises word", words[i])
			i++ // Пропускаем это слово, фильтр будет пустой, значит вытащим и так всё
			break
		}

		ex, ok := exerciseByLang[lang][words[i]]
		if !ok { // Если не распознано, мы наверное дошли до интервала, остановимся
			break
		}
		// Если слова ещё не закончились, проверим, что тут могло быть упражнение из нескольких слов
		multiwordExName := words[i]
		for len(words) > i+1 {
			multiwordExName = fmt.Sprintf("%s %s", multiwordExName, words[i+1])
			multiwordEx, exists := exerciseByLang[lang][multiwordExName]
			if !exists { // Если не распознано, мы наверное дошли до интервала, остановимся
				break
			}

			// Если оно реально состоит из 2х слов, снова сдвигаем i на следующее слово
			if ex == multiwordEx {
				i++
				continue
			}

			// Останавливаемся, если упражнения различаются. Мы захватили уже следующее
			break
		}

		exrs = append(exrs, ex)
		i++ // Сдвигаем на следующее слово после текущего упражнения
	}

	// Проверяем, если вышли, и не нашли ни одного упражнения
	if len(exrs) == 0 && i == 0 {
		return fmt.Sprintf("%s. %s: %s", messagesByLang[lang][cantRecognizeEx], messagesByLang[lang][listEx], allExTextByLang(lang)), nil
	}

	var (
		periodsFilter periods
		invText       string
	)

	// Смотрим, есть ли кусок фразы про весь период в тексте.
	// Если нет, то парсим каждый период.
	// Если да, или если не задан, то считаем, что нужно взять за всё время.
	periodWords := words[i:]
	if len(periodWords) > 0 {
		// Слепим оставшуюся подстроку под период
		periodLeftPart := strings.Join(periodWords, " ")
		// Если в ней нет спец фразы для всех упражнений
		if !textContainsAllPeriodWords(periodLeftPart, lang) {
			var invPeriods []string
			// То идём парсить каждый элеент
			periodsFilter, invPeriods = m.prepareCorrectAndInvalidPeriods(ctx, periodWords, lang)
			invText = strings.Join(invPeriods, ", ")
		}
	}

	// Сразу добавляем в результат нераспознаные периоды
	if invText != "" {
		res += fmt.Sprintf("%s: %s\n", messagesByLang[lang][periodsInvalid], invText)
	}

	// Теперь идём за статистикой
	s := db.GroupedStatisticSearch{
		StatisticSearch: db.StatisticSearch{
			TgUserID:  &tgUserID,
			Exercises: exrs.StringSlice(),
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

func (m *MessageHandler) clearFromTo(in string, lang language) string {
	for _, s := range cleanByLang[lang] {
		if strings.Contains(in, s) {
			in = strings.ReplaceAll(in, strings.TrimRight(s, " "), "")
		}
	}

	return in
}

func (m *MessageHandler) prepareCorrectAndInvalidPeriods(ctx context.Context, periods []string, lang language) (res periods, invalid []string) {
	// Проходимся по каждому периоду
	for i := range periods {
		// Если он текстовый
		isText := m.langReByLang(lang).MatchString(periods[i])

		// То обработаем, попробуем взять интервалы из текста
		if isText {
			p, ok := m.periodByText(periods[i], time.Now(), lang)
			if ok { // Если получилось, добавляем в результат
				res = append(res, p)
				continue
			}

			// Иначе добавляем в невалидные
			m.Print(ctx, "captured invalid text period", "period", periods[i])
			invalid = append(invalid, periods[i])
			continue
		}

		// Иначе это должны быть даты, обработаем их
		p, inv := m.periodByTime(ctx, periods[i])
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

		return period{from: t, to: t}, nil
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

		return period{from: from, to: to}, invalid
	}

	return period{}, []string{interval}
}

func (m *MessageHandler) parseDate(date string) (time.Time, error) {
	// Пробуем распарсить в формате с полным годом
	parsed, err := time.Parse("02012006", date)
	if err != nil {
		// Если не получилось, пробуем распарсить в формате с коротким годом
		parsed, err = time.Parse("020106", date)
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

	tmpl := fmt.Sprintf(""+
		"%s\t%s\t%s\n"+
		"{{ range .Stat }}"+
		"{{ .TranslatedExercise }}\t{{ .SumCount }}\t{{ .Sets }}\n"+
		"{{ end }}", messagesByLang[lang][tableExCol], messagesByLang[lang][tableCntCol], messagesByLang[lang][tableSetCol],
	)

	type Data struct {
		Stat []GroupedStatistic
	}

	return tableFromTemplate("feedTrafficDeviations", tmpl, Data{Stat: NewGroupedStatisticList(in, lang)})
}

func tableFromTemplate(name, tmpl string, data interface{}) (string, error) {
	t := template.Must(template.New(name).Parse(tmpl))
	var b strings.Builder
	b.WriteString("```\n")

	wr := tabwriter.NewWriter(&b, 0, 1, 4, ' ', 0)
	err := t.Execute(wr, data)
	if err != nil {
		return "", err
	}

	err = wr.Flush()
	if err != nil {
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
