package tg

import (
	"fmt"
	"testing"

	"github.com/DmiTryAgain/sports-statistics/pkg/db"
)

func TestExerciseInlineKeyboard_FirstPage(t *testing.T) {
	kb := exerciseInlineKeyboard(0, langRU, "add")

	// Считаем кнопки упражнений (все ряды кроме навигации и отмены)
	totalExButtons := 0
	for _, row := range kb.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil && len(*btn.CallbackData) > 3 && (*btn.CallbackData)[:3] == "ex|" {
				totalExButtons++
			}
		}
	}

	if totalExButtons != exercisesPerPage {
		t.Errorf("expected %d exercise buttons on first page, got %d", exercisesPerPage, totalExButtons)
	}

	// Проверяем, что первая кнопка — подтягивания
	firstBtn := kb.InlineKeyboard[0][0]
	if *firstBtn.CallbackData != "ex|jogging" {
		t.Errorf("expected first button callback ex|jogging, got %s", *firstBtn.CallbackData)
	}

	// Должна быть кнопка "Ещё >>" (pg|1|add)
	hasMore := false
	for _, row := range kb.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil && *btn.CallbackData == "pg|1|add" {
				hasMore = true
			}
		}
	}
	if !hasMore {
		t.Error("expected 'More' navigation button on first page")
	}

	// Не должно быть кнопки "Назад"
	hasBack := false
	for _, row := range kb.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil && len(*btn.CallbackData) > 3 && (*btn.CallbackData)[:3] == "pg|" && *btn.CallbackData != "pg|1|add" {
				hasBack = true
			}
		}
	}
	if hasBack {
		t.Error("should not have 'Back' button on first page")
	}

	// Должна быть кнопка отмены
	lastRow := kb.InlineKeyboard[len(kb.InlineKeyboard)-1]
	if *lastRow[0].CallbackData != "x" {
		t.Errorf("expected cancel button at the end, got %s", *lastRow[0].CallbackData)
	}
}

func TestExerciseInlineKeyboard_SecondPage(t *testing.T) {
	kb := exerciseInlineKeyboard(1, langRU, "add")

	// Должна быть кнопка "<< Назад" (pg|0|add)
	hasBack := false
	for _, row := range kb.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil && *btn.CallbackData == "pg|0|add" {
				hasBack = true
			}
		}
	}
	if !hasBack {
		t.Error("expected 'Back' navigation button on second page")
	}
}

func TestExerciseInlineKeyboard_LastPage(t *testing.T) {
	lastPage := (len(exerciseOrder) - 1) / exercisesPerPage
	kb := exerciseInlineKeyboard(lastPage, langRU, "add")

	// Считаем оставшиеся упражнения на последней странице
	remaining := len(exerciseOrder) - lastPage*exercisesPerPage
	totalExButtons := 0
	for _, row := range kb.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil && len(*btn.CallbackData) > 3 && (*btn.CallbackData)[:3] == "ex|" {
				totalExButtons++
			}
		}
	}
	if totalExButtons != remaining {
		t.Errorf("expected %d exercise buttons on last page, got %d", remaining, totalExButtons)
	}
}

func TestWeightInlineKeyboard_WithHistory(t *testing.T) {
	weights := []float64{100, 80, 60}
	kb := weightInlineKeyboard(weights, langRU)

	// Проверяем кнопки весов
	weightButtons := 0
	for _, row := range kb.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil && len(*btn.CallbackData) > 2 && (*btn.CallbackData)[:2] == "w|" {
				weightButtons++
			}
		}
	}
	if weightButtons != 3 {
		t.Errorf("expected 3 weight buttons, got %d", weightButtons)
	}

	// Проверяем callback_data
	firstBtn := kb.InlineKeyboard[0][0]
	if *firstBtn.CallbackData != "w|100" {
		t.Errorf("expected first weight callback w|100, got %s", *firstBtn.CallbackData)
	}

	// Проверяем наличие кнопки "Другой" и "Отмена"
	lastRow := kb.InlineKeyboard[len(kb.InlineKeyboard)-1]
	if *lastRow[0].CallbackData != "cu|w" {
		t.Errorf("expected custom input button cu|w, got %s", *lastRow[0].CallbackData)
	}
	if *lastRow[1].CallbackData != "x" {
		t.Errorf("expected cancel button x, got %s", *lastRow[1].CallbackData)
	}
}

func TestWeightInlineKeyboard_WithoutHistory(t *testing.T) {
	kb := weightInlineKeyboard(nil, langRU)

	// Без истории — стандартные значения: 20, 40, 60, 80, 100
	weightButtons := 0
	for _, row := range kb.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil && len(*btn.CallbackData) > 2 && (*btn.CallbackData)[:2] == "w|" {
				weightButtons++
			}
		}
	}
	if weightButtons != 5 {
		t.Errorf("expected 5 default weight buttons, got %d", weightButtons)
	}
}

func TestCountInlineKeyboard(t *testing.T) {
	kb := countInlineKeyboard(langRU)

	// 6 кнопок: 5, 8, 10, 12, 15, 20
	countButtons := 0
	for _, row := range kb.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil && len(*btn.CallbackData) > 2 && (*btn.CallbackData)[:2] == "c|" {
				countButtons++
			}
		}
	}
	if countButtons != 6 {
		t.Errorf("expected 6 count buttons, got %d", countButtons)
	}

	// Проверяем callback_data первой кнопки
	if *kb.InlineKeyboard[0][0].CallbackData != "c|5" {
		t.Errorf("expected first count callback c|5, got %s", *kb.InlineKeyboard[0][0].CallbackData)
	}
}

func TestDistanceInlineKeyboard(t *testing.T) {
	kb := distanceInlineKeyboard(langRU)

	distButtons := 0
	for _, row := range kb.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil && len(*btn.CallbackData) > 2 && (*btn.CallbackData)[:2] == "d|" {
				distButtons++
			}
		}
	}
	if distButtons != 6 {
		t.Errorf("expected 5 distance buttons, got %d", distButtons)
	}

	// Проверяем текст первой кнопки — 1000м = 1км
	if kb.InlineKeyboard[0][0].Text != "1км" {
		t.Errorf("expected first distance text 1км, got %s", kb.InlineKeyboard[0][0].Text)
	}
}

func TestDurationInlineKeyboard_Duration(t *testing.T) {
	kb := durationInlineKeyboard(CategoryDuration, langRU)

	durButtons := 0
	for _, row := range kb.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil && len(*btn.CallbackData) > 2 && (*btn.CallbackData)[:2] == "t|" {
				durButtons++
			}
		}
	}
	if durButtons != 6 {
		t.Errorf("expected 6 duration buttons for CategoryDuration, got %d", durButtons)
	}

	// Первая кнопка — 30сек
	if kb.InlineKeyboard[0][0].Text != "30сек" {
		t.Errorf("expected first duration text 30сек, got %s", kb.InlineKeyboard[0][0].Text)
	}
}

func TestDurationInlineKeyboard_DistTime(t *testing.T) {
	kb := durationInlineKeyboard(CategoryDistTime, langRU)

	durButtons := 0
	for _, row := range kb.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil && len(*btn.CallbackData) > 2 && (*btn.CallbackData)[:2] == "t|" {
				durButtons++
			}
		}
	}
	if durButtons != 6 {
		t.Errorf("expected 6 duration buttons for CategoryDistTime, got %d", durButtons)
	}

	// Первая кнопка — 900 = 15мин
	if kb.InlineKeyboard[0][0].Text != "15мин" {
		t.Errorf("expected first duration text 15мин, got %s", kb.InlineKeyboard[0][0].Text)
	}
}

func TestPeriodInlineKeyboard(t *testing.T) {
	kb := periodInlineKeyboard(langRU)

	periodButtons := 0
	for _, row := range kb.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil && len(*btn.CallbackData) > 3 && (*btn.CallbackData)[:3] == "sp|" {
				periodButtons++
			}
		}
	}
	// 9 периодов: today, yesterday, week, lastWeek, month, lastMonth, year, lastYear, all
	if periodButtons != 9 {
		t.Errorf("expected 9 period buttons, got %d", periodButtons)
	}

	// Проверяем первый callback
	if *kb.InlineKeyboard[0][0].CallbackData != "sp|today" {
		t.Errorf("expected first period callback sp|today, got %s", *kb.InlineKeyboard[0][0].CallbackData)
	}

	// Последняя строка — отмена
	lastRow := kb.InlineKeyboard[len(kb.InlineKeyboard)-1]
	if *lastRow[0].CallbackData != "x" {
		t.Errorf("expected cancel button, got %s", *lastRow[0].CallbackData)
	}
}

func TestShowExerciseInlineKeyboard(t *testing.T) {
	exercises := []string{"pullUp", "benchPress", "jogging"}
	kb := showExerciseInlineKeyboard(exercises, langRU)

	// 3 кнопки упражнений
	exButtons := 0
	for _, row := range kb.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil && len(*btn.CallbackData) > 3 && (*btn.CallbackData)[:3] == "se|" {
				exButtons++
			}
		}
	}
	if exButtons != 3 {
		t.Errorf("expected 3 exercise buttons, got %d", exButtons)
	}

	// Кнопка "Всё"
	hasAll := false
	for _, row := range kb.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil && *btn.CallbackData == "sa" {
				hasAll = true
			}
		}
	}
	if !hasAll {
		t.Error("expected 'All' button")
	}

	// Кнопка отмены
	lastRow := kb.InlineKeyboard[len(kb.InlineKeyboard)-1]
	if *lastRow[0].CallbackData != "x" {
		t.Errorf("expected cancel button, got %s", *lastRow[0].CallbackData)
	}
}

func TestQuickAddInlineKeyboard(t *testing.T) {
	freq := []db.FrequentExercise{
		{Exercise: "benchPress", Count: 10, Params: &db.StatisticParams{WeightKg: ptr[float64](80)}},
		{Exercise: "pullUp", Count: 12, Params: nil},
	}
	rows := quickAddInlineKeyboard(freq, langRU)

	if len(rows) != 2 {
		t.Fatalf("expected 2 quick-add rows, got %d", len(rows))
	}

	// Первая кнопка — жим
	if *rows[0][0].CallbackData != "qa|benchPress|10|80||" {
		t.Errorf("expected callback qa|benchPress|10|80||, got %s", *rows[0][0].CallbackData)
	}

	// Вторая кнопка — подтягивания
	if *rows[1][0].CallbackData != "qa|pullUp|12|||\x00"[:15] {
		// Проверяем что начинается с qa|pullUp|12
		data := *rows[1][0].CallbackData
		if len(data) < 12 || data[:12] != "qa|pullUp|12" {
			t.Errorf("expected callback starting with qa|pullUp|12, got %s", data)
		}
	}
}

func TestQuickAddInlineKeyboard_Empty(t *testing.T) {
	rows := quickAddInlineKeyboard(nil, langRU)
	if rows != nil {
		t.Errorf("expected nil for empty frequent, got %v", rows)
	}
}

func TestFormatQuickAddButton(t *testing.T) {
	tests := []struct {
		name string
		freq db.FrequentExercise
		lang language
		want string
	}{
		{
			name: "reps weight",
			freq: db.FrequentExercise{Exercise: "benchPress", Count: 10, Params: &db.StatisticParams{WeightKg: ptr[float64](80)}},
			lang: langRU,
			want: exTextByLang[langRU][benchPressEx] + " 80кг ×10",
		},
		{
			name: "reps only",
			freq: db.FrequentExercise{Exercise: "pullUp", Count: 12, Params: nil},
			lang: langRU,
			want: exTextByLang[langRU][pullUpEx] + " ×12",
		},
		{
			name: "duration",
			freq: db.FrequentExercise{Exercise: "plank", Count: 1, Params: &db.StatisticParams{DurationSec: ptr[float64](60)}},
			lang: langRU,
			want: exTextByLang[langRU][plankEx] + " " + formatDuration(60, langRU),
		},
		{
			name: "dist time",
			freq: db.FrequentExercise{Exercise: "jogging", Count: 1, Params: &db.StatisticParams{DistanceM: ptr[float64](5000), DurationSec: ptr[float64](1500)}},
			lang: langRU,
			want: exTextByLang[langRU][joggingEx] + " 5км 25мин",
		},
		{
			name: "reps weight en",
			freq: db.FrequentExercise{Exercise: "benchPress", Count: 10, Params: &db.StatisticParams{WeightKg: ptr[float64](80)}},
			lang: langEN,
			want: exTextByLang[langEN][benchPressEx] + " 80kg ×10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatQuickAddButton(tt.freq, tt.lang)
			if got != tt.want {
				t.Errorf("formatQuickAddButton() = %q, want %q", got, tt.want)
			}
		})
	}
}

func withQuickHint(confirm, quickCmd string, lang language) string {
	return confirm + "\n" + fmt.Sprintf(messagesByLang[lang][quickCopyHint], "`"+quickCmd+"`")
}

func TestFormatAddConfirmation(t *testing.T) {
	tests := []struct {
		name   string
		ex     Exercise
		cnt    float64
		params *db.StatisticParams
		lang   language
		want   string
	}{
		{
			name: "reps ru",
			ex:   pullUpEx,
			cnt:  12,
			lang: langRU,
			want: withQuickHint(
				fmt.Sprintf(messagesByLang[langRU][addedConfirmation], exTextByLang[langRU][pullUpEx], "×12"),
				"добавь подтягивания 12", langRU,
			),
		},
		{
			name:   "reps weight ru",
			ex:     benchPressEx,
			cnt:    10,
			params: &db.StatisticParams{WeightKg: ptr[float64](80)},
			lang:   langRU,
			want: withQuickHint(
				fmt.Sprintf(messagesByLang[langRU][addedConfirmation], exTextByLang[langRU][benchPressEx], "80кг × 10"),
				"добавь жим лёжа 80кг 10", langRU,
			),
		},
		{
			name:   "duration ru",
			ex:     plankEx,
			cnt:    1,
			params: &db.StatisticParams{DurationSec: ptr[float64](90)},
			lang:   langRU,
			want: withQuickHint(
				fmt.Sprintf(messagesByLang[langRU][addedConfirmation], exTextByLang[langRU][plankEx], formatDuration(90, langRU)),
				"добавь планка 1мин 30сек", langRU,
			),
		},
		{
			name:   "dist time ru",
			ex:     joggingEx,
			cnt:    1,
			params: &db.StatisticParams{DistanceM: ptr[float64](5000), DurationSec: ptr[float64](1500)},
			lang:   langRU,
			want: withQuickHint(
				fmt.Sprintf(messagesByLang[langRU][addedConfirmation], exTextByLang[langRU][joggingEx], "5км 25мин"),
				"добавь бег 5км 25мин", langRU,
			),
		},
		{
			name: "reps en",
			ex:   pullUpEx,
			cnt:  12,
			lang: langEN,
			want: withQuickHint(
				fmt.Sprintf(messagesByLang[langEN][addedConfirmation], exTextByLang[langEN][pullUpEx], "×12"),
				"add pull-ups 12", langEN,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatAddConfirmation(tt.ex, tt.cnt, tt.params, tt.lang)
			if got != tt.want {
				t.Errorf("formatAddConfirmation() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatQuickCopyCommand(t *testing.T) {
	tests := []struct {
		name   string
		ex     Exercise
		cnt    float64
		params *db.StatisticParams
		lang   language
		want   string
	}{
		{
			name: "reps only ru",
			ex:   pullUpEx,
			cnt:  10,
			lang: langRU,
			want: "добавь подтягивания 10",
		},
		{
			name:   "reps weight ru",
			ex:     benchPressEx,
			cnt:    8,
			params: &db.StatisticParams{WeightKg: ptr[float64](80)},
			lang:   langRU,
			want:   "добавь жим лёжа 80кг 8",
		},
		{
			name:   "dist time ru",
			ex:     joggingEx,
			cnt:    1,
			params: &db.StatisticParams{DistanceM: ptr[float64](5000), DurationSec: ptr[float64](1500)},
			lang:   langRU,
			want:   "добавь бег 5км 25мин",
		},
		{
			name:   "duration ru",
			ex:     plankEx,
			cnt:    1,
			params: &db.StatisticParams{DurationSec: ptr[float64](120)},
			lang:   langRU,
			want:   "добавь планка 2мин",
		},
		{
			name:   "duration weight ru",
			ex:     weightHoldEx,
			cnt:    1,
			params: &db.StatisticParams{WeightKg: ptr[float64](20), DurationSec: ptr[float64](60)},
			lang:   langRU,
			want:   "добавь удержание веса 20кг 1мин",
		},
		{
			name: "reps en",
			ex:   pushUpEx,
			cnt:  20,
			lang: langEN,
			want: "add push-ups 20",
		},
		{
			name:   "dist only en",
			ex:     walkingEx,
			cnt:    1,
			params: &db.StatisticParams{DistanceM: ptr[float64](3000)},
			lang:   langEN,
			want:   "add walking 3km",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatQuickCopyCommand(tt.ex, tt.cnt, tt.params, tt.lang)
			if got != tt.want {
				t.Errorf("formatQuickCopyCommand() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReplyKeyboard(t *testing.T) {
	kb := replyKeyboard(langRU)

	if len(kb.Keyboard) != 1 {
		t.Fatalf("expected 1 row, got %d", len(kb.Keyboard))
	}
	if len(kb.Keyboard[0]) != 3 {
		t.Fatalf("expected 3 buttons, got %d", len(kb.Keyboard[0]))
	}

	if kb.Keyboard[0][0].Text != "📝 Добавить" {
		t.Errorf("expected first button '📝 Добавить', got %s", kb.Keyboard[0][0].Text)
	}
	if kb.Keyboard[0][1].Text != "📊 Статистика" {
		t.Errorf("expected second button '📊 Статистика', got %s", kb.Keyboard[0][1].Text)
	}
	if kb.Keyboard[0][2].Text != "❓ Помощь" {
		t.Errorf("expected third button '❓ Помощь', got %s", kb.Keyboard[0][2].Text)
	}
}

func TestReplyKeyboard_EN(t *testing.T) {
	kb := replyKeyboard(langEN)

	if kb.Keyboard[0][0].Text != "📝 Add" {
		t.Errorf("expected first button '📝 Add', got %s", kb.Keyboard[0][0].Text)
	}
}
