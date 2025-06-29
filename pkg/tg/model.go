package tg

import (
	"strings"
	"time"

	"github.com/DmiTryAgain/sports-statistics/pkg/db"
)

type allChecker interface {
	isAll() bool
}

type language string

type cmd int

type Exercise string

func (e Exercise) String() string { return string(e) }

func (e Exercise) isAll() bool {
	return e == allEx
}

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
	_, ok := exHasCnt[e]
	return ok
}

type Exercises []Exercise

func (e Exercises) StringSlice() []string {
	res := make([]string, len(e))
	for i := range e {
		res[i] = e[i].String()
	}

	return res
}

func textContainsAllExerciseWords(text string, lang language) bool {
	return textContainsSubstringInMapInAllValByLang(text, exerciseByLang[lang])
}

type textPeriod string

func (tp textPeriod) isAll() bool {
	return tp == allPeriod
}

func textContainsAllPeriodWords(text string, lang language) bool {
	return textContainsSubstringInMapInAllValByLang(text, periodByLang[lang])
}

func textContainsSubstringInMapInAllValByLang[T allChecker](text string, m map[string]T) bool {
	for i, v := range m {
		if strings.Contains(text, i) && v.isAll() {
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

type GroupedStatistic struct {
	db.GroupedStatistic
	TranslatedExercise string
}

func NewGroupedStatistic(in db.GroupedStatistic, lang language) GroupedStatistic {
	return GroupedStatistic{
		GroupedStatistic:   in,
		TranslatedExercise: exTextByLang[lang][Exercise(in.Exercise)],
	}
}

func NewGroupedStatisticList(in []db.GroupedStatistic, lang language) []GroupedStatistic {
	res := make([]GroupedStatistic, len(in))
	for i := range in {
		res[i] = NewGroupedStatistic(in[i], lang)
	}

	return res
}
