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

		// Проверяем, что образались вообще к нам
		hasMention := m.hasBotMention(upd.Message.Text)
		if !hasMention && upd.FromChat().IsGroup() {
			continue // Скипаем, если к нам не обращались или не писали нам в личку
		}

		// Чистим текст от мусора
		msgText := m.clearRawMsg(upd.Message.Text)
		// Определяем язык
		lang, err := m.detectLang(msgText)
		if err != nil {
			m.sendMsg(upd, "Can't detect a language😶 Please, use the only one keyboard layout chars")
		}

		// Если ничего не осталось, отправляем соответствующий ответ
		if msgText == "" {
			m.sendMsg(upd, messagesByLang[lang][emptyMessage])
		}

		// Достаём пользователя
		userID := strconv.FormatInt(upd.Message.From.ID, 10)
		text, err := m.handle(ctx, msgText, userID, lang)
		if err != nil { // В случае ошибки сообщаем об этом
			text = messagesByLang[lang][errMsg]
			m.Error(ctx, "an error occurred", "message", msgText, "userID", userID, "err", err.Error()) // И логируем её
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
		return m.handleHelp(remainedText, lang)
	default:
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

	// Заменяем все отступы и переносы строк на одиночный пробел
	reSpaces := regexp.MustCompile(`\s+`)
	withoutSpaces := reSpaces.ReplaceAllString(withoutMention, " ")

	const dashPlaceHolder = "DASHPLACEHOLDER"
	// Делаем специальный плейсхолдер с тире, чтобы не удалить лишние тире
	reHyphen := regexp.MustCompile(`(\d)\s*-\s*(\d)`)
	withPlacehoder := reHyphen.ReplaceAllString(withoutSpaces, fmt.Sprintf("${1}%s${2}", dashPlaceHolder))

	// Убираем символы пунктуации
	rePunct := regexp.MustCompile(`[[:punct:]]`)
	withoutPuncts := rePunct.ReplaceAllString(withPlacehoder, "")

	// Теперь возвращаем тире обратно на место плейсхолдера
	withDashes := strings.ReplaceAll(withoutPuncts, dashPlaceHolder, "-")

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
		return fmt.Sprintf("%s. %s: %s", messagesByLang[lang][emptyEx], messagesByLang[lang][listEx], allExTextByLang(lang)), nil
	}

	words := strings.Split(rawMsg, " ")
	ex, ok := exerciseByLang[lang][words[0]]
	if !ok {
		return fmt.Sprintf("%s: %s. %s: %s", messagesByLang[lang][cantRecognizeEx], words[0], messagesByLang[lang][listEx], allExTextByLang(lang)), nil
	}

	// Если в упражнении должно быть задано количество
	if ex.mustHaveCnt() {
		if len(words) <= 1 {
			return messagesByLang[lang][cntRequired], nil
		}

		cnt, err := strconv.ParseFloat(words[1], 64)
		if err != nil {
			return fmt.Sprintf("%s: %s", messagesByLang[lang][cntInvalid], words[1]), nil //nolint:nilerr
		} else if cnt < 1 {
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
	for i = range words {
		if textContainsAllExerciseWords(words[i], lang) {
			i++ // Пропускаем это слово, фильтр будет пустой, значит вытащим и так всё
			break
		}

		ex, ok := exerciseByLang[lang][words[i]]
		if !ok { // Если не распознано, мы наверное дошли до интервала, остановимся
			break
		}

		exrs = append(exrs, ex)
		i++ // Сдвигаем на 1 для остатка слов
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
			periodsFilter, invPeriods = m.prepareCorrectAndInvalidPeriods(periodWords, lang)
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

	table, err := m.buildTableByStat(stats, lang)
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

func (m *MessageHandler) prepareCorrectAndInvalidPeriods(periods []string, lang language) (res periods, invalid []string) {
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
			invalid = append(invalid, periods[i])
			continue
		}

		// Иначе это должны быть даты, обработаем их
		p, inv := m.periodByTime(periods[i])
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

func (m *MessageHandler) periodByTime(interval string) (p period, invalid []string) {
	// И разибваем части по оставшемуся тире между датами
	intervals := strings.Split(interval, "-")

	// Если дата только одна, тогда from и to одинаковы
	if len(intervals) == 1 {
		t, err := m.parseDate(intervals[0])
		if err != nil {
			return period{}, []string{intervals[0]}
		}

		return period{from: t, to: t}, nil
	}

	// Если даты две
	if len(intervals) == 2 {
		from, err := m.parseDate(intervals[0])
		if err != nil {
			invalid = append(invalid, intervals[0])
		}

		to, err := m.parseDate(intervals[1])
		if err != nil {
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

func (m *MessageHandler) buildTableByStat(in []db.GroupedStatistic, lang language) (string, error) {
	if len(in) == 0 {
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

func (m *MessageHandler) handleHelp(rawMsg string, lang language) (string, error) {
	switch _, c := m.detectCmd(rawMsg, lang); c {
	case unknownCmd:
		return fmt.Sprintf(messagesByLang[lang][commonHelpMsg], m.cfg.Name), nil
	case addCmd:
		return fmt.Sprintf(messagesByLang[lang][addHelpMsg], m.cfg.Name) +
			fmt.Sprintf("%s: `%s`", messagesByLang[lang][listEx], allExTextByLang(lang)), nil
	case showCmd:
		return fmt.Sprintf(messagesByLang[lang][showHelpMsg], m.cfg.Name) +
			fmt.Sprintf("%s: %s", messagesByLang[lang][listPeriod], allPeriodsByLang(lang)), nil
	case helpCmd:
		return messagesByLang[lang][helpHelpMsg], nil
	default:
		return fmt.Sprintf("`%s` %s", rawMsg, messagesByLang[lang][cmdNotSupported]), nil
	}
}
