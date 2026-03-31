package tg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/DmiTryAgain/sports-statistics/pkg/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const buttonsPerRow = 3

func replyKeyboard(lang language) tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(messagesByLang[lang][addBtn]),
			tgbotapi.NewKeyboardButton(messagesByLang[lang][showBtn]),
			tgbotapi.NewKeyboardButton(messagesByLang[lang][helpBtn]),
		),
	)
	kb.ResizeKeyboard = true
	return kb
}

func exerciseInlineKeyboard(page int, lang language, cbCtx string) tgbotapi.InlineKeyboardMarkup {
	start := page * exercisesPerPage
	if start >= len(exerciseOrder) {
		start = 0
		page = 0
	}
	end := start + exercisesPerPage
	if end > len(exerciseOrder) {
		end = len(exerciseOrder)
	}

	exSlice := exerciseOrder[start:end]

	var rows [][]tgbotapi.InlineKeyboardButton
	row := make([]tgbotapi.InlineKeyboardButton, 0, len(exSlice))
	for i, ex := range exSlice {
		text := exTextByLang[lang][ex]
		data := encodeCallbackData("ex", string(ex))
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(text, data))
		if (i+1)%buttonsPerRow == 0 || i == len(exSlice)-1 {
			rows = append(rows, row)
			row = nil
		}
	}

	// Navigation
	var navRow []tgbotapi.InlineKeyboardButton
	if page > 0 {
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData(
			messagesByLang[lang][backBtn],
			encodeCallbackData("pg", strconv.Itoa(page-1), cbCtx),
		))
	}
	if end < len(exerciseOrder) {
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData(
			messagesByLang[lang][moreBtn],
			encodeCallbackData("pg", strconv.Itoa(page+1), cbCtx),
		))
	}
	if len(navRow) > 0 {
		rows = append(rows, navRow)
	}

	// Cancel
	rows = append(rows, cancelRow(lang))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func showExerciseInlineKeyboard(exercises []string, lang language) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	row := make([]tgbotapi.InlineKeyboardButton, 0, len(exercises))

	for i, exStr := range exercises {
		ex := Exercise(exStr)
		text := exTextByLang[lang][ex]
		if text == "" {
			text = exStr
		}
		data := encodeCallbackData("se", exStr)
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(text, data))
		if (i+1)%buttonsPerRow == 0 || i == len(exercises)-1 {
			rows = append(rows, row)
			row = nil
		}
	}

	// "All" button
	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(messagesByLang[lang][allExBtn], "sa"),
	})

	rows = append(rows, cancelRow(lang))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func weightInlineKeyboard(userWeights []float64, lang language) tgbotapi.InlineKeyboardMarkup {
	weights := userWeights
	if len(weights) == 0 {
		weights = []float64{20, 40, 60, 80, 100}
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	row := make([]tgbotapi.InlineKeyboardButton, 0, len(weights))
	for i, w := range weights {
		text := formatWeight(w, lang)
		data := encodeCallbackData("w", strconv.FormatFloat(w, 'f', -1, 64))
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(text, data))
		if (i+1)%buttonsPerRow == 0 || i == len(weights)-1 {
			rows = append(rows, row)
			row = nil
		}
	}

	rows = append(rows, customAndCancelRow(lang, "w"))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func countInlineKeyboard(lang language) tgbotapi.InlineKeyboardMarkup {
	counts := []int{5, 8, 10, 12, 15, 20}

	var rows [][]tgbotapi.InlineKeyboardButton
	row := make([]tgbotapi.InlineKeyboardButton, 0, len(counts))
	for i, c := range counts {
		text := strconv.Itoa(c)
		data := encodeCallbackData("c", text)
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(text, data))
		if (i+1)%buttonsPerRow == 0 || i == len(counts)-1 {
			rows = append(rows, row)
			row = nil
		}
	}

	rows = append(rows, customAndCancelRow(lang, "c"))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func distanceInlineKeyboard(lang language) tgbotapi.InlineKeyboardMarkup {
	distances := []float64{1000, 1500, 2000, 2500, 3000, 5000}

	var rows [][]tgbotapi.InlineKeyboardButton
	row := make([]tgbotapi.InlineKeyboardButton, 0, len(distances))
	for i, d := range distances {
		text := formatDistance(d, lang)
		data := encodeCallbackData("d", strconv.FormatFloat(d, 'f', -1, 64))
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(text, data))
		if (i+1)%buttonsPerRow == 0 || i == len(distances)-1 {
			rows = append(rows, row)
			row = nil
		}
	}

	rows = append(rows, customAndCancelRow(lang, "d"))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func durationInlineKeyboard(category ExerciseCategory, lang language) tgbotapi.InlineKeyboardMarkup {
	durations := []float64{900, 1200, 1500, 1800, 2700, 3600}
	if category == CategoryDuration || category == CategoryRepsOrDuration {
		durations = []float64{30, 45, 60, 90, 120, 180}
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	row := make([]tgbotapi.InlineKeyboardButton, 0, len(durations))
	for i, d := range durations {
		text := formatDuration(d, lang)
		data := encodeCallbackData("t", strconv.FormatFloat(d, 'f', -1, 64))
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(text, data))
		if (i+1)%buttonsPerRow == 0 || i == len(durations)-1 {
			rows = append(rows, row)
			row = nil
		}
	}

	rows = append(rows, customAndCancelRow(lang, "t"))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func periodInlineKeyboard(lang language) tgbotapi.InlineKeyboardMarkup {
	periodsOrder := []textPeriod{todayPeriod, yesterdayPeriod, weekPeriod, lastWeekPeriod, monthPeriod, lastMonthPeriod, yearPeriod, lastYearPeriod, allPeriod}

	var rows [][]tgbotapi.InlineKeyboardButton
	row := make([]tgbotapi.InlineKeyboardButton, 0, len(periodsOrder))
	for i, p := range periodsOrder {
		text := periodTextByLang[lang][p]
		data := encodeCallbackData("sp", string(p))
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(text, data))
		if (i+1)%buttonsPerRow == 0 || i == len(periodsOrder)-1 {
			rows = append(rows, row)
			row = nil
		}
	}

	rows = append(rows, cancelRow(lang))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func quickAddInlineKeyboard(frequent []db.FrequentExercise, lang language) [][]tgbotapi.InlineKeyboardButton {
	if len(frequent) == 0 {
		return nil
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(frequent))
	for _, f := range frequent {
		text := formatQuickAddButton(f, lang)
		data := encodeQuickAddData(Exercise(f.Exercise), f.Count, f.Params)
		rows = append(rows, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(text, data),
		})
	}
	return rows
}

func formatQuickAddButton(f db.FrequentExercise, lang language) string {
	ex := Exercise(f.Exercise)
	name := exTextByLang[lang][ex]
	if name == "" {
		name = f.Exercise
	}

	var weightKg, distanceM, durationSec *float64
	if f.Params != nil {
		weightKg, distanceM, durationSec = f.Params.WeightKg, f.Params.DistanceM, f.Params.DurationSec
	}
	showCount := f.Count > 1 || (durationSec == nil && distanceM == nil)

	parts := name
	if weightKg != nil {
		parts += " " + formatWeight(*weightKg, lang)
	}
	if showCount {
		parts += fmt.Sprintf(" ×%g", f.Count)
	}
	if distanceM != nil {
		parts += " " + formatDistance(*distanceM, lang)
	}
	if durationSec != nil {
		parts += " " + formatDuration(*durationSec, lang)
	}
	return parts
}

func confirmDetail(weightKg, distanceM, durationSec *float64, cnt float64, lang language) string {
	if weightKg != nil && durationSec == nil && distanceM == nil {
		return fmt.Sprintf("%s × %g", formatWeight(*weightKg, lang), cnt)
	}
	showCount := cnt > 1 || (durationSec == nil && distanceM == nil)
	var parts []string
	if weightKg != nil {
		parts = append(parts, formatWeight(*weightKg, lang))
	}
	if showCount {
		parts = append(parts, fmt.Sprintf("×%g", cnt))
	}
	if distanceM != nil {
		parts = append(parts, formatDistance(*distanceM, lang))
	}
	if durationSec != nil {
		parts = append(parts, formatDuration(*durationSec, lang))
	}
	return strings.Join(parts, " ")
}

func formatAddConfirmation(ex Exercise, cnt float64, params *db.StatisticParams, lang language) string {
	name := exTextByLang[lang][ex]
	if name == "" {
		name = ex.String()
	}
	var weightKg, distanceM, durationSec *float64
	if params != nil {
		weightKg, distanceM, durationSec = params.WeightKg, params.DistanceM, params.DurationSec
	}
	detail := confirmDetail(weightKg, distanceM, durationSec, cnt, lang)
	confirm := fmt.Sprintf(messagesByLang[lang][addedConfirmation], name, detail)
	quickCmd := formatQuickCopyCommand(ex, cnt, params, lang)
	return confirm + "\n" + fmt.Sprintf(messagesByLang[lang][quickCopyHint], "`"+quickCmd+"`")
}

// formatQuickCopyCommand формирует текстовую команду для быстрого копирования.
// Результат можно вставить в чат и отправить — парсер распознает его корректно.
func formatQuickCopyCommand(ex Exercise, cnt float64, params *db.StatisticParams, lang language) string {
	cmdWord := cmdTextByLang[lang][addCmd]
	exName := exTextByLang[lang][ex]
	if exName == "" {
		exName = ex.String()
	}

	parts := []string{cmdWord, exName}

	if params != nil && params.WeightKg != nil {
		parts = append(parts, formatWeight(*params.WeightKg, lang))
	}

	if cnt > 1 || params == nil || (params.DurationSec == nil && params.DistanceM == nil) {
		parts = append(parts, strconv.FormatFloat(cnt, 'f', -1, 64))
	}

	if params != nil && params.DistanceM != nil {
		parts = append(parts, formatDistance(*params.DistanceM, lang))
	}
	if params != nil && params.DurationSec != nil {
		parts = append(parts, formatDuration(*params.DurationSec, lang))
	}

	return strings.Join(parts, " ")
}

func cancelRow(lang language) []tgbotapi.InlineKeyboardButton {
	return []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(messagesByLang[lang][cancelBtn], "x"),
	}
}

func customAndCancelRow(lang language, target string) []tgbotapi.InlineKeyboardButton {
	return []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(messagesByLang[lang][customInputBtn], encodeCallbackData("cu", target)),
		tgbotapi.NewInlineKeyboardButtonData(messagesByLang[lang][cancelBtn], "x"),
	}
}

func optionalDistanceInlineKeyboard(lang language) tgbotapi.InlineKeyboardMarkup {
	distances := []float64{1000, 1500, 2000, 2500, 3000, 5000}

	var rows [][]tgbotapi.InlineKeyboardButton
	row := make([]tgbotapi.InlineKeyboardButton, 0, len(distances))
	for i, d := range distances {
		text := formatDistance(d, lang)
		data := encodeCallbackData("d", strconv.FormatFloat(d, 'f', -1, 64))
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(text, data))
		if (i+1)%buttonsPerRow == 0 || i == len(distances)-1 {
			rows = append(rows, row)
			row = nil
		}
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(messagesByLang[lang][customInputBtn], encodeCallbackData("cu", "d")),
		tgbotapi.NewInlineKeyboardButtonData(messagesByLang[lang][skipBtn], encodeCallbackData("sk", "d")),
	})
	rows = append(rows, cancelRow(lang))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func optionalDurationInlineKeyboard(category ExerciseCategory, lang language) tgbotapi.InlineKeyboardMarkup {
	durations := []float64{900, 1200, 1500, 1800, 2700, 3600}
	if category == CategoryDuration || category == CategoryDurationWeight || category == CategoryRepsOrDuration {
		durations = []float64{30, 45, 60, 90, 120, 180}
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	row := make([]tgbotapi.InlineKeyboardButton, 0, len(durations))
	for i, d := range durations {
		text := formatDuration(d, lang)
		data := encodeCallbackData("t", strconv.FormatFloat(d, 'f', -1, 64))
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(text, data))
		if (i+1)%buttonsPerRow == 0 || i == len(durations)-1 {
			rows = append(rows, row)
			row = nil
		}
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(messagesByLang[lang][customInputBtn], encodeCallbackData("cu", "t")),
		tgbotapi.NewInlineKeyboardButtonData(messagesByLang[lang][skipBtn], encodeCallbackData("sk", "t")),
	})
	rows = append(rows, cancelRow(lang))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func optionalWeightInlineKeyboard(userWeights []float64, lang language) tgbotapi.InlineKeyboardMarkup {
	weights := userWeights
	if len(weights) == 0 {
		weights = []float64{20, 40, 60, 80, 100}
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	row := make([]tgbotapi.InlineKeyboardButton, 0, len(weights))
	for i, w := range weights {
		text := formatWeight(w, lang)
		data := encodeCallbackData("w", strconv.FormatFloat(w, 'f', -1, 64))
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(text, data))
		if (i+1)%buttonsPerRow == 0 || i == len(weights)-1 {
			rows = append(rows, row)
			row = nil
		}
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(messagesByLang[lang][customInputBtn], encodeCallbackData("cu", "w")),
		tgbotapi.NewInlineKeyboardButtonData(messagesByLang[lang][skipBtn], encodeCallbackData("sk", "w")),
	})
	rows = append(rows, cancelRow(lang))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
