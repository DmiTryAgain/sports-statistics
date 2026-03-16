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
	plankEx   Exercise = "plank"

	benchPressEx       Exercise = "benchPress"
	deadliftEx         Exercise = "deadlift"
	barbellSquatEx     Exercise = "barbellSquat"
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

			// barbellSquatEx
			"присед со штангой":  barbellSquatEx,
			"приседы со штангой": barbellSquatEx,
			"присед штанга":      barbellSquatEx,

			// plankEx
			"планка":  plankEx,
			"планку":  plankEx,
			"планки":  plankEx,
			"планке":  plankEx,
			"планкой": plankEx,

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
			"верхний блок":        latPulldownEx,
			"верхняя тяга":        latPulldownEx,

			// legPressEx
			"жим ногами": legPressEx,
			"жим нагами": legPressEx,
			"жым ногами": legPressEx,
			"жым нагами": legPressEx,

			// preacherCurlEx
			"скамья скотта":   preacherCurlEx,
			"скамья скота":    preacherCurlEx,
			"скотта":          preacherCurlEx,
			"скота":           preacherCurlEx,
			"скамейка скотта": preacherCurlEx,
			"скамейка скота":  preacherCurlEx,

			// shoulderPressEx
			"жим стоя":        shoulderPressEx,
			"жым стоя":        shoulderPressEx,
			"армейский жим":   shoulderPressEx,
			"армейский жым":   shoulderPressEx,
			"армейскый жим":   shoulderPressEx,
			"армейскый жым":   shoulderPressEx,
			"жим над головой": shoulderPressEx,

			// bentOverRowEx
			"тяга в наклоне":        bentOverRowEx,
			"тяга штанги в наклоне": bentOverRowEx,
			"тяга штанги в наклони": bentOverRowEx,
			"тяга в наклони":        bentOverRowEx,

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
			"тяга нижнева блока":  seatedRowEx,
			"нижний блок":         seatedRowEx,
			"нижняя тяга":         seatedRowEx,
			"горизонтальная тяга": seatedRowEx,

			// chestFlyEx
			"сведение рук":   chestFlyEx,
			"сведения рук":   chestFlyEx,
			"бабочка":        chestFlyEx,
			"бабачка":        chestFlyEx,
			"разведение рук": chestFlyEx,
			"разведения рук": chestFlyEx,

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
			"махи":                      lateralRaiseEx,
			"разводка гантелей":         lateralRaiseEx,
			"разводка":                  lateralRaiseEx,
			"махи в стороны":            lateralRaiseEx,
			"подъём гантелей в стороны": lateralRaiseEx,
			"подъем гантелей в стороны": lateralRaiseEx,

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

			// benchPressEx
			"bench":       benchPressEx,
			"bench press": benchPressEx,
			"benchpress":  benchPressEx,

			// deadliftEx
			"deadlift":  deadliftEx,
			"deadlifts": deadliftEx,
			"dead lift": deadliftEx,

			// barbellSquatEx
			"barbell squat":  barbellSquatEx,
			"barbell squats": barbellSquatEx,
			"barbellsquat":   barbellSquatEx,

			// plankEx
			"plank":  plankEx,
			"planks": plankEx,

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
			plankEx:            "планка",
			benchPressEx:       "жим лёжа",
			deadliftEx:         "становая тяга",
			barbellSquatEx:     "присед со штангой",
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
			plankEx:            "plank",
			benchPressEx:       "bench press",
			deadliftEx:         "deadlift",
			barbellSquatEx:     "barbell squat",
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
	tableWeightCol
	tableDistCol
	tableTimeCol
	weightRequired
	durationRequired
	distanceRequired
	paramInvalid
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
			tableWeightCol:   "вес",
			tableDistCol:     "дистанция",
			tableTimeCol:     "время",
			weightRequired:   "Для этого упражнения нужно указать вес. Пример: жим 80кг 10",
			durationRequired: "Для этого упражнения нужно указать время. Пример: планка 90сек",
			distanceRequired: "Для этого упражнения нужно указать дистанцию. Пример: бег 5км 25мин",
			paramInvalid:     "Не удалось распознать параметр: %s",
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
				"*Поддерживаемые периоды:* сегодня, вчера, позавчера, неделя, месяц, год, всё время.\n" +
				"Также можно указать точную дату или интервал:\n" +
				"`за 15.10.2025`\n" +
				"`за 01.10.2025-10.10.2025`\n\n" +
				"Полный пример:\n" +
				"`@%[1]s покажи подтягивания отжимания за сегодня за 01.10.2025-10.10.2025`\n",
			helpHelpMsg: "Помощь к команде помощи не предусмотрена. Надо ж было додуматься попросить помощь к команде помощи🤔",
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
			tableWeightCol:   "weight",
			tableDistCol:     "distance",
			tableTimeCol:     "time",
			weightRequired:   "Weight is required for this exercise. Example: bench 80kg 10",
			durationRequired: "Duration is required for this exercise. Example: plank 90sec",
			distanceRequired: "Distance is required for this exercise. Example: run 5km 25min",
			paramInvalid:     "Can't recognize parameter: %s",
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
				"*Supported periods:* today, yesterday, week, month, year, all.\n" +
				"You can also specify an exact date or a range:\n" +
				"`for 15.10.2025`\n" +
				"`for 01.10.2025-10.10.2025`\n\n" +
				"Full example:\n" +
				"`@%[1]s show pull-ups push-ups for today for 01.10.2025-10.10.2025`\n",
			helpHelpMsg: "Help for the help command is not provided. How did you even think to ask for help on the help command?🤔",
			errMsg:      "❌ An error occurred. Try again later",
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

var exerciseCategoryMap = map[Exercise]ExerciseCategory{
	// CategoryRepsWeight
	benchPressEx:       CategoryRepsWeight,
	deadliftEx:         CategoryRepsWeight,
	barbellSquatEx:     CategoryRepsWeight,
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

	// CategoryDuration
	plankEx: CategoryDuration,
}

var unitSuffixByLang = map[language]map[string]UnitDef{
	langRU: {
		"кг":  {ParamType: ParamWeight, Multiplier: 1},
		"г":   {ParamType: ParamWeight, Multiplier: 0.001},
		"км":  {ParamType: ParamDistance, Multiplier: 1000},
		"м":   {ParamType: ParamDistance, Multiplier: 1},
		"ч":   {ParamType: ParamDuration, Multiplier: 3600},
		"мин": {ParamType: ParamDuration, Multiplier: 60},
		"сек": {ParamType: ParamDuration, Multiplier: 1},
		"с":   {ParamType: ParamDuration, Multiplier: 1},
		"раз": {ParamType: ParamCount, Multiplier: 1},
		"р":   {ParamType: ParamCount, Multiplier: 1},
	},
	langEN: {
		"kg":   {ParamType: ParamWeight, Multiplier: 1},
		"lbs":  {ParamType: ParamWeight, Multiplier: 0.453592},
		"lb":   {ParamType: ParamWeight, Multiplier: 0.453592},
		"g":    {ParamType: ParamWeight, Multiplier: 0.001},
		"km":   {ParamType: ParamDistance, Multiplier: 1000},
		"m":    {ParamType: ParamDistance, Multiplier: 1},
		"mi":   {ParamType: ParamDistance, Multiplier: 1609.34},
		"h":    {ParamType: ParamDuration, Multiplier: 3600},
		"hr":   {ParamType: ParamDuration, Multiplier: 3600},
		"min":  {ParamType: ParamDuration, Multiplier: 60},
		"sec":  {ParamType: ParamDuration, Multiplier: 1},
		"s":    {ParamType: ParamDuration, Multiplier: 1},
		"reps": {ParamType: ParamCount, Multiplier: 1},
		"rep":  {ParamType: ParamCount, Multiplier: 1},
		"x":    {ParamType: ParamCount, Multiplier: 1},
	},
}
