package tg

import (
	"strings"
	"time"
)

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
	pullUpEx         Exercise = "pullUp"
	muscleUpEx       Exercise = "muscleUp"
	pushUpEx         Exercise = "pushUp"
	dipsEx           Exercise = "dip"
	absEx            Exercise = "abs"
	squatEx          Exercise = "squat"
	lungeEx          Exercise = "lunge"
	burpeeEx         Exercise = "burpee"
	skippingRopeEx   Exercise = "skippingRope"
	hyperextensionEx Exercise = "hyperextension"
	legRaiseEx       Exercise = "legRaise"

	joggingEx Exercise = "jogging"
	walkingEx Exercise = "walking"

	plankEx      Exercise = "plank"
	wallSitEx    Exercise = "wallSit"
	hangEx       Exercise = "hang"
	hollowHoldEx Exercise = "hollowHold"
	supermanEx   Exercise = "superman"
	sidePlankEx  Exercise = "sidePlank"

	weightHoldEx Exercise = "weightHold"

	benchPressEx       Exercise = "benchPress"
	deadliftEx         Exercise = "deadlift"
	latPulldownEx      Exercise = "latPulldown"
	legPressEx         Exercise = "legPress"
	preacherCurlEx     Exercise = "preacherCurl"
	shoulderPressEx    Exercise = "shoulderPress"
	bentOverRowEx      Exercise = "bentOverRow"
	dumbbellCurlEx     Exercise = "dumbbellCurl"
	legExtensionEx     Exercise = "legExtension"
	legCurlEx          Exercise = "legCurl"
	seatedRowEx        Exercise = "seatedRow"
	chestFlyEx         Exercise = "chestFly"
	tricepPushdownEx   Exercise = "tricepPushdown"
	romanianDeadliftEx Exercise = "romanianDeadlift"
	hipThrustEx        Exercise = "hipThrust"
	lateralRaiseEx     Exercise = "lateralRaise"
	shrugEx            Exercise = "shrug"

	allEx Exercise = "all"
)

const (
	todayPeriod              textPeriod = "today"
	yesterdayPeriod          textPeriod = "yesterday"
	dayBeforeYesterdayPeriod textPeriod = "dayBeforeYesterday"
	weekPeriod               textPeriod = "week"
	lastWeekPeriod           textPeriod = "lastWeek"
	weekBeforeLastPeriod     textPeriod = "weekBeforeLast"
	monthPeriod              textPeriod = "month"
	lastMonthPeriod          textPeriod = "lastMonth"
	monthBeforeLastPeriod    textPeriod = "monthBeforeLast"
	yearPeriod               textPeriod = "year"
	lastYearPeriod           textPeriod = "lastYear"
	yearBeforeLastPeriod     textPeriod = "yearBeforeLast"
	allPeriod                textPeriod = "all"

	weekdayMondayPeriod    textPeriod = "weekdayMonday"
	weekdayTuesdayPeriod   textPeriod = "weekdayTuesday"
	weekdayWednesdayPeriod textPeriod = "weekdayWednesday"
	weekdayThursdayPeriod  textPeriod = "weekdayThursday"
	weekdayFridayPeriod    textPeriod = "weekdayFriday"
	weekdaySaturdayPeriod  textPeriod = "weekdaySaturday"
	weekdaySundayPeriod    textPeriod = "weekdaySunday"
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
			"патягивания":  pullUpEx,
			"потягивания":  pullUpEx,
			"потягевания":  pullUpEx,
			"патягевания":  pullUpEx,
			"патягеваня":   pullUpEx,
			"потягеваня":   pullUpEx,

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
			// barbellSquat синонимы → squat
			"присед со штангой":  squatEx,
			"приседы со штангой": squatEx,
			"присед штанга":      squatEx,

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
			"бег":      joggingEx,
			"бегал":    joggingEx,
			"пробежал": joggingEx,
			"пробежка": joggingEx,
			"пробежку": joggingEx,
			"бежал":    joggingEx,

			// walkingEx
			"ходьба":   walkingEx,
			"хотьба":   walkingEx,
			"ходьбу":   walkingEx,
			"хотьбу":   walkingEx,
			"прогулка": walkingEx,
			"прогулку": walkingEx,
			"ходил":    walkingEx,
			"хадил":    walkingEx,
			"гулял":    walkingEx,

			// benchPressEx
			"жим":      benchPressEx,
			"жим лёжа": benchPressEx,
			"жим лежа": benchPressEx,
			"жым":      benchPressEx,
			"жым лёжа": benchPressEx,
			"жым лежа": benchPressEx,

			// deadliftEx
			"становая":      deadliftEx,
			"становая тяга": deadliftEx,
			"становую":      deadliftEx,
			"становую тягу": deadliftEx,
			"станавая":      deadliftEx,
			"станавая тяга": deadliftEx,
			"станавую":      deadliftEx,
			"станавую тягу": deadliftEx,

			// plankEx
			"планка":  plankEx,
			"планку":  plankEx,
			"планки":  plankEx,
			"планке":  plankEx,
			"планкой": plankEx,

			// wallSitEx
			"стульчик":   wallSitEx,
			"стульчек":   wallSitEx,
			"стулчик":    wallSitEx,
			"стулчек":    wallSitEx,
			"стульчиком": wallSitEx,

			// hangEx
			"вис":                hangEx,
			"вис на турнике":     hangEx,
			"вис на перекладине": hangEx,
			"висел":              hangEx,

			// hollowHoldEx
			"лодочка": hollowHoldEx,
			"лодочку": hollowHoldEx,
			"лодачку": hollowHoldEx,
			"лодачка": hollowHoldEx,

			// supermanEx
			"супермен":  supermanEx,
			"суперман":  supermanEx,
			"супермена": supermanEx,

			// sidePlankEx
			"боковая планка": sidePlankEx,
			"баковая планка": sidePlankEx,
			"боковая планку": sidePlankEx,
			"боковую планку": sidePlankEx,
			"боковаю планку": sidePlankEx,

			// weightHoldEx
			"удержание веса": weightHoldEx,
			"удержание":      weightHoldEx,
			"удиржание веса": weightHoldEx,
			"удиржание":      weightHoldEx,

			// hyperextensionEx
			"гиперэкстензия": hyperextensionEx,
			"гиперэкстензии": hyperextensionEx,
			"гиперэкстензию": hyperextensionEx,
			"гиперэкстензий": hyperextensionEx,
			"гиперэкстэнзия": hyperextensionEx,
			"гиперэкстэнзии": hyperextensionEx,
			"гиперэкстэнзию": hyperextensionEx,
			"гиперикстензия": hyperextensionEx,
			"гиперикстензии": hyperextensionEx,
			"экстензия":      hyperextensionEx,
			"экстензии":      hyperextensionEx,

			// legRaiseEx
			"подъём ног":          legRaiseEx,
			"подъем ног":          legRaiseEx,
			"подъёмы ног":         legRaiseEx,
			"подъемы ног":         legRaiseEx,
			"подъём ног в висе":   legRaiseEx,
			"подъем ног в висе":   legRaiseEx,
			"подъёмы ног в висе":  legRaiseEx,
			"подъемы ног в висе":  legRaiseEx,
			"поднятие ног":        legRaiseEx,
			"поднятие ног в висе": legRaiseEx,

			// latPulldownEx
			"тягу верхнего блока": latPulldownEx,
			"тяга верхнего блока": latPulldownEx,
			"тяга верхнева блока": latPulldownEx,
			"тяга верхнего блоко": latPulldownEx,
			"тягу верхнева блока": latPulldownEx,
			"тягу верхнего блоко": latPulldownEx,
			"верхний блок":        latPulldownEx,
			"верхняя тяга":        latPulldownEx,
			"верхнюю тягу":        latPulldownEx,

			// legPressEx
			"жим ногами": legPressEx,
			"жим нагами": legPressEx,
			"жым ногами": legPressEx,
			"жым нагами": legPressEx,

			// preacherCurlEx
			"скамья скотта":   preacherCurlEx,
			"скамья скота":    preacherCurlEx,
			"скамейка скотта": preacherCurlEx,
			"скамейка скота":  preacherCurlEx,
			"скамью скотта":   preacherCurlEx,
			"скамью скота":    preacherCurlEx,
			"скамейку скотта": preacherCurlEx,
			"скамейку скота":  preacherCurlEx,
			"скотта":          preacherCurlEx,
			"скота":           preacherCurlEx,

			// shoulderPressEx
			"жим стоя":        shoulderPressEx,
			"жым стоя":        shoulderPressEx,
			"армейский жим":   shoulderPressEx,
			"армейский жым":   shoulderPressEx,
			"армейскый жим":   shoulderPressEx,
			"армейскый жым":   shoulderPressEx,
			"жим над головой": shoulderPressEx,
			"жим над галавой": shoulderPressEx,

			// bentOverRowEx
			"тяга в наклоне":        bentOverRowEx,
			"тяга штанги в наклоне": bentOverRowEx,
			"тяга штанги в наклони": bentOverRowEx,
			"тяга в наклони":        bentOverRowEx,
			"тягу в наклоне":        bentOverRowEx,
			"тягу штанги в наклоне": bentOverRowEx,
			"тягу штанги в наклони": bentOverRowEx,
			"тягу в наклони":        bentOverRowEx,

			// dumbbellCurlEx
			"подъём гантелей":           dumbbellCurlEx,
			"подъем гантелей":           dumbbellCurlEx,
			"подъём гантелей на бицепс": dumbbellCurlEx,
			"подъем гантелей на бицепс": dumbbellCurlEx,
			"бицепс гантели":            dumbbellCurlEx,
			"бицепс гантелями":          dumbbellCurlEx,
			"сгибание на бицепс":        dumbbellCurlEx,
			"сгибания на бицепс":        dumbbellCurlEx,

			// legExtensionEx
			"разгибание ног":             legExtensionEx,
			"разгибания ног":             legExtensionEx,
			"разгибание ног в тренажёре": legExtensionEx,
			"разгибание ног в тренажере": legExtensionEx,

			// legCurlEx
			"сгибание ног":             legCurlEx,
			"сгибания ног":             legCurlEx,
			"сгибание ног в тренажёре": legCurlEx,
			"сгибание ног в тренажере": legCurlEx,

			// seatedRowEx
			"тяга нижнего блока":  seatedRowEx,
			"тягу нижнего блока":  seatedRowEx,
			"тяга нижниго блока":  seatedRowEx,
			"тягу нижниго блока":  seatedRowEx,
			"тяга нижнева блока":  seatedRowEx,
			"тягу нижнева блока":  seatedRowEx,
			"тяга нижнива блока":  seatedRowEx,
			"тягу нижнива блока":  seatedRowEx,
			"нижний блок":         seatedRowEx,
			"нижняя тяга":         seatedRowEx,
			"нижнюю тягу":         seatedRowEx,
			"горизонтальная тяга": seatedRowEx,
			"горизонтальную тягу": seatedRowEx,
			"горезонтальную тягу": seatedRowEx,
			"гаризонтальную тягу": seatedRowEx,
			"гарезонтальную тягу": seatedRowEx,
			"горизантальная тяга": seatedRowEx,
			"горизантальную тягу": seatedRowEx,
			"горезантальную тягу": seatedRowEx,
			"гаризантальную тягу": seatedRowEx,
			"гарезантальную тягу": seatedRowEx,
			"горизонталная тяга":  seatedRowEx,
			"горизонталную тягу":  seatedRowEx,
			"горезонталную тягу":  seatedRowEx,
			"гаризонталную тягу":  seatedRowEx,
			"гарезонталную тягу":  seatedRowEx,
			"горизанталная тяга":  seatedRowEx,
			"горизанталную тягу":  seatedRowEx,
			"горезанталную тягу":  seatedRowEx,
			"гаризанталную тягу":  seatedRowEx,
			"гарезанталную тягу":  seatedRowEx,

			// chestFlyEx
			"сведение рук":   chestFlyEx,
			"сведения рук":   chestFlyEx,
			"бабочка":        chestFlyEx,
			"бабачка":        chestFlyEx,
			"бабочку":        chestFlyEx,
			"бабачку":        chestFlyEx,
			"разведение рук": chestFlyEx,
			"разведения рук": chestFlyEx,
			"розведение рук": chestFlyEx,
			"розведения рук": chestFlyEx,

			// tricepPushdownEx
			"разгибание на трицепс": tricepPushdownEx,
			"разгибания на трицепс": tricepPushdownEx,
			"трицепс на блоке":      tricepPushdownEx,
			"трицепс блок":          tricepPushdownEx,

			// romanianDeadliftEx
			"румынская тяга": romanianDeadliftEx,
			"румынская":      romanianDeadliftEx,
			"румынскую тягу": romanianDeadliftEx,
			"румынскую":      romanianDeadliftEx,
			"румынка":        romanianDeadliftEx,
			"румынку":        romanianDeadliftEx,
			"мёртвая тяга":   romanianDeadliftEx,
			"мертвая тяга":   romanianDeadliftEx,

			// hipThrustEx
			"ягодичный мост":   hipThrustEx,
			"ягодичный мостик": hipThrustEx,
			"ягадичный мост":   hipThrustEx,
			"ягадичный мостик": hipThrustEx,

			// lateralRaiseEx
			"махи гантелями":            lateralRaiseEx,
			"махи гонтелями":            lateralRaiseEx,
			"махи гантелей":             lateralRaiseEx,
			"махи гонтелей":             lateralRaiseEx,
			"махи гантелий":             lateralRaiseEx,
			"махи гонтелий":             lateralRaiseEx,
			"махи":                      lateralRaiseEx,
			"разводка":                  lateralRaiseEx,
			"разводка гантелей":         lateralRaiseEx,
			"разведение гантелей":       lateralRaiseEx,
			"разведения гантелей":       lateralRaiseEx,
			"розводка гантелей":         lateralRaiseEx,
			"розведение гантелей":       lateralRaiseEx,
			"розведения гантелей":       lateralRaiseEx,
			"розвидение гантелей":       lateralRaiseEx,
			"розвидения гантелей":       lateralRaiseEx,
			"розвидене гантелей":        lateralRaiseEx,
			"розвиденя гантелей":        lateralRaiseEx,
			"разводка гонтелей":         lateralRaiseEx,
			"разведение гонтелей":       lateralRaiseEx,
			"разведения гонтелей":       lateralRaiseEx,
			"розводка гонтелей":         lateralRaiseEx,
			"розведение гонтелей":       lateralRaiseEx,
			"розведения гонтелей":       lateralRaiseEx,
			"розвидение гонтелей":       lateralRaiseEx,
			"розвидения гонтелей":       lateralRaiseEx,
			"розвидене гонтелей":        lateralRaiseEx,
			"розвиденя гонтелей":        lateralRaiseEx,
			"разводка гонтелий":         lateralRaiseEx,
			"разведение гонтелий":       lateralRaiseEx,
			"разведения гонтелий":       lateralRaiseEx,
			"розводка гонтелий":         lateralRaiseEx,
			"розведение гонтелий":       lateralRaiseEx,
			"розведения гонтелий":       lateralRaiseEx,
			"розвидение гонтелий":       lateralRaiseEx,
			"розвидения гонтелий":       lateralRaiseEx,
			"розвидене гонтелий":        lateralRaiseEx,
			"розвиденя гонтелий":        lateralRaiseEx,
			"махи в стороны":            lateralRaiseEx,
			"махи в стораны":            lateralRaiseEx,
			"подъём гантелей в стороны": lateralRaiseEx,
			"подъем гантелей в стороны": lateralRaiseEx,
			"подъём гонтелей в стороны": lateralRaiseEx,
			"подъем гонтелей в стороны": lateralRaiseEx,
			"подъём гантелий в стороны": lateralRaiseEx,
			"подъем гантелий в стороны": lateralRaiseEx,
			"подъём гонтелий в стороны": lateralRaiseEx,
			"подъем гонтелий в стороны": lateralRaiseEx,
			"подъём гантелей в стораны": lateralRaiseEx,
			"подъем гантелей в стораны": lateralRaiseEx,
			"подъём гонтелей в стораны": lateralRaiseEx,
			"подъем гонтелей в стораны": lateralRaiseEx,
			"подъём гантелий в стораны": lateralRaiseEx,
			"подъем гантелий в стораны": lateralRaiseEx,
			"подъём гонтелий в стораны": lateralRaiseEx,
			"подъем гонтелий в стораны": lateralRaiseEx,

			// shrugEx
			"шраги":             shrugEx,
			"шраг":              shrugEx,
			"шрагов":            shrugEx,
			"шраги со штангой":  shrugEx,
			"шраги с гантелями": shrugEx,

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
			"pullup":    pullUpEx,
			"pullups":   pullUpEx,
			"pull up":   pullUpEx,
			"pull ups":  pullUpEx,
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
			// barbellSquat синонимы → squat
			"barbell squat":  squatEx,
			"barbell squats": squatEx,
			"barbellsquat":   squatEx,

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
			"jogging": joggingEx,
			"joging":  joggingEx,
			"joggin":  joggingEx,
			"run":     joggingEx,
			"running": joggingEx,
			"trot":    joggingEx,
			"sprint":  joggingEx,

			// walkingEx
			"walking": walkingEx,
			"walk":    walkingEx,
			"hike":    walkingEx,
			"hiking":  walkingEx,

			// benchPressEx
			"bench":       benchPressEx,
			"bench press": benchPressEx,
			"benchpress":  benchPressEx,

			// deadliftEx
			"deadlift":  deadliftEx,
			"deadlifts": deadliftEx,
			"dead lift": deadliftEx,

			// plankEx
			"plank":  plankEx,
			"planks": plankEx,

			// wallSitEx
			"wall sit":  wallSitEx,
			"wall sits": wallSitEx,
			"wallsit":   wallSitEx,

			// hangEx
			"hang":      hangEx,
			"dead hang": hangEx,
			"deadhang":  hangEx,
			"bar hang":  hangEx,

			// hollowHoldEx
			"hollow hold":  hollowHoldEx,
			"hollow holds": hollowHoldEx,
			"hollowhold":   hollowHoldEx,
			"hollow":       hollowHoldEx,

			// supermanEx
			"superman":      supermanEx,
			"superman hold": supermanEx,
			"supermans":     supermanEx,

			// sidePlankEx
			"side plank":  sidePlankEx,
			"side planks": sidePlankEx,
			"sideplank":   sidePlankEx,

			// weightHoldEx
			"weight hold":  weightHoldEx,
			"weight holds": weightHoldEx,
			"weighthold":   weightHoldEx,

			// hyperextensionEx
			"hyperextension":  hyperextensionEx,
			"hyperextensions": hyperextensionEx,
			"hyper extension": hyperextensionEx,
			"hyper":           hyperextensionEx,
			"back extension":  hyperextensionEx,
			"back extensions": hyperextensionEx,

			// legRaiseEx
			"leg raise":          legRaiseEx,
			"leg raises":         legRaiseEx,
			"legraise":           legRaiseEx,
			"legraises":          legRaiseEx,
			"hanging leg raise":  legRaiseEx,
			"hanging leg raises": legRaiseEx,

			// latPulldownEx
			"lat pulldown":  latPulldownEx,
			"lat pulldowns": latPulldownEx,
			"latpulldown":   latPulldownEx,
			"pulldown":      latPulldownEx,
			"pulldowns":     latPulldownEx,
			"pull down":     latPulldownEx,
			"pull downs":    latPulldownEx,

			// legPressEx
			"leg press": legPressEx,
			"legpress":  legPressEx,

			// preacherCurlEx
			"preacher curl":  preacherCurlEx,
			"preacher curls": preacherCurlEx,
			"preachercurl":   preacherCurlEx,
			"scott bench":    preacherCurlEx,
			"scott curl":     preacherCurlEx,
			"scott curls":    preacherCurlEx,

			// shoulderPressEx
			"shoulder press": shoulderPressEx,
			"shoulderpress":  shoulderPressEx,
			"overhead press": shoulderPressEx,
			"overheadpress":  shoulderPressEx,
			"ohp":            shoulderPressEx,
			"military press": shoulderPressEx,

			// bentOverRowEx
			"bent over row":  bentOverRowEx,
			"bent over rows": bentOverRowEx,
			"bentoverrow":    bentOverRowEx,
			"barbell row":    bentOverRowEx,
			"barbell rows":   bentOverRowEx,

			// dumbbellCurlEx
			"dumbbell curl":  dumbbellCurlEx,
			"dumbbell curls": dumbbellCurlEx,
			"dumbbellcurl":   dumbbellCurlEx,
			"bicep curl":     dumbbellCurlEx,
			"bicep curls":    dumbbellCurlEx,
			"biceps curl":    dumbbellCurlEx,
			"biceps curls":   dumbbellCurlEx,

			// legExtensionEx
			"leg extension":  legExtensionEx,
			"leg extensions": legExtensionEx,
			"legextension":   legExtensionEx,

			// legCurlEx
			"leg curl":  legCurlEx,
			"leg curls": legCurlEx,
			"legcurl":   legCurlEx,

			// seatedRowEx
			"seated row":  seatedRowEx,
			"seated rows": seatedRowEx,
			"seatedrow":   seatedRowEx,
			"cable row":   seatedRowEx,
			"cable rows":  seatedRowEx,

			// chestFlyEx
			"chest fly":   chestFlyEx,
			"chest flys":  chestFlyEx,
			"chest flies": chestFlyEx,
			"pec deck":    chestFlyEx,
			"pec fly":     chestFlyEx,
			"pec flys":    chestFlyEx,

			// tricepPushdownEx
			"tricep pushdown":   tricepPushdownEx,
			"tricep pushdowns":  tricepPushdownEx,
			"triceppushdown":    tricepPushdownEx,
			"tricep extension":  tricepPushdownEx,
			"tricep extensions": tricepPushdownEx,

			// romanianDeadliftEx
			"romanian deadlift":  romanianDeadliftEx,
			"romanian deadlifts": romanianDeadliftEx,
			"rdl":                romanianDeadliftEx,
			"stiff leg deadlift": romanianDeadliftEx,

			// hipThrustEx
			"hip thrust":   hipThrustEx,
			"hip thrusts":  hipThrustEx,
			"hipthrust":    hipThrustEx,
			"glute bridge": hipThrustEx,

			// lateralRaiseEx
			"lateral raise":  lateralRaiseEx,
			"lateral raises": lateralRaiseEx,
			"lateralraise":   lateralRaiseEx,
			"side raise":     lateralRaiseEx,
			"side raises":    lateralRaiseEx,

			// shrugEx
			"shrug":  shrugEx,
			"shrugs": shrugEx,

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

			"прошлую неделю": lastWeekPeriod,
			"прошлой неделе": lastWeekPeriod,
			"прошлой недели": lastWeekPeriod,
			"прошлая неделя": lastWeekPeriod,
			"неделю назад":   lastWeekPeriod,
			"неделя назад":   lastWeekPeriod,
			"прошлую ниделю": lastWeekPeriod,
			"прошлой нидели": lastWeekPeriod,
			"ниделю назад":   lastWeekPeriod,
			"нидели назад":   lastWeekPeriod,

			"позапрошлую неделю": weekBeforeLastPeriod,
			"позапрошлой неделе": weekBeforeLastPeriod,
			"позапрошлой недели": weekBeforeLastPeriod,
			"позапрошлая неделя": weekBeforeLastPeriod,
			"позапрошлую ниделю": weekBeforeLastPeriod,
			"позапрошлой нидели": weekBeforeLastPeriod,
			"пазапрошлую неделю": weekBeforeLastPeriod,
			"пазапрошлой недели": weekBeforeLastPeriod,

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

			"прошлый месяц":   lastMonthPeriod,
			"прошлого месяца": lastMonthPeriod,
			"месяц назад":     lastMonthPeriod,
			"прошлый месец":   lastMonthPeriod,
			"прошлый месиц":   lastMonthPeriod,
			"месец назад":     lastMonthPeriod,
			"месиц назад":     lastMonthPeriod,

			"позапрошлый месяц":   monthBeforeLastPeriod,
			"позапрошлого месяца": monthBeforeLastPeriod,
			"пазапрошлый месяц":   monthBeforeLastPeriod,
			"позапрошлый месец":   monthBeforeLastPeriod,
			"пазапрошлый месец":   monthBeforeLastPeriod,

			"год": yearPeriod,
			"гот": yearPeriod,

			"прошлый год":   lastYearPeriod,
			"прошлого года": lastYearPeriod,
			"год назад":     lastYearPeriod,
			"прошлый гот":   lastYearPeriod,
			"гот назад":     lastYearPeriod,

			"позапрошлый год":   yearBeforeLastPeriod,
			"позапрошлого года": yearBeforeLastPeriod,
			"пазапрошлый год":   yearBeforeLastPeriod,
			"позапрошлый гот":   yearBeforeLastPeriod,
			"пазапрошлый гот":   yearBeforeLastPeriod,

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

			"понедельник":  weekdayMondayPeriod,
			"понедельника": weekdayMondayPeriod,
			"панедельник":  weekdayMondayPeriod,
			"понидельник":  weekdayMondayPeriod,
			"пн":           weekdayMondayPeriod,

			"вторник":  weekdayTuesdayPeriod,
			"вторника": weekdayTuesdayPeriod,
			"фторник":  weekdayTuesdayPeriod,
			"вт":       weekdayTuesdayPeriod,

			"среда": weekdayWednesdayPeriod,
			"среду": weekdayWednesdayPeriod,
			"среды": weekdayWednesdayPeriod,
			"сриду": weekdayWednesdayPeriod,
			"срида": weekdayWednesdayPeriod,
			"ср":    weekdayWednesdayPeriod,

			"четверг":  weekdayThursdayPeriod,
			"четверга": weekdayThursdayPeriod,
			"читверг":  weekdayThursdayPeriod,
			"четьверг": weekdayThursdayPeriod,
			"чт":       weekdayThursdayPeriod,

			"пятница": weekdayFridayPeriod,
			"пятницу": weekdayFridayPeriod,
			"пятницы": weekdayFridayPeriod,
			"питница": weekdayFridayPeriod,
			"пятнецу": weekdayFridayPeriod,
			"пт":      weekdayFridayPeriod,

			"суббота": weekdaySaturdayPeriod,
			"субботу": weekdaySaturdayPeriod,
			"субботы": weekdaySaturdayPeriod,
			"субота":  weekdaySaturdayPeriod,
			"суботу":  weekdaySaturdayPeriod,
			"суботы":  weekdaySaturdayPeriod,
			"сб":      weekdaySaturdayPeriod,

			"воскресенье": weekdaySundayPeriod,
			"воскресенья": weekdaySundayPeriod,
			"воскресения": weekdaySundayPeriod,
			"воскресение": weekdaySundayPeriod,
			"васкресенье": weekdaySundayPeriod,
			"воскрисенье": weekdaySundayPeriod,
			"вс":          weekdaySundayPeriod,
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

			"last week":     lastWeekPeriod,
			"lastweek":      lastWeekPeriod,
			"prev week":     lastWeekPeriod,
			"previous week": lastWeekPeriod,
			"week ago":      lastWeekPeriod,
			"a week ago":    lastWeekPeriod,

			"week before last": weekBeforeLastPeriod,
			"2 weeks ago":      weekBeforeLastPeriod,
			"2weeks ago":       weekBeforeLastPeriod,

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

			"last month":     lastMonthPeriod,
			"lastmonth":      lastMonthPeriod,
			"prev month":     lastMonthPeriod,
			"previous month": lastMonthPeriod,
			"month ago":      lastMonthPeriod,
			"a month ago":    lastMonthPeriod,

			"month before last": monthBeforeLastPeriod,
			"2 months ago":      monthBeforeLastPeriod,
			"2months ago":       monthBeforeLastPeriod,

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

			"last year":     lastYearPeriod,
			"lastyear":      lastYearPeriod,
			"prev year":     lastYearPeriod,
			"previous year": lastYearPeriod,
			"year ago":      lastYearPeriod,
			"a year ago":    lastYearPeriod,

			"year before last": yearBeforeLastPeriod,
			"2 years ago":      yearBeforeLastPeriod,
			"2years ago":       yearBeforeLastPeriod,

			"all":        allPeriod,
			"everything": allPeriod,
			"total":      allPeriod,
			"full":       allPeriod,
			"al":         allPeriod,
			"aall":       allPeriod,

			"monday":    weekdayMondayPeriod,
			"mon":       weekdayMondayPeriod,
			"tuesday":   weekdayTuesdayPeriod,
			"tue":       weekdayTuesdayPeriod,
			"tues":      weekdayTuesdayPeriod,
			"wednesday": weekdayWednesdayPeriod,
			"wed":       weekdayWednesdayPeriod,
			"thursday":  weekdayThursdayPeriod,
			"thu":       weekdayThursdayPeriod,
			"thur":      weekdayThursdayPeriod,
			"thurs":     weekdayThursdayPeriod,
			"friday":    weekdayFridayPeriod,
			"fri":       weekdayFridayPeriod,
			"saturday":  weekdaySaturdayPeriod,
			"sat":       weekdaySaturdayPeriod,
			"sunday":    weekdaySundayPeriod,
			"sun":       weekdaySundayPeriod,
		},
	}

	weekdayByPeriod = map[textPeriod]time.Weekday{
		weekdayMondayPeriod:    time.Monday,
		weekdayTuesdayPeriod:   time.Tuesday,
		weekdayWednesdayPeriod: time.Wednesday,
		weekdayThursdayPeriod:  time.Thursday,
		weekdayFridayPeriod:    time.Friday,
		weekdaySaturdayPeriod:  time.Saturday,
		weekdaySundayPeriod:    time.Sunday,
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
			pullUpEx:           "подтягивания",
			muscleUpEx:         "выход силы",
			pushUpEx:           "отжимания",
			dipsEx:             "брусья",
			absEx:              "пресс",
			squatEx:            "приседания",
			lungeEx:            "выпады",
			burpeeEx:           "бёрпи",
			skippingRopeEx:     "скакалка",
			hyperextensionEx:   "гиперэкстензия",
			legRaiseEx:         "подъём ног",
			joggingEx:          "бег",
			walkingEx:          "ходьба",
			plankEx:            "планка",
			wallSitEx:          "стульчик",
			hangEx:             "вис",
			hollowHoldEx:       "лодочка",
			supermanEx:         "супермен",
			sidePlankEx:        "боковая планка",
			weightHoldEx:       "удержание веса",
			benchPressEx:       "жим лёжа",
			deadliftEx:         "становая тяга",
			latPulldownEx:      "тяга верхнего блока",
			legPressEx:         "жим ногами",
			preacherCurlEx:     "скамья Скотта",
			shoulderPressEx:    "жим стоя",
			bentOverRowEx:      "тяга в наклоне",
			dumbbellCurlEx:     "подъём на бицепс",
			legExtensionEx:     "разгибание ног",
			legCurlEx:          "сгибание ног",
			seatedRowEx:        "тяга нижнего блока",
			chestFlyEx:         "сведение рук",
			tricepPushdownEx:   "разгибание на трицепс",
			romanianDeadliftEx: "румынская тяга",
			hipThrustEx:        "ягодичный мост",
			lateralRaiseEx:     "махи гантелями",
			shrugEx:            "шраги",
		},
		langEN: {
			pullUpEx:           "pull-ups",
			muscleUpEx:         "muscle-ups",
			pushUpEx:           "push-ups",
			dipsEx:             "dips",
			absEx:              "abs",
			squatEx:            "squats",
			lungeEx:            "lunges",
			burpeeEx:           "burpee",
			skippingRopeEx:     "skipping rope",
			hyperextensionEx:   "hyperextension",
			legRaiseEx:         "leg raise",
			joggingEx:          "jogging",
			walkingEx:          "walking",
			plankEx:            "plank",
			wallSitEx:          "wall sit",
			hangEx:             "hang",
			hollowHoldEx:       "hollow hold",
			supermanEx:         "superman",
			sidePlankEx:        "side plank",
			weightHoldEx:       "weight hold",
			benchPressEx:       "bench press",
			deadliftEx:         "deadlift",
			latPulldownEx:      "lat pulldown",
			legPressEx:         "leg press",
			preacherCurlEx:     "preacher curl",
			shoulderPressEx:    "shoulder press",
			bentOverRowEx:      "bent-over row",
			dumbbellCurlEx:     "dumbbell curl",
			legExtensionEx:     "leg extension",
			legCurlEx:          "leg curl",
			seatedRowEx:        "seated row",
			chestFlyEx:         "chest fly",
			tricepPushdownEx:   "tricep pushdown",
			romanianDeadliftEx: "romanian deadlift",
			hipThrustEx:        "hip thrust",
			lateralRaiseEx:     "lateral raise",
			shrugEx:            "shrugs",
		},
	}
	periodTextByLang = map[language]map[textPeriod]string{
		langRU: {
			todayPeriod:              "сегодня",
			yesterdayPeriod:          "вчера",
			dayBeforeYesterdayPeriod: "позавчера",
			weekPeriod:               "неделя",
			lastWeekPeriod:           "прошлая неделя",
			weekBeforeLastPeriod:     "позапрошлая неделя",
			monthPeriod:              "месяц",
			lastMonthPeriod:          "прошлый месяц",
			monthBeforeLastPeriod:    "позапрошлый месяц",
			yearPeriod:               "год",
			lastYearPeriod:           "прошлый год",
			yearBeforeLastPeriod:     "позапрошлый год",
			allPeriod:                "всё время",
			weekdayMondayPeriod:      "понедельник",
			weekdayTuesdayPeriod:     "вторник",
			weekdayWednesdayPeriod:   "среда",
			weekdayThursdayPeriod:    "четверг",
			weekdayFridayPeriod:      "пятница",
			weekdaySaturdayPeriod:    "суббота",
			weekdaySundayPeriod:      "воскресенье",
		},
		langEN: {
			todayPeriod:              "today",
			yesterdayPeriod:          "yesterday",
			dayBeforeYesterdayPeriod: "a day before yesterday",
			weekPeriod:               "week",
			lastWeekPeriod:           "last week",
			weekBeforeLastPeriod:     "week before last",
			monthPeriod:              "month",
			lastMonthPeriod:          "last month",
			monthBeforeLastPeriod:    "month before last",
			yearPeriod:               "year",
			lastYearPeriod:           "last year",
			yearBeforeLastPeriod:     "year before last",
			allPeriod:                "all",
			weekdayMondayPeriod:      "monday",
			weekdayTuesdayPeriod:     "tuesday",
			weekdayWednesdayPeriod:   "wednesday",
			weekdayThursdayPeriod:    "thursday",
			weekdayFridayPeriod:      "friday",
			weekdaySaturdayPeriod:    "saturday",
			weekdaySundayPeriod:      "sunday",
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
	tableWeightCol
	tableDistCol
	tableTimeCol
	weightRequired
	durationRequired
	distanceRequired
	distOrTimeRequired
	weightAndDurationRequired
	paramInvalid
	commonHelpMsg
	addHelpMsg
	showHelpMsg
	helpHelpMsg
	errMsg
	sessionExpired
	msgTooLong

	// ReplyKeyboard кнопки
	addBtn
	showBtn
	helpBtn

	// Сообщения для пошагового диалога
	welcomeMsg
	chooseExercise
	chooseExerciseOrText
	yourFrequent
	chooseWeight
	chooseCount
	chooseDistance
	chooseDuration
	enterCustomWeight
	enterCustomCount
	enterCustomDistance
	enterCustomDuration
	customInputBtn
	cancelBtn
	moreBtn
	backBtn
	allExBtn
	choosePeriod
	addedConfirmation
	quickCopyHint
	orWriteText
	skipBtn
	chooseOptionalWeight
	chooseOptionalDistance
	chooseOptionalDuration
)

var (
	messagesByLang = map[language]map[int]string{
		langRU: { //nolint:dupl
			emptyMessage:              "Чё?",
			listCmd:                   "Список поддерживаемых команд",
			listEx:                    "Список поддерживаемых упражнений",
			listPeriod:                "Список поддерживаемых текстовых периодов",
			cantRecognizeCmd:          "Команда не распознана",
			cmdNotSupported:           "Команда не поддерживается",
			emptyEx:                   "Упражнение не задано",
			cantRecognizeEx:           "Упражнение не распознано",
			cntRequired:               "Для этого упражнения требуется ввести количество повторений",
			cntInvalid:                "Указано некорректное количество повторений",
			cntGE:                     "Количество повторений должно быть от 1 и более",
			exAdded:                   "Добавлено ✅",
			periodsInvalid:            "Нераспознаные периоды",
			nothingFound:              "Ничего не найдено 😢",
			tableExCol:                "упражнение",
			tableCntCol:               "кол-во",
			tableSetCol:               "подходы",
			tableWeightCol:            "вес",
			tableDistCol:              "дистанция",
			tableTimeCol:              "время",
			weightRequired:            "Для этого упражнения нужно указать вес. Пример: жим 80кг 10",
			durationRequired:          "Для этого упражнения нужно указать время. Пример: планка 90сек",
			distanceRequired:          "Для этого упражнения нужно указать дистанцию. Пример: бег 5км 25мин",
			distOrTimeRequired:        "Нужно указать хотя бы дистанцию или время. Пример: бег 5км или бег 25мин",
			weightAndDurationRequired: "Нужно указать вес и время. Пример: удержание 40кг 30сек",
			paramInvalid:              "Не удалось распознать параметр: %s",
			commonHelpMsg: "Привет! Я помогу вести статистику твоих спортивных упражнений.\n" +
				"Ты же ведь занимаешься спортом, верно?🤔\n\n" +
				"Пиши мне в личные сообщения. В группах обращайся ко мне вот так: `@%s`\n\n" +
				"Список поддерживаемых команд:\n" +
				"• Добавить результат: `Сделал` или `Добавь`\n" +
				"• Показать статистику: `Покажи`\n" +
				"• Справка: `Помощь`\n\n" +
				"Чтобы посмотреть помощь по конкретной команде, отправь: `помощь` *название команды*\n" +
				"Например: `Помощь Добавь`",
			addHelpMsg: "Чтобы записать результаты, напиши команду `сделал`, затем упражнение и параметры.\n\n" +
				"*Обычные упражнения* — укажи количество повторений:\n" +
				"`@%[1]s сделал подтягивания 10`\n" +
				"`@%[1]s сделал отжимания 20`\n\n" +
				"*Упражнения с весом* — укажи вес с суффиксом (кг, г) и количество:\n" +
				"`@%[1]s сделал жим 80кг 10`\n" +
				"`@%[1]s сделал становую тягу 100кг 5`\n" +
				"`@%[1]s сделал тягу верхнего блока 60кг 12`\n\n" +
				"*Бег* — укажи дистанцию (км, м) и время (ч, мин, сек):\n" +
				"`@%[1]s сделал бег 5км 25мин`\n\n" +
				"*Планка* — укажи время:\n" +
				"`@%[1]s сделал планку 1мин 30сек`\n",
			showHelpMsg: "Чтобы показать статистику, напиши команду `Покажи`, затем упражнение и период.\n\n" +
				"*Можно указать несколько упражнений через пробел:*\n" +
				"`@%[1]s покажи подтягивания отжимания за неделю`\n\n" +
				"*Или всё сразу:*\n" +
				"`@%[1]s покажи всё за сегодня`\n\n" +
				"*Поддерживаемые периоды:* сегодня, вчера, позавчера, неделя, прошлая неделя, позапрошлая неделя, месяц, прошлый месяц, позапрошлый месяц, год, прошлый год, позапрошлый год, всё время.\n" +
				"Также можно указать день недели: понедельник / пн, вторник / вт, среда / ср, четверг / чт, пятница / пт, суббота / сб, воскресенье / вс.\n" +
				"Или точную дату или интервал:\n" +
				"`за 15.10.2025`\n" +
				"`за 01.10.2025-10.10.2025`\n\n" +
				"Полный пример:\n" +
				"`@%[1]s покажи подтягивания отжимания за сегодня за 01.10.2025-10.10.2025`\n",
			helpHelpMsg:    "Помощь к команде помощи не предусмотрена. Надо ж было додуматься попросить помощь к команде помощи🤔",
			errMsg:         "❌ Произошла ошибка. Попробуйте позже",
			sessionExpired: "Сессия истекла, начните заново",
			msgTooLong:     "Сообщение слишком длинное",

			addBtn:  "📝 Добавить",
			showBtn: "📊 Статистика",
			helpBtn: "❓ Помощь",

			welcomeMsg: "Привет! Я помогу вести статистику твоих спортивных упражнений.\n" +
				"Используй кнопки внизу для быстрого доступа к командам.",
			chooseExercise:         "Выбери упражнение:",
			chooseExerciseOrText:   "Выбери упражнение или напиши текстом.",
			yourFrequent:           "Твои частые:",
			chooseWeight:           "%s — укажи вес:",
			chooseCount:            "%s — сколько повторений?",
			chooseDistance:         "%s — укажи дистанцию:",
			chooseDuration:         "%s — укажи время:",
			enterCustomWeight:      "Введи вес (например: 85кг или 85)",
			enterCustomCount:       "Введи количество повторений",
			enterCustomDistance:    "Введи дистанцию (например: 5км или 5000м)",
			enterCustomDuration:    "Введи время (например: 25мин или 90сек)",
			customInputBtn:         "Другой",
			cancelBtn:              "Отмена",
			moreBtn:                "Ещё >>",
			backBtn:                "<< Назад",
			allExBtn:               "Всё",
			choosePeriod:           "%s — за какой период?",
			addedConfirmation:      "Добавлено ✅ %s: %s",
			quickCopyHint:          "Скопируй для быстрой вставки: %s",
			orWriteText:            "Или напиши текстом",
			skipBtn:                "Пропустить",
			chooseOptionalWeight:   "%s — добавить вес? (необязательно)",
			chooseOptionalDistance: "%s — указать дистанцию? (необязательно)",
			chooseOptionalDuration: "%s — указать время? (необязательно)",
		},
		langEN: { //nolint:dupl
			emptyMessage:              "What?",
			listCmd:                   "Supported commands",
			listEx:                    "Supported exercises",
			listPeriod:                "Text period list",
			cantRecognizeCmd:          "Can't recognize the command",
			cmdNotSupported:           "Command is not supported",
			emptyEx:                   "Exercise is not assigned",
			cantRecognizeEx:           "Can't recognize the exercise",
			cntRequired:               "This exercise requires you to enter the number of repetitions",
			cntInvalid:                "Incorrect number of repetitions",
			cntGE:                     "The number of repetitions should be 1 or more",
			exAdded:                   "Added ✅",
			periodsInvalid:            "Invalid periods",
			nothingFound:              "Nothing found 😢",
			tableExCol:                "exercise",
			tableCntCol:               "reps",
			tableSetCol:               "sets",
			tableWeightCol:            "weight",
			tableDistCol:              "distance",
			tableTimeCol:              "time",
			weightRequired:            "Weight is required for this exercise. Example: bench 80kg 10",
			durationRequired:          "Duration is required for this exercise. Example: plank 90sec",
			distanceRequired:          "Distance is required for this exercise. Example: run 5km 25min",
			distOrTimeRequired:        "Distance or duration is required. Example: run 5km or run 25min",
			weightAndDurationRequired: "Weight and duration are required. Example: weight hold 40kg 30sec",
			paramInvalid:              "Can't recognize parameter: %s",
			commonHelpMsg: "Hi there! I can keep your training statistics.\n" +
				"You do sports, right?🤔\n\n" +
				"Write me direct messages. In groups, mention me like this: `@%s`\n\n" +
				"Supported commands:\n" +
				"• Add a result: `Add` or `Done`\n" +
				"• Show statistics: `Show`\n" +
				"• Help: `Help`\n\n" +
				"To get help for a specific command, send: `help` *command name*\n" +
				"For example: `Help add`",
			addHelpMsg: "To record a result, write the command `add`, then the exercise and its parameters.\n\n" +
				"*Basic exercises* — specify the number of reps:\n" +
				"`@%[1]s add pull-ups 10`\n" +
				"`@%[1]s add push-ups 20`\n\n" +
				"*Weighted exercises* — specify weight with a suffix (kg, lbs) and reps:\n" +
				"`@%[1]s add bench press 80kg 10`\n" +
				"`@%[1]s add deadlift 100kg 5`\n" +
				"`@%[1]s add lat pulldown 60kg 12`\n\n" +
				"*Running* — specify distance (km, m) and time (h, min, sec):\n" +
				"`@%[1]s add jogging 5km 25min`\n\n" +
				"*Plank* — specify duration:\n" +
				"`@%[1]s add plank 1min 30sec`\n",
			showHelpMsg: "To show statistics, write the command `Show`, then an exercise and a period.\n\n" +
				"*You can specify multiple exercises separated by spaces:*\n" +
				"`@%[1]s show pull-ups push-ups for week`\n\n" +
				"*Or everything at once:*\n" +
				"`@%[1]s show all for today`\n\n" +
				"*Supported periods:* today, yesterday, week, last week, week before last, month, last month, month before last, year, last year, year before last, all.\n" +
				"You can also specify a weekday: monday / mon, tuesday / tue, wednesday / wed, thursday / thu, friday / fri, saturday / sat, sunday / sun.\n" +
				"Or an exact date or a range:\n" +
				"`for 15.10.2025`\n" +
				"`for 01.10.2025-10.10.2025`\n\n" +
				"Full example:\n" +
				"`@%[1]s show pull-ups push-ups for today for 01.10.2025-10.10.2025`\n",
			helpHelpMsg:    "Help for the help command is not provided. How did you even think to ask for help on the help command?🤔",
			errMsg:         "❌ An error occurred. Try again later",
			sessionExpired: "Session expired, please start again",
			msgTooLong:     "Message is too long",

			addBtn:  "📝 Add",
			showBtn: "📊 Statistics",
			helpBtn: "❓ Help",

			welcomeMsg: "Hi there! I can keep your training statistics.\n" +
				"Use the buttons below for quick access to commands.",
			chooseExercise:         "Choose an exercise:",
			chooseExerciseOrText:   "Choose an exercise or type it.",
			yourFrequent:           "Your frequent:",
			chooseWeight:           "%s — enter weight:",
			chooseCount:            "%s — how many reps?",
			chooseDistance:         "%s — enter distance:",
			chooseDuration:         "%s — enter duration:",
			enterCustomWeight:      "Enter weight (e.g. 85kg or 85)",
			enterCustomCount:       "Enter the number of reps",
			enterCustomDistance:    "Enter distance (e.g. 5km or 5000m)",
			enterCustomDuration:    "Enter duration (e.g. 25min or 90sec)",
			customInputBtn:         "Other",
			cancelBtn:              "Cancel",
			moreBtn:                "More >>",
			backBtn:                "<< Back",
			allExBtn:               "All",
			choosePeriod:           "%s — for what period?",
			addedConfirmation:      "Added ✅ %s: %s",
			quickCopyHint:          "Copy for quick paste: %s",
			orWriteText:            "Or type it",
			skipBtn:                "Skip",
			chooseOptionalWeight:   "%s — add weight? (optional)",
			chooseOptionalDistance: "%s — add distance? (optional)",
			chooseOptionalDuration: "%s — add duration? (optional)",
		},
	}

	prepositionByLang = map[language]map[string]struct{}{
		langRU: {
			"за": struct{}{},
			"с":  struct{}{},
			"по": struct{}{},
			"до": struct{}{},
		},
		langEN: {
			"from": struct{}{},
			"for":  struct{}{},
			"to":   struct{}{},
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

var replyButtonCmd = map[language]map[string]cmd{
	langRU: {
		"📝 Добавить":   addCmd,
		"📊 Статистика": showCmd,
		"❓ Помощь":     helpCmd,
	},
	langEN: {
		"📝 Add":        addCmd,
		"📊 Statistics": showCmd,
		"❓ Help":       helpCmd,
	},
}

var exerciseOrder = []Exercise{
	joggingEx, walkingEx,
	pullUpEx, pushUpEx, dipsEx, absEx, squatEx,
	benchPressEx, deadliftEx, legPressEx,

	plankEx, lungeEx, muscleUpEx, burpeeEx,
	skippingRopeEx, hyperextensionEx, legRaiseEx,
	chestFlyEx, hangEx, shoulderPressEx,

	bentOverRowEx, latPulldownEx, seatedRowEx,
	dumbbellCurlEx, preacherCurlEx, tricepPushdownEx,
	legExtensionEx, legCurlEx,
	romanianDeadliftEx, hipThrustEx,

	lateralRaiseEx, shrugEx, wallSitEx,
	hollowHoldEx, supermanEx, sidePlankEx, weightHoldEx,
}

const exercisesPerPage = 9 // количество кнопок упражнений на одной странице (3 ряда по 3)

var exerciseCategoryMap = map[Exercise]ExerciseCategory{
	// CategoryRepsWeight
	benchPressEx:       CategoryRepsWeight,
	deadliftEx:         CategoryRepsWeight,
	latPulldownEx:      CategoryRepsWeight,
	legPressEx:         CategoryRepsWeight,
	preacherCurlEx:     CategoryRepsWeight,
	shoulderPressEx:    CategoryRepsWeight,
	bentOverRowEx:      CategoryRepsWeight,
	dumbbellCurlEx:     CategoryRepsWeight,
	legExtensionEx:     CategoryRepsWeight,
	legCurlEx:          CategoryRepsWeight,
	seatedRowEx:        CategoryRepsWeight,
	chestFlyEx:         CategoryRepsWeight,
	tricepPushdownEx:   CategoryRepsWeight,
	romanianDeadliftEx: CategoryRepsWeight,
	hipThrustEx:        CategoryRepsWeight,
	lateralRaiseEx:     CategoryRepsWeight,
	shrugEx:            CategoryRepsWeight,

	// CategoryDistTime
	joggingEx: CategoryDistTime,
	walkingEx: CategoryDistTime,

	// CategoryDuration
	plankEx:      CategoryDuration,
	wallSitEx:    CategoryDuration,
	hangEx:       CategoryDuration,
	hollowHoldEx: CategoryDuration,
	supermanEx:   CategoryDuration,
	sidePlankEx:  CategoryDuration,

	// CategoryDurationWeight
	weightHoldEx: CategoryDurationWeight,
}

// exerciseOptionalParamsMap — необязательные параметры для конкретных упражнений
var exerciseOptionalParamsMap = map[Exercise][]ParamType{
	// CategoryReps с опциональным весом
	pullUpEx:         {ParamWeight},
	pushUpEx:         {ParamWeight},
	dipsEx:           {ParamWeight},
	squatEx:          {ParamWeight},
	lungeEx:          {ParamWeight},
	hyperextensionEx: {ParamWeight},

	// CategoryDuration с опциональным весом
	plankEx:     {ParamWeight},
	hangEx:      {ParamWeight},
	wallSitEx:   {ParamWeight},
	sidePlankEx: {ParamWeight},

	// CategoryDistTime с опциональным весом
	joggingEx: {ParamWeight},
	walkingEx: {ParamWeight},
}

var unitSuffixByLang = map[language]map[string]UnitDef{
	langRU: {
		// Вес — килограммы
		"кг":          {ParamType: ParamWeight, Multiplier: 1},
		"кило":        {ParamType: ParamWeight, Multiplier: 1},
		"килограмм":   {ParamType: ParamWeight, Multiplier: 1},
		"килограмма":  {ParamType: ParamWeight, Multiplier: 1},
		"килограммов": {ParamType: ParamWeight, Multiplier: 1},
		"килограм":    {ParamType: ParamWeight, Multiplier: 1},
		"килограма":   {ParamType: ParamWeight, Multiplier: 1},
		"килограмов":  {ParamType: ParamWeight, Multiplier: 1},

		// Вес — граммы
		"г":       {ParamType: ParamWeight, Multiplier: 0.001},
		"гр":      {ParamType: ParamWeight, Multiplier: 0.001},
		"грамм":   {ParamType: ParamWeight, Multiplier: 0.001},
		"грамма":  {ParamType: ParamWeight, Multiplier: 0.001},
		"граммов": {ParamType: ParamWeight, Multiplier: 0.001},
		"грам":    {ParamType: ParamWeight, Multiplier: 0.001},
		"грама":   {ParamType: ParamWeight, Multiplier: 0.001},
		"грамов":  {ParamType: ParamWeight, Multiplier: 0.001},

		// Дистанция — километры
		"км":         {ParamType: ParamDistance, Multiplier: 1000},
		"километр":   {ParamType: ParamDistance, Multiplier: 1000},
		"километра":  {ParamType: ParamDistance, Multiplier: 1000},
		"километров": {ParamType: ParamDistance, Multiplier: 1000},
		"киламетр":   {ParamType: ParamDistance, Multiplier: 1000},
		"киламетра":  {ParamType: ParamDistance, Multiplier: 1000},
		"киламетров": {ParamType: ParamDistance, Multiplier: 1000},
		"келометр":   {ParamType: ParamDistance, Multiplier: 1000},
		"келометра":  {ParamType: ParamDistance, Multiplier: 1000},
		"келометров": {ParamType: ParamDistance, Multiplier: 1000},
		"келаметр":   {ParamType: ParamDistance, Multiplier: 1000},
		"келаметра":  {ParamType: ParamDistance, Multiplier: 1000},
		"келаметров": {ParamType: ParamDistance, Multiplier: 1000},

		// Дистанция — метры
		"м":      {ParamType: ParamDistance, Multiplier: 1},
		"метр":   {ParamType: ParamDistance, Multiplier: 1},
		"метра":  {ParamType: ParamDistance, Multiplier: 1},
		"метров": {ParamType: ParamDistance, Multiplier: 1},

		// Время — часы
		"ч":     {ParamType: ParamDuration, Multiplier: 3600},
		"час":   {ParamType: ParamDuration, Multiplier: 3600},
		"часа":  {ParamType: ParamDuration, Multiplier: 3600},
		"часов": {ParamType: ParamDuration, Multiplier: 3600},
		"чиса":  {ParamType: ParamDuration, Multiplier: 3600},
		"чисов": {ParamType: ParamDuration, Multiplier: 3600},
		"чеса":  {ParamType: ParamDuration, Multiplier: 3600},
		"чесов": {ParamType: ParamDuration, Multiplier: 3600},

		// Время — минуты
		"мин":    {ParamType: ParamDuration, Multiplier: 60},
		"минут":  {ParamType: ParamDuration, Multiplier: 60},
		"минута": {ParamType: ParamDuration, Multiplier: 60},
		"минуты": {ParamType: ParamDuration, Multiplier: 60},
		"минуту": {ParamType: ParamDuration, Multiplier: 60},

		// Время — секунды
		"сек":     {ParamType: ParamDuration, Multiplier: 1},
		"с":       {ParamType: ParamDuration, Multiplier: 1},
		"секунд":  {ParamType: ParamDuration, Multiplier: 1},
		"секунда": {ParamType: ParamDuration, Multiplier: 1},
		"секунды": {ParamType: ParamDuration, Multiplier: 1},
		"секунду": {ParamType: ParamDuration, Multiplier: 1},
		"сикунд":  {ParamType: ParamDuration, Multiplier: 1},
		"сикунда": {ParamType: ParamDuration, Multiplier: 1},
		"сикунды": {ParamType: ParamDuration, Multiplier: 1},
		"сикунду": {ParamType: ParamDuration, Multiplier: 1},

		// Повторения
		"раз":        {ParamType: ParamCount, Multiplier: 1},
		"р":          {ParamType: ParamCount, Multiplier: 1},
		"повтор":     {ParamType: ParamCount, Multiplier: 1},
		"повтора":    {ParamType: ParamCount, Multiplier: 1},
		"повторов":   {ParamType: ParamCount, Multiplier: 1},
		"повторение": {ParamType: ParamCount, Multiplier: 1},
		"повторения": {ParamType: ParamCount, Multiplier: 1},
		"повторений": {ParamType: ParamCount, Multiplier: 1},
	},
	langEN: {
		// Weight — kilograms
		"kg":          {ParamType: ParamWeight, Multiplier: 1},
		"kgs":         {ParamType: ParamWeight, Multiplier: 1},
		"kilo":        {ParamType: ParamWeight, Multiplier: 1},
		"kilos":       {ParamType: ParamWeight, Multiplier: 1},
		"kilogram":    {ParamType: ParamWeight, Multiplier: 1},
		"kilograms":   {ParamType: ParamWeight, Multiplier: 1},
		"kilogramme":  {ParamType: ParamWeight, Multiplier: 1},
		"killo":       {ParamType: ParamWeight, Multiplier: 1},
		"killos":      {ParamType: ParamWeight, Multiplier: 1},
		"killogram":   {ParamType: ParamWeight, Multiplier: 1},
		"killograms":  {ParamType: ParamWeight, Multiplier: 1},
		"killogramme": {ParamType: ParamWeight, Multiplier: 1},

		// Weight — pounds
		"lbs":    {ParamType: ParamWeight, Multiplier: 0.453592},
		"lb":     {ParamType: ParamWeight, Multiplier: 0.453592},
		"pound":  {ParamType: ParamWeight, Multiplier: 0.453592},
		"pounds": {ParamType: ParamWeight, Multiplier: 0.453592},

		// Weight — grams
		"g":      {ParamType: ParamWeight, Multiplier: 0.001},
		"gs":     {ParamType: ParamWeight, Multiplier: 0.001},
		"gram":   {ParamType: ParamWeight, Multiplier: 0.001},
		"grams":  {ParamType: ParamWeight, Multiplier: 0.001},
		"gramm":  {ParamType: ParamWeight, Multiplier: 0.001},
		"gramms": {ParamType: ParamWeight, Multiplier: 0.001},

		// Distance — kilometers
		"km":          {ParamType: ParamDistance, Multiplier: 1000},
		"kms":         {ParamType: ParamDistance, Multiplier: 1000},
		"kilometer":   {ParamType: ParamDistance, Multiplier: 1000},
		"kilometers":  {ParamType: ParamDistance, Multiplier: 1000},
		"kilometre":   {ParamType: ParamDistance, Multiplier: 1000},
		"kilometres":  {ParamType: ParamDistance, Multiplier: 1000},
		"killometer":  {ParamType: ParamDistance, Multiplier: 1000},
		"killometers": {ParamType: ParamDistance, Multiplier: 1000},
		"killometre":  {ParamType: ParamDistance, Multiplier: 1000},
		"killometres": {ParamType: ParamDistance, Multiplier: 1000},

		// Distance — meters
		"m":      {ParamType: ParamDistance, Multiplier: 1},
		"meter":  {ParamType: ParamDistance, Multiplier: 1},
		"meters": {ParamType: ParamDistance, Multiplier: 1},
		"metre":  {ParamType: ParamDistance, Multiplier: 1},
		"metres": {ParamType: ParamDistance, Multiplier: 1},

		// Distance — miles
		"mi":     {ParamType: ParamDistance, Multiplier: 1609.34},
		"mile":   {ParamType: ParamDistance, Multiplier: 1609.34},
		"miles":  {ParamType: ParamDistance, Multiplier: 1609.34},
		"mille":  {ParamType: ParamDistance, Multiplier: 1609.34},
		"milles": {ParamType: ParamDistance, Multiplier: 1609.34},

		// Duration — hours
		"h":     {ParamType: ParamDuration, Multiplier: 3600},
		"hr":    {ParamType: ParamDuration, Multiplier: 3600},
		"hrs":   {ParamType: ParamDuration, Multiplier: 3600},
		"hour":  {ParamType: ParamDuration, Multiplier: 3600},
		"hours": {ParamType: ParamDuration, Multiplier: 3600},

		// Duration — minutes
		"min":     {ParamType: ParamDuration, Multiplier: 60},
		"mins":    {ParamType: ParamDuration, Multiplier: 60},
		"minute":  {ParamType: ParamDuration, Multiplier: 60},
		"minutes": {ParamType: ParamDuration, Multiplier: 60},

		// Duration — seconds
		"sec":     {ParamType: ParamDuration, Multiplier: 1},
		"secs":    {ParamType: ParamDuration, Multiplier: 1},
		"s":       {ParamType: ParamDuration, Multiplier: 1},
		"second":  {ParamType: ParamDuration, Multiplier: 1},
		"seconds": {ParamType: ParamDuration, Multiplier: 1},

		// Reps
		"reps":  {ParamType: ParamCount, Multiplier: 1},
		"rep":   {ParamType: ParamCount, Multiplier: 1},
		"x":     {ParamType: ParamCount, Multiplier: 1},
		"time":  {ParamType: ParamCount, Multiplier: 1},
		"times": {ParamType: ParamCount, Multiplier: 1},
	},
}
