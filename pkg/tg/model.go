package tg

import "strings"

type cmd int

const (
	unknownCmd = iota
	addCmd
	showCmd
	helpCmd
)

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

func (e Exercise) hasDistance() bool {
	_, ok := hasDistance[e]
	return ok
}

func (e Exercise) hasWeight() bool {
	_, ok := hasWeight[e]
	return ok
}

const (
	unknownEx      Exercise = "неизвестно"
	pullUpEx       Exercise = "подтягивание"
	muscleUpEx     Exercise = "выход"
	pushUpEx       Exercise = "отжимание"
	dipsEx         Exercise = "брусья"
	pressEx        Exercise = "пресс"
	squatsEx       Exercise = "приседания"
	lungesEx       Exercise = "выпады"
	burpeeEx       Exercise = "бёрпи"
	skippingRopeEx Exercise = "скакалка"
	joggingEx      Exercise = "бег"
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

func exercises() Exercises {
	return Exercises{
		pullUpEx,
		muscleUpEx,
		pushUpEx,
		dipsEx,
		pressEx,
		squatsEx,
		lungesEx,
		burpeeEx,
		skippingRopeEx,
		joggingEx,
	}
}

var (
	hasDistance = map[Exercise]struct{}{
		joggingEx: {},
	}
	hasWeight = map[Exercise]struct{}{
		pullUpEx:   {},
		muscleUpEx: {},
		dipsEx:     {},
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

		// pressEx
		"пресс":    pressEx,
		"прес":     pressEx,
		"пресуха":  pressEx,
		"прессуха": pressEx,

		// squatsEx
		"приседания": squatsEx,
		"приседаня":  squatsEx,
		"приседание": squatsEx,
		"приседане":  squatsEx,

		// lungesEx
		"выпады": lungesEx,
		"выпад":  lungesEx,

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
)

type Response struct {
	Message string `json:"message"`
}
