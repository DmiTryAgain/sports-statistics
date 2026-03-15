package tg

import (
	"context"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/DmiTryAgain/sports-statistics/pkg/db"

	"github.com/go-pg/pg/v10"
	"github.com/vmkteam/embedlog"
)

var (
	dbConn = env("DB_CONN", "postgres://postgres:postgres@localhost:5432/sport_statsrv?sslmode=disable")
	botCfg = Bot{
		Token:       "test",
		Name:        "test",
		ReplyFormat: "markdown",
		Debug:       false,
		Timeout:     Duration{Duration: 30 * time.Second},
	}
	pgConn *pg.DB
	dbc    db.DB

	mh *MessageHandler
)

func env(v, def string) string {
	if r := os.Getenv(v); r != "" {
		return r
	}

	return def
}

func initTestDB() db.DB {
	cfg, err := pg.ParseURL(dbConn)
	if err != nil {
		panic(err)
	}
	pgConn = pg.Connect(cfg)
	return db.New(pgConn)
}

func TestMain(m *testing.M) {
	//bot, err := tgbotapi.NewBotAPI(botCfg.Token)
	//if err != nil {
	//	panic(fmt.Errorf("create tgbot, err=%w", err))
	//}

	dbc = initTestDB()
	mh = New(context.Background(), embedlog.NewDevLogger(), dbc, db.NewStatisticRepo(pgConn), nil, botCfg)
	m.Run()
}

func TestMessageHandler_parseRawMsgAsExercisesAndPeriods(t *testing.T) {
	tests := []struct {
		name               string
		rawMsg             string
		lang               language
		wantExercises      Exercises
		wantPeriods        periods
		wantInvalidPeriods []string
		wantErr            bool
	}{
		{
			name:               "all ex all periods ru",
			rawMsg:             "всё за всё время",
			lang:               langRU,
			wantExercises:      nil,
			wantPeriods:        nil,
			wantInvalidPeriods: nil,
			wantErr:            false,
		},
		{
			name:               "push ups ex all periods ru",
			rawMsg:             "отжимания за всё время",
			lang:               langRU,
			wantExercises:      Exercises{pushUpEx},
			wantPeriods:        nil,
			wantInvalidPeriods: nil,
			wantErr:            false,
		},
		{
			name:               "push ups and pull ups ex all periods ru",
			rawMsg:             "отжимания патягивания за всё время",
			lang:               langRU,
			wantExercises:      Exercises{pushUpEx, pullUpEx},
			wantPeriods:        nil,
			wantInvalidPeriods: nil,
			wantErr:            false,
		},
		{
			name:               "push ups and pull ups ex all periods ru",
			rawMsg:             "push ups pull ups for all",
			lang:               langEN,
			wantExercises:      Exercises{pushUpEx, pullUpEx},
			wantPeriods:        nil,
			wantInvalidPeriods: nil,
			wantErr:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExercises, gotPeriods, gotInvalidPeriods, err := mh.parseRawMsgAsExercisesAndPeriods(t.Context(), tt.rawMsg, tt.lang)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRawMsgAsExercisesAndPeriods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotExercises, tt.wantExercises) {
				t.Errorf("parseRawMsgAsExercisesAndPeriods() gotExercises = %v, want %v", gotExercises, tt.wantExercises)
			}
			if !reflect.DeepEqual(gotPeriods, tt.wantPeriods) {
				t.Errorf("parseRawMsgAsExercisesAndPeriods() gotPeriods = %v, want %v", gotPeriods, tt.wantPeriods)
			}
			if !reflect.DeepEqual(gotInvalidPeriods, tt.wantInvalidPeriods) {
				t.Errorf("parseRawMsgAsExercisesAndPeriods() gotInvalidPeriods = %v, want %v", gotInvalidPeriods, tt.wantInvalidPeriods)
			}
		})
	}
}

func TestMessageHandler_multiWordsEx(t *testing.T) {
	tests := []struct {
		name         string
		words        []string
		lang         language
		wantExercise Exercise
		wantOk       bool
		wantExIdx    int
	}{
		{
			name:         "empty sclice",
			words:        []string{},
			lang:         langRU,
			wantExercise: "",
			wantOk:       false,
			wantExIdx:    0,
		},
		{
			name:         "one word ex ru",
			words:        []string{"подтягивания"},
			lang:         langRU,
			wantExercise: pullUpEx,
			wantOk:       true,
			wantExIdx:    0,
		},
		{
			name:         "two word ex ru with the first correct",
			words:        []string{"выход", "силы"},
			lang:         langRU,
			wantExercise: muscleUpEx,
			wantOk:       true,
			wantExIdx:    1,
		},
		{
			name:         "two word ex + another one exercise ru with the first correct",
			words:        []string{"выход", "силы", "подтягивания"},
			lang:         langRU,
			wantExercise: muscleUpEx,
			wantOk:       true,
			wantExIdx:    1,
		},
		{
			name:         "two word ex + another one exercise ru between",
			words:        []string{"выход", "подтягивания", "силы"},
			lang:         langRU,
			wantExercise: muscleUpEx,
			wantOk:       true,
			wantExIdx:    0,
		},
		{
			name:         "two word ex duplicated",
			words:        []string{"выход", "силы", "выход", "силы"},
			lang:         langRU,
			wantExercise: muscleUpEx,
			wantOk:       true,
			wantExIdx:    1,
		},
		{
			name:         "two word ex duplicated",
			words:        []string{"выход", "силы"},
			lang:         langRU,
			wantExercise: muscleUpEx,
			wantOk:       true,
			wantExIdx:    1,
		},
		{
			name:         "two word ex en with first incorrect but generally correct",
			words:        []string{"pull", "up"},
			lang:         langEN,
			wantExercise: pullUpEx,
			wantOk:       true,
			wantExIdx:    1,
		},
		{
			name:         "no valid exercises",
			words:        []string{"ahahahah", "invalid", "228"},
			lang:         langEN,
			wantExercise: "",
			wantOk:       false,
			wantExIdx:    2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExercise, gotOk, gotExIdx := mh.extractExerciseAndItsPosition(tt.words, tt.lang)
			if gotExercise != tt.wantExercise {
				t.Errorf("extractExerciseAndItsPosition() gotExercise = %v, want %v", gotExercise, tt.wantExercise)
			}
			if gotOk != tt.wantOk {
				t.Errorf("extractExerciseAndItsPosition() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if gotExIdx != tt.wantExIdx {
				t.Errorf("extractExerciseAndItsPosition() gotExIdx = %v, want %v", gotExIdx, tt.wantExIdx)
			}
		})
	}
}

func TestMessageHandler_clearRawMsg(t *testing.T) {
	tests := []struct {
		name   string
		rawMsg string
		want   string
	}{
		{
			name:   "empty string",
			rawMsg: "",
			want:   "",
		},
		{
			name:   "with mention valid",
			rawMsg: "@" + mh.cfg.Name + "сделал подтягивания 5",
			want:   "сделал подтягивания 5",
		},
		{
			name:   "with mention and puncts",
			rawMsg: "@" + mh.cfg.Name + "сделал подтягивания 5.0",
			want:   "сделал подтягивания 5.0",
		},
		{
			name:   "with mention and puncts",
			rawMsg: "@" + mh.cfg.Name + "сделал подтягивания 5.0",
			want:   "сделал подтягивания 5.0",
		},
		{
			name:   "without mention and extra puncts",
			rawMsg: "   сделал ,. .-подтягивания -    5.0",
			want:   "сделал подтягивания 5.0",
		},
		{
			name:   "with period",
			rawMsg: " покажи  ......всё за 15.10.2025  - 20.10.2025, 30.10.2025 ",
			want:   "покажи всё за 15.10.2025-20.10.2025 30.10.2025",
		},
		{
			name:   "with period",
			rawMsg: " покажи  ......всё за 15.10.2025  - 20.10.25, 30.10.2025 ",
			want:   "покажи всё за 15.10.2025-20.10.25 30.10.2025",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mh.clearRawMsg(tt.rawMsg); got != tt.want {
				t.Errorf("clearRawMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseValueWithUnit(t *testing.T) {
	tests := []struct {
		name      string
		word      string
		lang      language
		wantValue float64
		wantType  ParamType
		wantHas   bool
		wantErr   bool
	}{
		{"weight kg ru", "80кг", langRU, 80, ParamWeight, true, false},
		{"weight g ru", "500г", langRU, 500, ParamWeight, true, false},
		{"distance km ru", "5км", langRU, 5, ParamDistance, true, false},
		{"distance km float ru", "5.5км", langRU, 5.5, ParamDistance, true, false},
		{"duration min ru", "30мин", langRU, 30, ParamDuration, true, false},
		{"duration sec ru", "90сек", langRU, 90, ParamDuration, true, false},
		{"count explicit ru", "10раз", langRU, 10, ParamCount, true, false},
		{"bare number ru", "10", langRU, 10, 0, false, false},
		{"bare float ru", "5.5", langRU, 5.5, 0, false, false},
		{"weight kg en", "80kg", langEN, 80, ParamWeight, true, false},
		{"weight lbs en", "150lbs", langEN, 150, ParamWeight, true, false},
		{"distance km en", "5km", langEN, 5, ParamDistance, true, false},
		{"duration min en", "30min", langEN, 30, ParamDuration, true, false},
		{"not a number", "abc", langRU, 0, 0, false, true},
		{"unknown suffix", "80xyz", langRU, 0, 0, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, unit, hasUnit, err := parseValueWithUnit(tt.word, tt.lang)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseValueWithUnit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if value != tt.wantValue {
				t.Errorf("parseValueWithUnit() value = %v, want %v", value, tt.wantValue)
			}
			if hasUnit != tt.wantHas {
				t.Errorf("parseValueWithUnit() hasUnit = %v, want %v", hasUnit, tt.wantHas)
			}
			if hasUnit && unit.ParamType != tt.wantType {
				t.Errorf("parseValueWithUnit() paramType = %v, want %v", unit.ParamType, tt.wantType)
			}
		})
	}
}

func TestParseExerciseParams(t *testing.T) {
	tests := []struct {
		name       string
		words      []string
		category   ExerciseCategory
		lang       language
		wantCount  *float64
		wantWeight *float64
		wantDist   *float64
		wantDur    *float64
		wantErr    bool
	}{
		{
			name:  "reps backward compat",
			words: []string{"10"}, category: CategoryReps, lang: langRU,
			wantCount: ptr[float64](10),
		},
		{
			name:  "reps weight joined suffix",
			words: []string{"80кг", "10"}, category: CategoryRepsWeight, lang: langRU,
			wantCount: ptr[float64](10), wantWeight: ptr[float64](80),
		},
		{
			name:  "reps weight separate suffix",
			words: []string{"80", "кг", "10"}, category: CategoryRepsWeight, lang: langRU,
			wantCount: ptr[float64](10), wantWeight: ptr[float64](80),
		},
		{
			name:  "reps weight explicit count suffix",
			words: []string{"80кг", "10раз"}, category: CategoryRepsWeight, lang: langRU,
			wantCount: ptr[float64](10), wantWeight: ptr[float64](80),
		},
		{
			name:  "dist time composite",
			words: []string{"5км", "1ч", "30мин"}, category: CategoryDistTime, lang: langRU,
			wantDist: ptr[float64](5000), wantDur: ptr[float64](5400),
		},
		{
			name:  "dist time separate",
			words: []string{"5", "км", "25", "мин"}, category: CategoryDistTime, lang: langRU,
			wantDist: ptr[float64](5000), wantDur: ptr[float64](1500),
		},
		{
			name:  "duration only",
			words: []string{"90сек"}, category: CategoryDuration, lang: langRU,
			wantDur: ptr[float64](90),
		},
		{
			name:  "duration composite hours and min",
			words: []string{"1ч", "30мин"}, category: CategoryDuration, lang: langRU,
			wantDur: ptr[float64](5400),
		},
		{
			name:  "en weight reps",
			words: []string{"80kg", "10"}, category: CategoryRepsWeight, lang: langEN,
			wantCount: ptr[float64](10), wantWeight: ptr[float64](80),
		},
		{
			name:  "missing count for reps weight",
			words: []string{"80кг"}, category: CategoryRepsWeight, lang: langRU,
			wantErr: true,
		},
		{
			name:  "missing weight for reps weight",
			words: []string{"10"}, category: CategoryRepsWeight, lang: langRU,
			wantErr: true,
		},
		{
			name:  "missing duration for dist time",
			words: []string{"5км"}, category: CategoryDistTime, lang: langRU,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pp, err := parseExerciseParams(tt.words, tt.category, tt.lang)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseExerciseParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !floatPtrEqual(pp.Count, tt.wantCount) {
				t.Errorf("Count = %v, want %v", ptr(pp.Count), ptr(tt.wantCount))
			}
			if !floatPtrEqual(pp.WeightKg, tt.wantWeight) {
				t.Errorf("WeightKg = %v, want %v", ptr(pp.WeightKg), ptr(tt.wantWeight))
			}
			if !floatPtrEqual(pp.DistanceM, tt.wantDist) {
				t.Errorf("DistanceM = %v, want %v", ptr(pp.DistanceM), ptr(tt.wantDist))
			}
			if !floatPtrEqual(pp.DurationSec, tt.wantDur) {
				t.Errorf("DurationSec = %v, want %v", ptr(pp.DurationSec), ptr(tt.wantDur))
			}
		})
	}
}

func TestFormatWeight(t *testing.T) {
	tests := []struct {
		kg   float64
		lang language
		want string
	}{
		{80, langRU, "80кг"},
		{0.5, langRU, "500г"},
		{80.5, langRU, "80.5кг"},
		{80, langEN, "80kg"},
		{0.5, langEN, "500g"},
	}
	for _, tt := range tests {
		if got := formatWeight(tt.kg, tt.lang); got != tt.want {
			t.Errorf("formatWeight(%v, %v) = %v, want %v", tt.kg, tt.lang, got, tt.want)
		}
	}
}

func TestFormatDistance(t *testing.T) {
	tests := []struct {
		meters float64
		lang   language
		want   string
	}{
		{5000, langRU, "5км"},
		{800, langRU, "800м"},
		{2500, langRU, "2.5км"},
		{5000, langEN, "5km"},
	}
	for _, tt := range tests {
		if got := formatDistance(tt.meters, tt.lang); got != tt.want {
			t.Errorf("formatDistance(%v, %v) = %v, want %v", tt.meters, tt.lang, got, tt.want)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		sec  float64
		lang language
		want string
	}{
		{45, langRU, "45сек"},
		{1500, langRU, "25мин"},
		{5400, langRU, "1ч 30мин"},
		{1530, langRU, "25мин 30сек"},
		{3600, langRU, "1ч"},
		{45, langEN, "45sec"},
		{1500, langEN, "25min"},
		{5400, langEN, "1h 30min"},
	}
	for _, tt := range tests {
		if got := formatDuration(tt.sec, tt.lang); got != tt.want {
			t.Errorf("formatDuration(%v, %v) = %v, want %v", tt.sec, tt.lang, got, tt.want)
		}
	}
}

func floatPtrEqual(a, b *float64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func TestMessageHandler_handleAdd(t *testing.T) {
	t.Skip()

	ctx := t.Context()
	tests := []struct {
		name    string
		rawMsg  string
		lang    language
		want    string
		wantErr bool
	}{
		{
			name:    "add valid exercise ru",
			rawMsg:  "подтягивания 5",
			lang:    langRU,
			want:    messagesByLang[langRU][exAdded],
			wantErr: false,
		},
		{
			name:    "add valid double words with valid the first word exercise ru",
			rawMsg:  "выход силы 5",
			lang:    langRU,
			want:    messagesByLang[langRU][exAdded],
			wantErr: false,
		},
		{
			name:    "add valid double words with valid the first word exercise ru with float count",
			rawMsg:  "выход силы 5.5",
			lang:    langRU,
			want:    messagesByLang[langRU][exAdded],
			wantErr: false,
		},
		{
			name:    "add valid double words exercise with invalid the first word en",
			rawMsg:  "pull ups 5",
			lang:    langEN,
			want:    messagesByLang[langEN][exAdded],
			wantErr: false,
		},
		{
			name:    "invalid exercise ru",
			rawMsg:  "ыыыщ 5",
			lang:    langRU,
			want:    messagesByLang[langRU][cantRecognizeEx],
			wantErr: false,
		},
		{
			name:    "empty count ru",
			rawMsg:  "выход силы",
			lang:    langRU,
			want:    messagesByLang[langRU][cntRequired],
			wantErr: false,
		},
		{
			name:    "small count ru",
			rawMsg:  "выход силы 0",
			lang:    langRU,
			want:    messagesByLang[langRU][cntGE],
			wantErr: false,
		},
		// New exercises
		{
			name:    "bench press ru",
			rawMsg:  "жим 80кг 10",
			lang:    langRU,
			want:    messagesByLang[langRU][exAdded],
			wantErr: false,
		},
		{
			name:    "bench press en",
			rawMsg:  "bench 80kg 10",
			lang:    langEN,
			want:    messagesByLang[langEN][exAdded],
			wantErr: false,
		},
		{
			name:    "plank ru",
			rawMsg:  "планка 90сек",
			lang:    langRU,
			want:    messagesByLang[langRU][exAdded],
			wantErr: false,
		},
		{
			name:    "jogging ru",
			rawMsg:  "бег 5км 25мин",
			lang:    langRU,
			want:    messagesByLang[langRU][exAdded],
			wantErr: false,
		},
		{
			name:    "bench press missing weight ru",
			rawMsg:  "жим 10",
			lang:    langRU,
			want:    messagesByLang[langRU][weightRequired],
			wantErr: false,
		},
		{
			name:    "bench press missing count ru",
			rawMsg:  "жим 80кг",
			lang:    langRU,
			want:    messagesByLang[langRU][weightRequired],
			wantErr: false,
		},
		{
			name:    "plank missing duration ru",
			rawMsg:  "планка",
			lang:    langRU,
			want:    messagesByLang[langRU][durationRequired],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mh.handleAdd(ctx, tt.rawMsg, "testuser", tt.lang)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(got, tt.want) {
				t.Errorf("handleAdd() got = %v, want %v", got, tt.want)
			}
		})
	}
}
