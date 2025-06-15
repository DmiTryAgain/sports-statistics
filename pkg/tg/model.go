package tg

import (
	"strings"
	"time"

	"github.com/DmiTryAgain/sports-statistics/pkg/db"
)

const (
	unknownCmd cmd = iota
	addCmd
	showCmd
	helpCmd
)

type cmd int

var (
	cmdByWord = map[string]cmd{
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
		"разебал":   addCmd,
		"разыбал":   addCmd,

		// showCmd
		"покажи":     showCmd,
		"покажы":     showCmd,
		"пакажи":     showCmd,
		"пакажы":     showCmd,
		"покожы":     showCmd,
		"покожи":     showCmd,
		"статистика": showCmd,
		"стата":      showCmd,

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
	}
)

type Exercise string

func (e Exercise) String() string { return string(e) }

// TODO: implement later
//func (e Exercise) hasDistance() bool {
//	_, ok := hasDistance[e]
//	return ok
//}

// TODO: implement weight exercises
//func (e Exercise) hasWeight() bool {
//	_, ok := hasWeight[e]
//	return ok
//}

func (e Exercise) mustHaveCnt() bool {
	_, ok := hasCnt[e]
	return ok
}

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
	joggingEx      Exercise = "jogging"
	allEx          Exercise = "all"
)

type Exercises []Exercise

func (e Exercises) String() string {
	sb := strings.Builder{}
	for _, ex := range e {
		sb.WriteString(ex.String())
		sb.WriteString(", ")
	}

	return sb.String()
}

func (e Exercises) StringSlice() []string {
	res := make([]string, len(e))
	for i := range e {
		res[i] = e[i].String()
	}

	return res
}

func exercises() Exercises {
	return Exercises{
		pullUpEx,
		muscleUpEx,
		pushUpEx,
		dipsEx,
		absEx,
		squatEx,
		lungeEx,
		burpeeEx,
		skippingRopeEx,
		joggingEx,
	}
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
	hasCnt = map[Exercise]struct{}{
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

var (
	exerciseByWord = map[string]Exercise{
		//	pullUpEx
		"подтягивание": pullUpEx,
		"подтягивания": pullUpEx,
		"падтягивания": pullUpEx,
		"падтягивание": pullUpEx,
		"потягивание":  pullUpEx,
		"патягивание":  pullUpEx,
		"подтянулся":   pullUpEx,

		//	muscleUpEx
		"выход": muscleUpEx,

		//	pushUpEx
		"отжимание": pushUpEx,
		"отжимания": pushUpEx,
		"анжуманя":  pushUpEx,

		// dipsEx
		"брусья":  dipsEx,
		"брусьях": dipsEx,

		// absEx
		"пресс":    absEx,
		"прес":     absEx,
		"пресуха":  absEx,
		"прессуха": absEx,

		// squatEx
		"приседания": squatEx,
		"приседаня":  squatEx,
		"приседание": squatEx,
		"приседане":  squatEx,

		// lungeEx
		"выпады": lungeEx,
		"выпад":  lungeEx,

		// burpeeEx
		"бёрпи": burpeeEx,
		"берпи": burpeeEx,

		// skippingRopeEx
		"скакалка": skippingRopeEx,

		// joggingEx
		"бег":      joggingEx,
		"бегал":    joggingEx,
		"пробежал": joggingEx,
		"пробежка": joggingEx,
	}

	allExercisesByWord = map[string]Exercise{
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
	}
)

func textContainsAllExerciseWords(text string) bool {
	return textContainsSubstringInMap(text, allExercisesByWord)
}

const (
	todayPeriod              textPeriod = "сегодня"
	yesterdayPeriod          textPeriod = "вчера"
	dayBeforeYesterdayPeriod textPeriod = "позавчера"
	weekPeriod               textPeriod = "неделя"
	monthPeriod              textPeriod = "месяц"
	yearPeriod               textPeriod = "год"
	allPeriod                textPeriod = "всё время"
)

type textPeriod string

var (
	periodByWord = map[string]textPeriod{
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
	}

	allPeriodByWord = map[string]textPeriod{
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
	}
)

func textContainsAllPeriodWords(text string) bool {
	return textContainsSubstringInMap(text, allPeriodByWord)
}

func textContainsSubstringInMap[T any](text string, m map[string]T) bool {
	for i := range m {
		if strings.Contains(text, i) {
			return true
		}
	}

	return false
}

type period struct {
	from, to time.Time
}

func (p period) IsZero() bool {
	return p.from.IsZero() && p.to.IsZero()
}

func (p period) ToDB() db.Period {
	return db.Period{
		From: p.from,
		To:   p.to,
	}
}

type periods []period

func (ps periods) ToDB() []db.Period {
	res := make([]db.Period, len(ps))
	for i := range ps {
		res[i] = ps[i].ToDB()
	}

	return res
}
