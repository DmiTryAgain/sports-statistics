package tg

import (
	"context"
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

type MessageHandler struct {
	embedlog.Logger

	dbc      *pg.DB
	statRepo db.StatisticRepo
	tgBot    *tgbotapi.BotAPI
	cfg      Bot

	commonHelpMsg string
	addHelpMsg    string
	showHelpMsg   string
	helpHelpMsg   string
	errMsg        string
}

func New(logger embedlog.Logger, dbc *pg.DB, tgBot *tgbotapi.BotAPI, cfg Bot) *MessageHandler {
	h := &MessageHandler{
		Logger:   logger,
		dbc:      dbc,
		cfg:      cfg,
		tgBot:    tgBot,
		statRepo: db.NewStatisticRepo(dbc),
	}

	h.initMessages()

	return h
}

func (m *MessageHandler) initMessages() {
	m.commonHelpMsg = fmt.Sprintf(
		"Привет! Я помогу вести статистику твоих спортивных упражнений."+
			"Ты же ведь занимаешься спортом, верно?🤔\n"+
			"Пиши мне в личные сообщения. В группах обращайся ко мне вот так: `@%s`"+
			"Список поддерживаемых команд: \n"+
			"На добавление: `Сделал` или `Добавь` \n"+
			"На показ статистики: `Покажи` \n"+
			"Чтобы посмотреть помощь по каждой комманде, отправь: `помощь` *название команды*\n"+
			"Например: `Помощь Добавь`",
		m.cfg.Name,
	)

	m.addHelpMsg = fmt.Sprintf(
		"Чтобы записать результаты, отправь команду на добавление упражнения (`сделал`). Затем, через "+
			"пробел укажи упражнение, которое сделал. Далее через пробел укажи сделанное количество \n"+
			"Например, ты сделал подход из 10 подтягиваний. Чтобы я всё корректно записал, напиши мне "+
			"`@%s сделал подтягивание 10`\n"+
			"Список доступных упражнений: `%s`",
		m.cfg.Name,
		exercises().String(),
	)

	m.showHelpMsg =
		"Чтобы показать статистику, отправь команду `Покажи`. Затем укажи название упражнения." +
			"*Можно ввести несколько, через запятую*, например, `подтягивание, отжимание`.\n" +
			"Далее укажи период, за который ты хочешь посмотреть статистику. Период будет корректно " +
			"распознан, если после указанных упражнений последует предлог *за*. Периодов можно указывать " +
			"несколько через запятую. Для каждого периода нужно так же нужен предлог *за*.\n" +
			"Например, нужно вывести статистику по подтягиваниям за сегодня, за 15.10.2022, " +
			"за период с 01.10.2022 по 10.10.2022. Чтобы периоды обработались корректно, введи периоды" +
			"следующим образом:\n" +
			"`за сегодня, за 15.10.2022, за 01.10.2022-10.10.2022`\n" +
			"Если период будет указан некорректно, результат будет без учёта некорректного периода. Если при " +
			"вводе интервала дата *от* окажется больше даты *до*, они поменяются местами и результат за этот " +
			"период будет найден корректно.\n" +
			"В итоге корректная команда будет выглядеть следующим образом: \n" +
			"`@%s покажи подтягивание, отжимание за сегодня, за 15.10.2022, за 01.10.2022-10.10.2022`\n" +
			"Список поддерживаемых текстовых периодов: " //TODO

	m.helpHelpMsg = "Помощь к команде помощи не предусмотрена. Надо ж было додуматься попросить помощь команде помощи🤔"

	m.errMsg = "❌ Произошла ошибка при просмотре статистики. Попробуйте позже"
}

func (m *MessageHandler) ListenAndHandle(ctx context.Context) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = int(m.cfg.Timeout.Duration.Seconds())

	// Get updates chan to listen to them
	updates := m.tgBot.GetUpdatesChan(updateConfig)

	// Listen messages
	for update := range updates {
		if update.Message == nil {
			continue
		}

		text, err := m.Handle(ctx, update)
		if err != nil {
			text = m.errMsg // TODO: handle better
			m.Error(ctx, err.Error())
		} else if text == "" {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := m.tgBot.Send(msg); err != nil {
			// TODO: make retries
			m.Errorf("failed to send message: %v", err)
		}
	}
}

func (m *MessageHandler) Handle(ctx context.Context, upd tgbotapi.Update) (string, error) {
	// Проверяем, что образались вообще к нам
	hasMention := m.hasBotMention(upd.Message.Text)
	if !hasMention && upd.FromChat().IsGroup() {
		return "", nil // Скипаем, если к нам не обращались или не писали нам в личку
	}

	msgText := m.clearRawMsg(upd.Message.Text)
	// Обрабатываем, если ничего не осталось
	if msgText == "" {
		return "Чё?", nil
	}

	userID := strconv.FormatInt(upd.Message.From.ID, 10)

	switch remainedText, c := m.evaluateCmd(msgText); c {
	case addCmd:
		return m.handleAdd(ctx, remainedText, userID)
	case showCmd:
		return m.handleShow(ctx, remainedText, userID)
	case helpCmd:
		return m.handleHelp(remainedText)
	default:
		return "Не могу обработать введёную Вами команду", nil
	}
}

// hasBotMention Проверяет, был ли бот заменшенен
func (m *MessageHandler) hasBotMention(msgTxt string) bool {
	return strings.Contains(msgTxt, "@"+strings.ToLower(m.cfg.Name))
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

// evaluateCmd Рассчитывает, какого типа команда, строку без названия команды и саму команду
func (m *MessageHandler) evaluateCmd(rawMsg string) (cleaned string, cmd cmd) {
	// Берём первое слово, чтобы понять, что за команда
	words := strings.SplitN(rawMsg, " ", 2)
	if len(words) == 0 {
		return rawMsg, unknownCmd
	}

	if len(words) > 1 {
		cleaned = words[1]
	}

	return cleaned, cmdByWord[strings.ToLower(words[0])]
}

func (m *MessageHandler) handleAdd(ctx context.Context, rawMsg string, tgUserID string) (string, error) {
	if rawMsg == "" {
		return "Упражнение не задано", nil
	}

	words := strings.Split(rawMsg, " ")
	ex, ok := exerciseByWord[words[0]]
	if !ok {
		return fmt.Sprintf("Неизвестное упражнение: %s", words[0]), nil
	}

	// Если в упражнении должн
	if ex.mustHaveCnt() {
		if len(words) <= 1 {
			return fmt.Sprintf("Для упражнения `%s` должно быть указано количество повторений", words[0]), nil
		}

		cnt, err := strconv.ParseFloat(words[1], 64)
		if err != nil {
			return fmt.Sprintf("Указано некорректное количество повторений: %s", words[1]), nil //nolint:nilerr
		} else if cnt < 1 {
			return "Количество повторений должно быть от 1 и более", nil
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

	return "Добавлено ✅", nil
}

func (m *MessageHandler) handleShow(ctx context.Context, rawMsg string, tgUserID string) (res string, err error) {
	if rawMsg == "" {
		return "Упражнение не задано", nil
	}

	// Удаляем ненужный предлог
	rawMsg = strings.ReplaceAll(rawMsg, "за", "")
	// Разбиваем по пробелам
	words := strings.Split(rawMsg, " ")

	var (
		exrs Exercises // Сюда запишем список упражнений, по которым надо будет фильтрануть
		i    int       // Здесь запомним, на каком элементе выйдем из цикла
	)

	// Идём по каждому слову и ищем упражнения, который надо достать до первого фейла
	for i = range words {
		if textContainsAllExerciseWords(words[i]) {
			i++ // Пропускаем это слово, фильтр будет пустой, значит вытащим и так всё
			break
		}

		ex, ok := exerciseByWord[words[i]]
		if !ok { // Если не распознано, мы наверное дошли до интервала, остановимся
			break
		}

		exrs = append(exrs, ex)
	}

	var (
		periodsFilter periods
		invText       string
	)

	// Смотрим, есть ли кусок фразы про весь период в тексте.
	// Если нет, то парсим каждый период.
	// Если да, или если не задан, то считаем, что нужно взять за всё время.
	if len(words[i:]) > 0 {
		// Слепим оставшуюся подстроку под период
		periodLeftPart := strings.Join(words[i:], " ")
		// Если в ней нет спец фразы для всех упражнений
		if !textContainsAllPeriodWords(periodLeftPart) {
			var invPeriods []string
			// То идём парсить каждый элеент
			periodsFilter, invPeriods = m.prepareCorrectAndInvalidPeriods(words[i:])
			invText = strings.Join(invPeriods, ", ")
		}
	}

	// Сразу добавляем в результат нераспознаные периоды
	if invText != "" {
		res += fmt.Sprintf("Нераспознаные периоды: %s\n", invText)
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
		return res + "Ничего не найдено 😢", nil
	}

	table, err := m.buildTableByStat(stats)
	if err != nil {
		return "", fmt.Errorf("build table by stat, err=%w", err)
	}

	res += table

	return res, nil
}

func (m *MessageHandler) prepareCorrectAndInvalidPeriods(periods []string) (res periods, invalid []string) {
	// Проходимся по каждому периоду
	for i := range periods {
		// Если он текстовый
		reWords := regexp.MustCompile(`^[а-яА-ЯёЁ]+$`)
		isText := reWords.MatchString(periods[i])

		// То обработаем, попробуем взять интервалы из текста
		if isText {
			p, ok := m.periodByText(periods[i], time.Now())
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

func (m *MessageHandler) periodByText(text string, now time.Time) (p period, ok bool) {
	switch periodByWord[text] {
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

func (m *MessageHandler) buildTableByStat(in []db.GroupedStatistic) (string, error) {
	if len(in) == 0 {
		return "", nil
	}

	const tmpl = "" +
		"упражнение\tкол-во\tподходы\n" +
		"{{ range .Stat }}" +
		"{{ .Exercise }}\t{{ .SumCount }}\t{{ .Sets }}\n" +
		"{{ end }}"

	type Data struct {
		Stat []db.GroupedStatistic
	}

	return tableFromTemplate("feedTrafficDeviations", tmpl, Data{Stat: in})
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

func (m *MessageHandler) handleHelp(rawMsg string) (string, error) {
	switch _, c := m.evaluateCmd(rawMsg); c {
	case unknownCmd:
		return m.commonHelpMsg, nil
	case addCmd:
		return m.addHelpMsg, nil
	case showCmd:
		return m.showHelpMsg, nil
	case helpCmd:
		return m.helpHelpMsg, nil
	default:
		return fmt.Sprintf("Команда `%s` не поддерживается", rawMsg), nil
	}
}
