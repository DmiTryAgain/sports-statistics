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

func TestPeriodFromTextPeriod(t *testing.T) {
	// Среда 2025-10-15 12:00 UTC (weekday=3)
	now := time.Date(2025, 10, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name   string
		tp     string
		lang   language
		want   period
		wantOk bool
	}{
		{
			name: "lastWeek",
			tp:   "прошлую неделю",
			lang: langRU,
			want: period{
				from: time.Date(2025, 10, 6, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 10, 13, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			name: "weekBeforeLast",
			tp:   "позапрошлую неделю",
			lang: langRU,
			want: period{
				from: time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 10, 6, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			name: "lastMonth",
			tp:   "прошлый месяц",
			lang: langRU,
			want: period{
				from: time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			name: "monthBeforeLast",
			tp:   "позапрошлый месяц",
			lang: langRU,
			want: period{
				from: time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			name: "lastYear",
			tp:   "прошлый год",
			lang: langRU,
			want: period{
				from: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			name: "yearBeforeLast",
			tp:   "позапрошлый год",
			lang: langRU,
			want: period{
				from: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			// Понедельник уже прошёл на текущей неделе (пн 2025-10-13)
			name: "weekday monday (passed this week)",
			tp:   "пн",
			lang: langRU,
			want: period{
				from: time.Date(2025, 10, 13, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 10, 14, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			// Среда — сегодня
			name: "weekday wednesday (today)",
			tp:   "среда",
			lang: langRU,
			want: period{
				from: time.Date(2025, 10, 15, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 10, 15, 12, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			// Пятница ещё не наступила → прошлая неделя (2025-10-10)
			name: "weekday friday (not yet this week)",
			tp:   "пятнецу",
			lang: langRU,
			want: period{
				from: time.Date(2025, 10, 10, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 10, 11, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			name:   "unknown period",
			tp:     "неизвестно",
			lang:   langRU,
			wantOk: false,
		},
		{
			name: "lastWeek",
			tp:   "last week",
			lang: langEN,
			want: period{
				from: time.Date(2025, 10, 6, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 10, 13, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			name: "weekBeforeLast",
			tp:   "week before last",
			lang: langEN,
			want: period{
				from: time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 10, 6, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			name: "lastMonth",
			tp:   "last month",
			lang: langEN,
			want: period{
				from: time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			name: "monthBeforeLast",
			tp:   "month before last",
			lang: langEN,
			want: period{
				from: time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			name: "lastYear",
			tp:   "last year",
			lang: langEN,
			want: period{
				from: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			name: "yearBeforeLast",
			tp:   "year before last",
			lang: langEN,
			want: period{
				from: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			// Понедельник уже прошёл на текущей неделе (пн 2025-10-13)
			name: "weekday monday (passed this week)",
			tp:   "mon",
			lang: langEN,
			want: period{
				from: time.Date(2025, 10, 13, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 10, 14, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			// Среда — сегодня
			name: "weekday wednesday (today)",
			tp:   "wed",
			lang: langEN,
			want: period{
				from: time.Date(2025, 10, 15, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 10, 15, 12, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			// Пятница ещё не наступила → прошлая неделя (2025-10-10)
			name: "weekday friday (not yet this week)",
			tp:   "fri",
			lang: langEN,
			want: period{
				from: time.Date(2025, 10, 10, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 10, 11, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			name:   "unknown period",
			tp:     "unknown",
			lang:   langEN,
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := mh.periodByText(tt.tp, now, tt.lang)
			if ok != tt.wantOk {
				t.Fatalf("ok = %v, wantOk %v", ok, tt.wantOk)
			}
			if got != tt.want {
				t.Errorf("period = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriodForWeekday(t *testing.T) {
	// Фиксируем: среда 2025-10-15 14:30:00 UTC (weekday=3)
	now := time.Date(2025, 10, 15, 14, 30, 0, 0, time.UTC)

	tests := []struct {
		name   string
		word   string
		lang   language
		want   period
		wantOk bool
	}{
		{
			// Понедельник — уже прошёл на этой неделе (2 дня назад)
			name: "monday - passed this week",
			word: "пн",
			lang: langRU,
			want: period{
				from: time.Date(2025, 10, 13, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 10, 14, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			// Вторник — уже прошёл на этой неделе (1 день назад)
			name: "tuesday - passed this week",
			word: "вт",
			lang: langRU,
			want: period{
				from: time.Date(2025, 10, 14, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 10, 15, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			// Среда — это сегодня, to = текущий момент
			name: "wednesday - today",
			word: "ср",
			lang: langRU,
			want: period{
				from: time.Date(2025, 10, 15, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 10, 15, 14, 30, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			// Четверг — ещё не наступил, берём прошлую неделю (2025-10-09)
			name: "thursday - not yet this week, last week",
			word: "чт",
			lang: langRU,
			want: period{
				from: time.Date(2025, 10, 9, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 10, 10, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			// Пятница — ещё не наступила, берём прошлую неделю (2025-10-10)
			name: "friday - not yet this week, last week",
			word: "пт",
			lang: langRU,
			want: period{
				from: time.Date(2025, 10, 10, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 10, 11, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			// Суббота — ещё не наступила, берём прошлую неделю (2025-10-11)
			name: "saturday - not yet this week, last week",
			word: "сб",
			lang: langRU,
			want: period{
				from: time.Date(2025, 10, 11, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 10, 12, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			// Воскресенье — ещё не наступило, берём прошлую неделю (2025-10-12)
			name: "sunday - not yet this week, last week",
			word: "вс",
			lang: langRU,
			want: period{
				from: time.Date(2025, 10, 12, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2025, 10, 13, 0, 0, 0, 0, time.UTC),
			},
			wantOk: true,
		},
		{
			name:   "unknown period",
			word:   "unknown",
			lang:   langRU,
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := periodForWeekday(tt.word, tt.lang, now)
			if ok != tt.wantOk {
				t.Fatalf("periodForWeekday() ok = %v, wantOk %v", ok, tt.wantOk)
			}
			if got != tt.want {
				t.Errorf("from = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriodForWeekday_Sunday(t *testing.T) {
	// Воскресенье 2025-10-19 10:00 UTC (ISO weekday=7)
	now := time.Date(2025, 10, 19, 10, 0, 0, 0, time.UTC)

	// Воскресенье == сегодня
	got, ok := periodForWeekday("вс", langRU, now)
	if !ok {
		t.Fatal("expected ok for sunday == today")
	}
	want1 := period{
		from: time.Date(2025, 10, 19, 0, 0, 0, 0, time.UTC),
		to:   time.Date(2025, 10, 19, 10, 0, 0, 0, time.UTC),
	}
	if got != want1 {
		t.Errorf("from = %v, want %v", got, want1)
	}

	want2 := period{
		from: time.Date(2025, 10, 13, 0, 0, 0, 0, time.UTC),
		to:   time.Date(2025, 10, 14, 0, 0, 0, 0, time.UTC),
	}
	// Понедельник — ещё не наступил (воскресенье = последний день недели),
	// берём прошлую неделю (2025-10-13)
	got, ok = periodForWeekday("пн", langRU, now)
	if !ok {
		t.Fatal("expected ok for monday")
	}
	if got != want2 {
		t.Errorf("from = %v, want %v", got, want2)
	}
}

// TestParseNewPeriodKeywords проверяет, что новые ключевые слова периодов и дней недели
// парсятся без ошибок (exact даты зависят от now, поэтому проверяем только кол-во и отсутствие invalid).
func TestParseNewPeriodKeywords(t *testing.T) {
	tests := []struct {
		name   string
		rawMsg string
		lang   language
		wantEx Exercises
	}{
		{name: "last week ru", rawMsg: "подтягивания за прошлую неделю", lang: langRU, wantEx: Exercises{pullUpEx}},
		{name: "week before last ru", rawMsg: "подтягивания за позапрошлую неделю", lang: langRU, wantEx: Exercises{pullUpEx}},
		{name: "last month ru", rawMsg: "подтягивания за прошлый месяц", lang: langRU, wantEx: Exercises{pullUpEx}},
		{name: "month before last ru", rawMsg: "подтягивания за позапрошлый месяц", lang: langRU, wantEx: Exercises{pullUpEx}},
		{name: "last year ru", rawMsg: "подтягивания за прошлый год", lang: langRU, wantEx: Exercises{pullUpEx}},
		{name: "year before last ru", rawMsg: "подтягивания за позапрошлый год", lang: langRU, wantEx: Exercises{pullUpEx}},
		{name: "monday ru", rawMsg: "подтягивания за понедельник", lang: langRU, wantEx: Exercises{pullUpEx}},
		{name: "пн abbreviation ru", rawMsg: "подтягивания за пн", lang: langRU, wantEx: Exercises{pullUpEx}},
		{name: "среда ru", rawMsg: "подтягивания за среду", lang: langRU, wantEx: Exercises{pullUpEx}},
		{name: "вс ru", rawMsg: "подтягивания за вс", lang: langRU, wantEx: Exercises{pullUpEx}},
		{name: "last week en", rawMsg: "pull ups for last week", lang: langEN, wantEx: Exercises{pullUpEx}},
		{name: "last month en", rawMsg: "pull ups for last month", lang: langEN, wantEx: Exercises{pullUpEx}},
		{name: "last year en", rawMsg: "pull ups for last year", lang: langEN, wantEx: Exercises{pullUpEx}},
		{name: "friday en", rawMsg: "pull ups for friday", lang: langEN, wantEx: Exercises{pullUpEx}},
		{name: "mon en", rawMsg: "pull ups for mon", lang: langEN, wantEx: Exercises{pullUpEx}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEx, gotPeriods, gotInvalid, err := mh.parseRawMsgAsExercisesAndPeriods(t.Context(), tt.rawMsg, tt.lang)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(gotEx, tt.wantEx) {
				t.Errorf("exercises = %v, want %v", gotEx, tt.wantEx)
			}
			if len(gotPeriods) != 1 {
				t.Errorf("expected 1 period, got %d: %v", len(gotPeriods), gotPeriods)
			}
			if len(gotInvalid) != 0 {
				t.Errorf("unexpected invalid periods: %v", gotInvalid)
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
			name:         "multiword ex where first word is different ex ru (жим ногами)",
			words:        []string{"жим", "ногами", "100кг", "10"},
			lang:         langRU,
			wantExercise: legPressEx,
			wantOk:       true,
			wantExIdx:    1,
		},
		{
			name:         "multiword ex where first word is different ex ru (жим стоя)",
			words:        []string{"жим", "стоя", "50кг", "8"},
			lang:         langRU,
			wantExercise: shoulderPressEx,
			wantOk:       true,
			wantExIdx:    1,
		},
		{
			name:         "short ex still works when no longer match ru (жим)",
			words:        []string{"жим", "80кг", "10"},
			lang:         langRU,
			wantExercise: benchPressEx,
			wantOk:       true,
			wantExIdx:    0,
		},
		{
			name:         "multiword ex where first word is not ex ru (армейский жим)",
			words:        []string{"армейский", "жим", "40кг", "6"},
			lang:         langRU,
			wantExercise: shoulderPressEx,
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
		{"weight g ru", "500г", langRU, 0.5, ParamWeight, true, false},
		{"distance km ru", "5км", langRU, 5000, ParamDistance, true, false},
		{"distance km float ru", "5.5км", langRU, 5500, ParamDistance, true, false},
		{"duration min ru", "30мин", langRU, 1800, ParamDuration, true, false},
		{"duration sec ru", "90сек", langRU, 90, ParamDuration, true, false},
		{"count explicit ru", "10раз", langRU, 10, ParamCount, true, false},
		{"bare number ru", "10", langRU, 10, 0, false, false},
		{"bare float ru", "5.5", langRU, 5.5, 0, false, false},
		{"weight kg en", "80kg", langEN, 80, ParamWeight, true, false},
		{"weight lbs en", "150lbs", langEN, 150 * 0.453592, ParamWeight, true, false},
		{"distance km en", "5km", langEN, 5000, ParamDistance, true, false},
		{"duration min en", "30min", langEN, 1800, ParamDuration, true, false},
		{"not a number", "abc", langRU, 0, 0, false, true},
		{"unknown suffix", "80xyz", langRU, 0, 0, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pv, err := parseValueWithUnit(tt.word, tt.lang)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseValueWithUnit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if pv.Value != tt.wantValue {
				t.Errorf("parseValueWithUnit() value = %v, want %v", pv.Value, tt.wantValue)
			}
			if pv.HasUnit() != tt.wantHas {
				t.Errorf("parseValueWithUnit() hasUnit = %v, want %v", pv.HasUnit(), tt.wantHas)
			}
			if pv.HasUnit() && pv.Unit.ParamType != tt.wantType {
				t.Errorf("parseValueWithUnit() paramType = %v, want %v", pv.Unit.ParamType, tt.wantType)
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
			name:  "dist only for dist time is valid",
			words: []string{"5км"}, category: CategoryDistTime, lang: langRU,
			wantErr:  false,
			wantDist: ptr(5000.0),
		},
		{
			name:  "duration only for dist time is valid",
			words: []string{"25мин"}, category: CategoryDistTime, lang: langRU,
			wantErr: false,
			wantDur: ptr(1500.0),
		},
		{
			name:  "nothing for dist time is invalid",
			words: []string{}, category: CategoryDistTime, lang: langRU,
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
		// CategoryDistTime: хотя бы одно из dist/time
		{
			name:    "jogging dist only ru",
			rawMsg:  "бег 5км",
			lang:    langRU,
			want:    messagesByLang[langRU][exAdded],
			wantErr: false,
		},
		{
			name:    "jogging duration only ru",
			rawMsg:  "бег 25мин",
			lang:    langRU,
			want:    messagesByLang[langRU][exAdded],
			wantErr: false,
		},
		{
			name:    "jogging nothing ru",
			rawMsg:  "бег",
			lang:    langRU,
			want:    messagesByLang[langRU][distOrTimeRequired],
			wantErr: false,
		},
		// Walking
		{
			name:    "walking dist ru",
			rawMsg:  "ходьба 3км",
			lang:    langRU,
			want:    messagesByLang[langRU][exAdded],
			wantErr: false,
		},
		{
			name:    "walking duration en",
			rawMsg:  "walking 30min",
			lang:    langEN,
			want:    messagesByLang[langEN][exAdded],
			wantErr: false,
		},
		// CategoryDurationWeight
		{
			name:    "weight hold ru",
			rawMsg:  "удержание 40кг 30сек",
			lang:    langRU,
			want:    messagesByLang[langRU][exAdded],
			wantErr: false,
		},
		{
			name:    "weight hold missing weight ru",
			rawMsg:  "удержание 30сек",
			lang:    langRU,
			want:    messagesByLang[langRU][weightAndDurationRequired],
			wantErr: false,
		},
		{
			name:    "weight hold missing duration ru",
			rawMsg:  "удержание 40кг",
			lang:    langRU,
			want:    messagesByLang[langRU][weightAndDurationRequired],
			wantErr: false,
		},
		// New duration exercises
		{
			name:    "wall sit ru",
			rawMsg:  "стульчик 60сек",
			lang:    langRU,
			want:    messagesByLang[langRU][exAdded],
			wantErr: false,
		},
		{
			name:    "hang ru",
			rawMsg:  "вис 45сек",
			lang:    langRU,
			want:    messagesByLang[langRU][exAdded],
			wantErr: false,
		},
		// Squat with optional weight (text input)
		{
			name:    "squat with weight ru",
			rawMsg:  "приседания 80кг 10",
			lang:    langRU,
			want:    messagesByLang[langRU][exAdded],
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
