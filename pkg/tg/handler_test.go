package tg

import (
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
	mh = New(embedlog.NewDevLogger(), dbc, db.NewStatisticRepo(pgConn), nil, botCfg)
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
			name:    "invalid count ru",
			rawMsg:  "выход силы 5..0",
			lang:    langRU,
			want:    messagesByLang[langRU][cntInvalid],
			wantErr: false,
		},
		{
			name:    "small count ru",
			rawMsg:  "выход силы 0",
			lang:    langRU,
			want:    messagesByLang[langRU][cntGE],
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
