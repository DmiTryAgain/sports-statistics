package tg

import "strings"

const (
	langRU language = "RU"
	langEN language = "EN"
)

const (
	unknownCmd cmd = iota
	addCmd
	showCmd
	helpCmd
)

const (
	pullUpEx       Exercise = "pullUp"
	muscleUpEx     Exercise = "muscleUp"
	pushUpEx       Exercise = "pushUp"
	dipsEx         Exercise = "dip"
	absEx          Exercise = "abs"
	squatEx        Exercise = "squat"
	lungeEx        Exercise = "lunge"
	burpeeEx       Exercise = "burpee"
	skippingRopeEx Exercise = "skippingRope"
	//joggingEx      Exercise = "jogging"
	allEx Exercise = "all"
)

const (
	todayPeriod              textPeriod = "today"
	yesterdayPeriod          textPeriod = "yesterday"
	dayBeforeYesterdayPeriod textPeriod = "dayBeforeYesterday"
	weekPeriod               textPeriod = "week"
	monthPeriod              textPeriod = "month"
	yearPeriod               textPeriod = "year"
	allPeriod                textPeriod = "all"
)

var (
	cmdByLang = map[language]map[string]cmd{
		langRU: {
			// addCmd
			"сделал":    addCmd,
			"сделол":    addCmd,
			"делал":     addCmd,
			"делол":     addCmd,
			"совершил":  addCmd,
			"совершыл":  addCmd,
			"савершил":  addCmd,
			"савершыл":  addCmd,
			"намутил":   addCmd,
			"номутил":   addCmd,
			"замутил":   addCmd,
			"зомутил":   addCmd,
			"нафигачил": addCmd,
			"нахерачил": addCmd,
			"нахуячил":  addCmd,
			"добавь":    addCmd,
			"дабавь":    addCmd,
			"добавьте":  addCmd,
			"дабавьте":  addCmd,
			"добавте":   addCmd,
			"дабавте":   addCmd,
			"добавить":  addCmd,
			"дабавить":  addCmd,
			"выполнил":  addCmd,
			"виполнил":  addCmd,
			"выпалнил":  addCmd,
			"випалнил":  addCmd,
			"въебал":    addCmd,
			"вебал":     addCmd,
			"разъебал":  addCmd,
			"разебал":   addCmd,
			"разыбал":   addCmd,

			// showCmd
			"покажи":     showCmd,
			"покаж":      showCmd,
			"покажь":     showCmd,
			"покажы":     showCmd,
			"пакажи":     showCmd,
			"пакаж":      showCmd,
			"пакажь":     showCmd,
			"пакажы":     showCmd,
			"покожы":     showCmd,
			"покожи":     showCmd,
			"покож":      showCmd,
			"покожь":     showCmd,
			"статистика": showCmd,
			"стата":      showCmd,
			"выведи":     showCmd,
			"вывиди":     showCmd,

			// helpCmd
			"помоги":   helpCmd,
			"памаги":   helpCmd,
			"помогите": helpCmd,
			"памагити": helpCmd,
			"помагити": helpCmd,
			"помагите": helpCmd,
			"помощь":   helpCmd,
			"помощ":    helpCmd,
			"хелп":     helpCmd,
		},
		langEN: {
			// addCmd
			"add":    addCmd,
			"ad":     addCmd,
			"insert": addCmd,
			"ins":    addCmd,
			"store":  addCmd,
			"save":   addCmd,
			"record": addCmd,
			"track":  addCmd,
			"enter":  addCmd,
			"done":   addCmd,

			// showCmd
			"show":    showCmd,
			"showme":  showCmd,
			"display": showCmd,
			"view":    showCmd,
			"list":    showCmd,
			"fetch":   showCmd,
			"see":     showCmd,
			"peek":    showCmd,
			"shoe":    showCmd,
			"shew":    showCmd,
			"watch":   showCmd,

			// helpCmd
			"help":    helpCmd,
			"hlp":     helpCmd,
			"hulp":    helpCmd,
			"guide":   helpCmd,
			"manual":  helpCmd,
			"support": helpCmd,
			"suport":  helpCmd,
			"info":    helpCmd,
			"inf":     helpCmd,
			"explain": helpCmd,
		},
	}

	exerciseByLang = map[language]map[string]Exercise{
		langRU: {
			//	pullUpEx
			"подтягивание": pullUpEx,
			"подтягивания": pullUpEx,
			"падтягивания": pullUpEx,
			"падтягивание": pullUpEx,
			"подтягиваний": pullUpEx,
			"падтягиваний": pullUpEx,
			"подтягиванье": pullUpEx,
			"подтягиванья": pullUpEx,
			"падтягиванья": pullUpEx,
			"падтягиванье": pullUpEx,
			"потягивание":  pullUpEx,
			"потягиваний":  pullUpEx,
			"патягивание":  pullUpEx,
			"патягиваний":  pullUpEx,
			"потягиванье":  pullUpEx,
			"патягиванье":  pullUpEx,
			"потягиваня":   pullUpEx,
			"патягиваня":   pullUpEx,
			"подтянулся":   pullUpEx,
			"падтянулся":   pullUpEx,

			// muscleUpEx
			"выход":        muscleUpEx,
			"выхад":        muscleUpEx,
			"выхот":        muscleUpEx,
			"выхат":        muscleUpEx,
			"выход силы":   muscleUpEx,
			"выхад силы":   muscleUpEx,
			"выхот силы":   muscleUpEx,
			"выхат силы":   muscleUpEx,
			"виход":        muscleUpEx,
			"вихад":        muscleUpEx,
			"вихот":        muscleUpEx,
			"вихат":        muscleUpEx,
			"виход силы":   muscleUpEx,
			"вихад силы":   muscleUpEx,
			"вихот силы":   muscleUpEx,
			"вихат силы":   muscleUpEx,
			"выходов":      muscleUpEx,
			"выхадов":      muscleUpEx,
			"выхотов":      muscleUpEx,
			"выхатов":      muscleUpEx,
			"выходов силы": muscleUpEx,
			"выхадов силы": muscleUpEx,
			"выхотов силы": muscleUpEx,
			"выхатов силы": muscleUpEx,
			"виходов":      muscleUpEx,
			"вихадов":      muscleUpEx,
			"вихотов":      muscleUpEx,
			"вихатов":      muscleUpEx,
			"виходов силы": muscleUpEx,
			"вихадов силы": muscleUpEx,
			"вихотов силы": muscleUpEx,
			"вихатов силы": muscleUpEx,

			// pushUpEx
			"отжимание": pushUpEx,
			"отжимания": pushUpEx,
			"атжимание": pushUpEx,
			"атжимания": pushUpEx,
			"ажимание":  pushUpEx,
			"ажимания":  pushUpEx,
			"оджимание": pushUpEx,
			"оджимания": pushUpEx,
			"отжиманье": pushUpEx,
			"отжиманья": pushUpEx,
			"атжиманье": pushUpEx,
			"атжиманья": pushUpEx,
			"ажиманье":  pushUpEx,
			"ажиманья":  pushUpEx,
			"оджиманье": pushUpEx,
			"оджиманья": pushUpEx,
			"анжуманя":  pushUpEx,
			"ажимане":   pushUpEx,
			"ажиманя":   pushUpEx,
			"отжиманий": pushUpEx,
			"оджиманий": pushUpEx,
			"анжуманий": pushUpEx,
			"ажиманей":  pushUpEx,
			"ажиманий":  pushUpEx,

			// dipsEx
			"брусья":  dipsEx,
			"бруся":   dipsEx,
			"брусьях": dipsEx,
			"брусьев": dipsEx,

			// absEx
			"пресс":    absEx,
			"прес":     absEx,
			"пресса":   absEx,
			"преса":    absEx,
			"пресуха":  absEx,
			"прессуха": absEx,
			"пресуху":  absEx,
			"прессуху": absEx,

			// squatEx
			"приседания": squatEx,
			"приседанья": squatEx,
			"приседаня":  squatEx,
			"приседание": squatEx,
			"приседанье": squatEx,
			"приседане":  squatEx,
			"приседаний": squatEx,
			"присидания": squatEx,
			"присиданья": squatEx,
			"присиданя":  squatEx,
			"присидание": squatEx,
			"присиданье": squatEx,
			"присидане":  squatEx,
			"присиданий": squatEx,

			// lungeEx
			"выпады":  lungeEx,
			"выпадов": lungeEx,
			"выпада":  lungeEx,
			"выпад":   lungeEx,

			// burpeeEx
			"бёрпи": burpeeEx,
			"берпи": burpeeEx,

			// skippingRopeEx
			"скакалка": skippingRopeEx,
			"скакалку": skippingRopeEx,
			"скокалка": skippingRopeEx,
			"скокалку": skippingRopeEx,
			"скакалки": skippingRopeEx,
			"скокалки": skippingRopeEx,

			// joggingEx
			//"бег":      joggingEx,
			//"бегал":    joggingEx,
			//"пробежал": joggingEx,
			//"пробежка": joggingEx,

			// all
			"всё":            allEx,
			"все":            allEx,
			"фсe":            allEx,
			"фсё":            allEx,
			"всё упражнения": allEx,
			"все упражнения": allEx,
			"фсe упражнения": allEx,
			"фсё упражнения": allEx,
			"всё упрожнения": allEx,
			"все упрожнения": allEx,
			"фсe упрожнения": allEx,
			"фсё упрожнения": allEx,
			"всё упражненя":  allEx,
			"все упражненя":  allEx,
			"фсe упражненя":  allEx,
			"фсё упражненя":  allEx,
			"всё упрожненя":  allEx,
			"все упрожненя":  allEx,
			"фсe упрожненя":  allEx,
			"фсё упрожненя":  allEx,
			"вся активность": allEx,
			"фся активность": allEx,
			"вся октивность": allEx,
			"фся октивность": allEx,
			"вся активнасть": allEx,
			"фся активнасть": allEx,
			"вся октивнасть": allEx,
			"фся октивнасть": allEx,
		},
		langEN: {
			//	pullUpEx
			"pull":      pullUpEx,
			"pulls":     pullUpEx,
			"pull-up":   pullUpEx,
			"pulls-up":  pullUpEx,
			"pull-ups":  pullUpEx,
			"pullup":    pullUpEx,
			"pullups":   pullUpEx,
			"pull up":   pullUpEx,
			"pull ups":  pullUpEx,
			"chin-up":   pullUpEx,
			"chin-ups":  pullUpEx,
			"chinup":    pullUpEx,
			"chinups":   pullUpEx,
			"chin up":   pullUpEx,
			"chin ups":  pullUpEx,
			"chinning":  pullUpEx,
			"chinnings": pullUpEx,
			"pulup":     pullUpEx,
			"pulups":    pullUpEx,
			"poolup":    pullUpEx,
			"poolups":   pullUpEx,
			"pullip":    pullUpEx,
			"pullips":   pullUpEx,

			// muscleUpEx
			"muscleup":   muscleUpEx,
			"muscleups":  muscleUpEx,
			"muscle-up":  muscleUpEx,
			"muscle-ups": muscleUpEx,
			"muscle up":  muscleUpEx,
			"muscle ups": muscleUpEx,
			"mucleup":    muscleUpEx,
			"mucleups":   muscleUpEx,
			"mucle-up":   muscleUpEx,
			"mucle-ups":  muscleUpEx,
			"mucle up":   muscleUpEx,
			"mucle ups":  muscleUpEx,
			"musclup":    muscleUpEx,
			"musclups":   muscleUpEx,
			"muscl-up":   muscleUpEx,
			"muscl-ups":  muscleUpEx,
			"muscl up":   muscleUpEx,
			"muscl ups":  muscleUpEx,

			// pushUpEx
			"pushup":    pushUpEx,
			"pushups":   pushUpEx,
			"push-up":   pushUpEx,
			"push-ups":  pushUpEx,
			"push up":   pushUpEx,
			"push ups":  pushUpEx,
			"pressup":   pushUpEx,
			"pressups":  pushUpEx,
			"press-up":  pushUpEx,
			"press-ups": pushUpEx,
			"press up":  pushUpEx,
			"press ups": pushUpEx,
			"puship":    pushUpEx,
			"puships":   pushUpEx,
			"push-ip":   pushUpEx,
			"push-ips":  pushUpEx,
			"push ip":   pushUpEx,
			"push ips":  pushUpEx,
			"pusshup":   pushUpEx,
			"pusshups":  pushUpEx,
			"pussh-up":  pushUpEx,
			"pussh-ups": pushUpEx,
			"pussh up":  pushUpEx,
			"pussh ups": pushUpEx,

			// dipsEx
			"dip":           dipsEx,
			"dips":          dipsEx,
			"parallel bars": dipsEx,
			"dipp":          dipsEx,
			"dipps":         dipsEx,
			"deep":          dipsEx,
			"deeps":         dipsEx,

			// absEx
			"abs":       absEx,
			"abdominal": absEx,
			"core":      absEx,
			"crunches":  absEx,
			"abbs":      absEx,
			"aps":       absEx,

			// squatEx
			"squat":   squatEx,
			"squats":  squatEx,
			"sqat":    squatEx,
			"sqats":   squatEx,
			"squaut":  squatEx,
			"squauts": squatEx,

			// lungeEx
			"lunge":  lungeEx,
			"lunges": lungeEx,
			"lunje":  lungeEx,
			"lunjes": lungeEx,
			"longe":  lungeEx,
			"longes": lungeEx,

			// burpeeEx
			"burpee":  burpeeEx,
			"burpees": burpeeEx,
			"burpe":   burpeeEx,
			"burpes":  burpeeEx,
			"burpy":   burpeeEx,
			"burpys":  burpeeEx,

			// skippingRopeEx
			"skippingrope":   skippingRopeEx,
			"skippingropes":  skippingRopeEx,
			"skipping-rope":  skippingRopeEx,
			"skipping-ropes": skippingRopeEx,
			"skipping rope":  skippingRopeEx,
			"skipping ropes": skippingRopeEx,
			"jumprope":       skippingRopeEx,
			"jumpropes":      skippingRopeEx,
			"jump-rope":      skippingRopeEx,
			"jump-ropes":     skippingRopeEx,
			"jump rope":      skippingRopeEx,
			"jump ropes":     skippingRopeEx,
			"skipingrope":    skippingRopeEx,
			"skipingropes":   skippingRopeEx,
			"skiping rope":   skippingRopeEx,
			"skiping ropes":  skippingRopeEx,
			"skiprope":       skippingRopeEx,
			"skipropes":      skippingRopeEx,
			"skip-rope":      skippingRopeEx,
			"skip-ropes":     skippingRopeEx,
			"skip rope":      skippingRopeEx,
			"skip ropes":     skippingRopeEx,

			// joggingEx
			//"jogging": joggingEx,
			//"joging":  joggingEx,
			//"joggin":  joggingEx,
			//"run":     joggingEx,
			//"running": joggingEx,
			//"trot":    joggingEx,
			//"sprint":  joggingEx,

			// all
			"all":        allEx,
			"everything": allEx,
			"total":      allEx,
			"full":       allEx,
			"al":         allEx,
			"aall":       allEx,
		},
	}

	periodByLang = map[language]map[string]textPeriod{
		langRU: {
			"сегодня": todayPeriod,
			"севодня": todayPeriod,
			"сиводня": todayPeriod,

			"вчера": yesterdayPeriod,
			"вчира": yesterdayPeriod,
			"фчира": yesterdayPeriod,
			"фчера": yesterdayPeriod,

			"позавчера": dayBeforeYesterdayPeriod,
			"позавчира": dayBeforeYesterdayPeriod,
			"позафчира": dayBeforeYesterdayPeriod,
			"позафчера": dayBeforeYesterdayPeriod,
			"пазафчера": dayBeforeYesterdayPeriod,
			"пазавчера": dayBeforeYesterdayPeriod,
			"пазафчира": dayBeforeYesterdayPeriod,

			"неделя": weekPeriod,
			"неделю": weekPeriod,
			"неделе": weekPeriod,
			"недели": weekPeriod,
			"ниделя": weekPeriod,
			"ниделю": weekPeriod,
			"ниделе": weekPeriod,
			"нидели": weekPeriod,

			"месяц":   monthPeriod,
			"месяца":  monthPeriod,
			"месяцев": monthPeriod,
			"месяцы":  monthPeriod,
			"месяци":  monthPeriod,
			"месец":   monthPeriod,
			"месеца":  monthPeriod,
			"месецев": monthPeriod,
			"месецы":  monthPeriod,
			"месеци":  monthPeriod,
			"месиц":   monthPeriod,
			"месица":  monthPeriod,
			"месицев": monthPeriod,
			"месицы":  monthPeriod,
			"месици":  monthPeriod,

			"год": yearPeriod,
			"гот": yearPeriod,

			"всё время":   allPeriod,
			"все время":   allPeriod,
			"всегда":      allPeriod,
			"всигда":      allPeriod,
			"всекда":      allPeriod,
			"всикда":      allPeriod,
			"весь период": allPeriod,
			"весь периуд": allPeriod,
			"весь периут": allPeriod,
			"весь пириод": allPeriod,
			"весь пириуд": allPeriod,
			"весь пириут": allPeriod,
		},
		langEN: {
			"today":       todayPeriod,
			"tdy":         todayPeriod,
			"tod":         todayPeriod,
			"2day":        todayPeriod,
			"this day":    todayPeriod,
			"current day": todayPeriod,
			"curr day":    todayPeriod,
			"currday":     todayPeriod,

			"yesterday":    yesterdayPeriod,
			"yday":         yesterdayPeriod,
			"ystrdy":       yesterdayPeriod,
			"last day":     yesterdayPeriod,
			"previous day": yesterdayPeriod,
			"prev day":     yesterdayPeriod,
			"prevday":      yesterdayPeriod,
			"yesturday":    yesterdayPeriod,
			"yeterday":     yesterdayPeriod,

			"daybeforeyesterday":     dayBeforeYesterdayPeriod,
			"day before yesterday":   dayBeforeYesterdayPeriod,
			"a day before yesterday": dayBeforeYesterdayPeriod,
			"2 days ago":             dayBeforeYesterdayPeriod,
			"2days ago":              dayBeforeYesterdayPeriod,
			"dayBeforeYest":          dayBeforeYesterdayPeriod,
			"dbYesterday":            dayBeforeYesterdayPeriod,

			"week":         weekPeriod,
			"wk":           weekPeriod,
			"7 days":       weekPeriod,
			"7days":        weekPeriod,
			"weekly":       weekPeriod,
			"this week":    weekPeriod,
			"current week": weekPeriod,
			"cur week":     weekPeriod,
			"wekk":         weekPeriod,
			"weak":         weekPeriod,

			"month":          monthPeriod,
			"mth":            monthPeriod,
			"30 days":        monthPeriod,
			"30days":         monthPeriod,
			"calendar month": monthPeriod,
			"this month":     monthPeriod,
			"current month":  monthPeriod,
			"cur month":      monthPeriod,
			"moneth":         monthPeriod,
			"mounth":         monthPeriod,

			"year":         yearPeriod,
			"yr":           yearPeriod,
			"12 months":    yearPeriod,
			"12months":     yearPeriod,
			"annual":       yearPeriod,
			"this year":    yearPeriod,
			"current year": yearPeriod,
			"cur year":     yearPeriod,
			"yaer":         yearPeriod,
			"yera":         yearPeriod,

			"all":        allPeriod,
			"everything": allPeriod,
			"total":      allPeriod,
			"full":       allPeriod,
			"al":         allPeriod,
			"aall":       allPeriod,
		},
	}

	cmdTextByLang = map[language]map[cmd]string{
		langRU: {
			addCmd:  "добавь",
			showCmd: "покажи",
			helpCmd: "помощь",
		},
		langEN: {
			addCmd:  "add",
			showCmd: "show",
			helpCmd: "help",
		},
	}

	exTextByLang = map[language]map[Exercise]string{
		langRU: {
			pullUpEx:       "подтягивания",
			muscleUpEx:     "выход силы",
			pushUpEx:       "отжимания",
			dipsEx:         "брусья",
			absEx:          "пресс",
			squatEx:        "приседания",
			lungeEx:        "выпады",
			burpeeEx:       "бёрпи",
			skippingRopeEx: "скакалка",
			//joggingEx:      "бег",
		},
		langEN: {
			pullUpEx:       "pull-ups",
			muscleUpEx:     "muscle-ups",
			pushUpEx:       "push-ups",
			dipsEx:         "dips",
			absEx:          "abs",
			squatEx:        "squats",
			lungeEx:        "lunges",
			burpeeEx:       "burpee",
			skippingRopeEx: "skipping rope",
			//joggingEx:      "jogging",
		},
	}
	periodTextByLang = map[language]map[textPeriod]string{
		langRU: {
			todayPeriod:              "сегодня",
			yesterdayPeriod:          "вчера",
			dayBeforeYesterdayPeriod: "позавчера",
			weekPeriod:               "неделя",
			monthPeriod:              "месяц",
			yearPeriod:               "год",
			allPeriod:                "всё время",
		},
		langEN: {
			todayPeriod:              "today",
			yesterdayPeriod:          "yesterday",
			dayBeforeYesterdayPeriod: "a day before yesterday",
			weekPeriod:               "week",
			monthPeriod:              "month",
			yearPeriod:               "year",
			allPeriod:                "all",
		},
	}
)

const (
	emptyMessage = iota
	listCmd
	listEx
	listPeriod
	cantRecognizeCmd
	cmdNotSupported
	emptyEx
	cantRecognizeEx
	cntRequired
	cntInvalid
	cntGE
	exAdded
	periodsInvalid
	nothingFound
	tableExCol
	tableCntCol
	tableSetCol
	commonHelpMsg
	addHelpMsg
	showHelpMsg
	helpHelpMsg
	errMsg
)

var (
	messagesByLang = map[language]map[int]string{
		langRU: {
			emptyMessage:     "Чё?",
			listCmd:          "Список поддерживаемых команд",
			listEx:           "Список поддерживаемых упражнений",
			listPeriod:       "Список поддерживаемых текстовых периодов",
			cantRecognizeCmd: "Команда не распознана",
			cmdNotSupported:  "Команда не поддерживается",
			emptyEx:          "Упражнение не задано",
			cantRecognizeEx:  "Упражнение не распознано",
			cntRequired:      "Для этого упражнения требуется ввести количество повторений",
			cntInvalid:       "Указано некорректное количество повторений",
			cntGE:            "Количество повторений должно быть от 1 и более",
			exAdded:          "Добавлено ✅",
			periodsInvalid:   "Нераспознаные периоды",
			nothingFound:     "Ничего не найдено 😢",
			tableExCol:       "упражнение",
			tableCntCol:      "кол-во",
			tableSetCol:      "подходы",
			commonHelpMsg: "Привет! Я помогу вести статистику твоих спортивных упражнений.\n" +
				"Ты же ведь занимаешься спортом, верно?🤔\n" +
				"Пиши мне в личные сообщения. В группах обращайся ко мне вот так: `@%s`\n" +
				"Список поддерживаемых команд:\n" +
				"На добавление: `Сделал` или `Добавь`\n" +
				"На показ статистики: `Покажи`\n" +
				"Чтобы посмотреть помощь по каждой комманде, отправь: `помощь` *название команды*\n" +
				"Например: `Помощь Добавь`",
			addHelpMsg: "Чтобы записать результаты, напиши ключевое слово (команду) на добавление упражнения `сделал`. Затем, через " +
				"пробел укажи выполненное упражнение. Далее укажи сделанное количество.\n" +
				"Например, ты подтянулся 10 раз. Чтобы я всё корректно записал, напиши мне:\n" +
				"`@%s сделал подтягивание 10`",
			showHelpMsg: "Чтобы показать статистику, напиши ключевое слово (команду) `Покажи`. Затем укажи упражнение.\n" +
				"*Можно ввести несколько, через запятую*, например, `подтягивание, отжимание`.\n" +
				"Далее укажи период, за который ты хочешь посмотреть статистику. Периодов можно указывать " +
				"несколько через используя пробелы.\n" +
				"Например, нужно вывести статистику по подтягиваниям за сегодня, за 15.10.2022, " +
				"за период с 01.10.2022 по 10.10.2022. Чтобы периоды обработались корректно, введи периоды" +
				"следующим образом:\n" +
				"`за сегодня, за 15.10.2022, за 01.10.2022-10.10.2022`\n" +
				"Некорректные периоды будут проигнорированы и выведены. Если при " +
				"вводе интервала дата *от* окажется больше даты *до*, они поменяются местами и результат за этот " +
				"период будет найден корректно.\n" +
				"В итоге полная корректная команда будет выглядеть следующим образом: \n" +
				"`@%s покажи подтягивание, отжимание за сегодня, за 15.10.2022, за 01.10.2022-10.10.2022`\n",
			helpHelpMsg: "Помощь к команде помощи не предусмотрена. Надо ж было додуматься попросить помощь команде помощи🤔",
			errMsg:      "❌ Произошла ошибка. Попробуйте позже",
		},
		langEN: {
			emptyMessage:     "What?",
			listCmd:          "Supported commands",
			listEx:           "Supported exercises",
			listPeriod:       "Text period list",
			cantRecognizeCmd: "Can't recognize the command",
			cmdNotSupported:  "Command is not supported",
			emptyEx:          "Exercise is not assigned",
			cantRecognizeEx:  "Can't recognize the exercise",
			cntRequired:      "This exercise requires you to enter the number of repetitions",
			cntInvalid:       "Incorrect number of repetitions",
			cntGE:            "The number of repetitions should be 1 or more",
			exAdded:          "Added ✅",
			periodsInvalid:   "Invalid periods",
			nothingFound:     "Nothing found 😢",
			tableExCol:       "exercise",
			tableCntCol:      "reps",
			tableSetCol:      "sets",
			commonHelpMsg: "Hi there! I can keep your training statistic.\n" +
				"You do sports, right?🤔\n" +
				"Write me direct messages. Mention me in groups like this: `@%s`\n" +
				"List supported commands: \n" +
				"To add: `Add` or `Store` \n" +
				"To show statistic: `Show` \n" +
				"To get help for each command send: `help` *cmd name*\n" +
				"For instance: `Help add`",
			addHelpMsg: "To write statistic, write a message with key word (command) `add`. Then using spaces" +
				"write an exercise you done. Then write done reps. \n" +
				"For instance, you did 10 pull-ups. To store it correctly, write me\n" +
				"`@%s add push-ups 10`",
			showHelpMsg: "Write `Show` to show statistic, then write an exercise.\n" +
				"*you can write some exercises using spaces*, e.g. `pull-ups, push-ups`.\n" +
				"Then write a period you want to show a statistic. You can write some periods using spaces.\n " +
				"For instance, you want to watch a statistic of pull-ups for today, for 15.10.2022 and " +
				"from 01.10.2022 to 10.10.2022. I can handle it correctly if you write them like this:\n" +
				"`for today, for 15.10.2022, 01.10.2022-10.10.2022`\n" +
				"A wrong period will be ignored and printed." +
				"If you write an interval where *from* date later than *to* date, they will swap places and the result" +
				"will be found correctly." +
				"The full correct example is below:\n" +
				"`@%s show push-ups, pull-ups for today, for 15.10.2022, 01.10.2022-10.10.2022`\n",
			helpHelpMsg: "Help for the help command is not provided. How did you guess to ask help to help command?🤔",
			errMsg:      "❌ An error occurred. Try again later",
		},
	}

	cleanByLang = map[language][]string{
		langRU: {
			" за ",
			" с ",
			" по ",
			" до ",
		},
		langEN: {
			" from ",
			" for ",
			" to ",
		},
	}
)

func allCmdTextByLang(lang language) string {
	b := strings.Builder{}
	var i int
	textByCmd := cmdTextByLang[lang]
	for _, v := range textByCmd {
		i++
		b.WriteString("`")
		b.WriteString(v)
		b.WriteString("`")
		if i != len(textByCmd) {
			b.WriteString(", ")
		}
	}

	return b.String()
}

func allExTextByLang(lang language) string {
	b := strings.Builder{}
	var i int
	textByCmd := exTextByLang[lang]
	for _, v := range textByCmd {
		i++
		b.WriteString("`")
		b.WriteString(v)
		b.WriteString("`")
		if i != len(textByCmd) {
			b.WriteString(", ")
		}
	}

	return b.String()
}

func allPeriodsByLang(lang language) string {
	b := strings.Builder{}
	var i int
	textByCmd := periodTextByLang[lang]
	for _, v := range textByCmd {
		i++
		b.WriteString("`")
		b.WriteString(v)
		b.WriteString("`")
		if i != len(textByCmd) {
			b.WriteString(", ")
		}
	}

	return b.String()
}

var (
	// TODO: implement later
	//hasDistance = map[Exercise]struct{}{
	//	joggingEx: {},
	//}
	//hasWeight = map[Exercise]struct{}{
	//	pullUpEx:   {},
	//	muscleUpEx: {},
	//	dipsEx:     {},
	//}
	exHasCnt = map[Exercise]struct{}{
		pullUpEx:       {},
		muscleUpEx:     {},
		pushUpEx:       {},
		dipsEx:         {},
		absEx:          {},
		squatEx:        {},
		lungeEx:        {},
		burpeeEx:       {},
		skippingRopeEx: {},
	}
)
