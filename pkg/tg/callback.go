package tg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/DmiTryAgain/sports-statistics/pkg/db"
)

// CallbackType — тип callback-действия
type CallbackType int

const (
	cbSelectExercise CallbackType = iota
	cbSelectWeight
	cbSelectCount
	cbSelectDistance
	cbSelectDuration
	cbCustomInput
	cbQuickAdd
	cbShowExercise
	cbShowAll
	cbShowPeriod
	cbExercisePage
	cbCancel
	cbSkipParam
)

var errInvalidCallbackData = errors.New("invalid callback data")

// CallbackAction — результат парсинга callback_data
type CallbackAction struct {
	Type         CallbackType
	Exercise     Exercise
	Value        float64
	Count        float64
	Params       *db.StatisticParams
	Period       textPeriod
	Page         int
	Context      string // "add" или "show"
	CustomTarget CustomValueTarget
}

func parseCallbackData(data string) (CallbackAction, error) {
	parts := strings.Split(data, "|")
	if len(parts) == 0 || parts[0] == "" {
		return CallbackAction{}, errInvalidCallbackData
	}

	switch parts[0] {
	case "ex":
		return parseExerciseCallback(parts)
	case "w":
		return parseFloatCallback(parts, cbSelectWeight)
	case "c":
		return parseFloatCallback(parts, cbSelectCount)
	case "d":
		return parseFloatCallback(parts, cbSelectDistance)
	case "t":
		return parseFloatCallback(parts, cbSelectDuration)
	case "cu":
		return parseCustomInputCallback(parts)
	case "qa":
		return parseQuickAddCallback(parts)
	case "se":
		return parseShowExerciseCallback(parts)
	case "sa":
		return CallbackAction{Type: cbShowAll}, nil
	case "sp":
		return parseShowPeriodCallback(parts)
	case "pg":
		return parsePageCallback(parts)
	case "sk":
		return parseSkipParamCallback(parts)
	case "x":
		return CallbackAction{Type: cbCancel}, nil
	default:
		return CallbackAction{}, fmt.Errorf("%w: unknown type %q", errInvalidCallbackData, parts[0])
	}
}

func encodeCallbackData(parts ...string) string {
	return strings.Join(parts, "|")
}

func parseExerciseCallback(parts []string) (CallbackAction, error) {
	if len(parts) < 2 {
		return CallbackAction{}, fmt.Errorf("%w: missing exercise", errInvalidCallbackData)
	}
	return CallbackAction{Type: cbSelectExercise, Exercise: Exercise(parts[1])}, nil
}

func parseFloatCallback(parts []string, cbType CallbackType) (CallbackAction, error) {
	if len(parts) < 2 {
		return CallbackAction{}, fmt.Errorf("%w: missing value", errInvalidCallbackData)
	}
	v, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return CallbackAction{}, fmt.Errorf("%w: %w", errInvalidCallbackData, err)
	}
	return CallbackAction{Type: cbType, Value: v}, nil
}

func parseCustomInputCallback(parts []string) (CallbackAction, error) {
	if len(parts) < 2 {
		return CallbackAction{}, fmt.Errorf("%w: missing target", errInvalidCallbackData)
	}
	var target CustomValueTarget
	switch parts[1] {
	case "w":
		target = TargetWeight
	case "c":
		target = TargetCount
	case "d":
		target = TargetDistance
	case "t":
		target = TargetDuration
	default:
		return CallbackAction{}, fmt.Errorf("%w: unknown target %q", errInvalidCallbackData, parts[1])
	}
	return CallbackAction{Type: cbCustomInput, CustomTarget: target}, nil
}

func parseSkipParamCallback(parts []string) (CallbackAction, error) {
	if len(parts) < 2 {
		return CallbackAction{}, fmt.Errorf("%w: missing param type", errInvalidCallbackData)
	}
	var target CustomValueTarget
	switch parts[1] {
	case "w":
		target = TargetWeight
	case "c":
		target = TargetCount
	case "d":
		target = TargetDistance
	case "t":
		target = TargetDuration
	default:
		return CallbackAction{}, fmt.Errorf("%w: unknown skip target %q", errInvalidCallbackData, parts[1])
	}
	return CallbackAction{Type: cbSkipParam, CustomTarget: target}, nil
}

func parseQuickAddCallback(parts []string) (CallbackAction, error) {
	// qa|EXERCISE|COUNT|WEIGHT|DISTANCE|DURATION
	if len(parts) < 6 {
		return CallbackAction{}, fmt.Errorf("%w: quick add needs 6 parts", errInvalidCallbackData)
	}

	ex := Exercise(parts[1])

	cnt, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return CallbackAction{}, fmt.Errorf("%w: parse count: %w", errInvalidCallbackData, err)
	}

	var params *db.StatisticParams
	w := parseOptionalFloat(parts[3])
	d := parseOptionalFloat(parts[4])
	t := parseOptionalFloat(parts[5])

	if w != nil || d != nil || t != nil {
		params = &db.StatisticParams{
			WeightKg:    w,
			DistanceM:   d,
			DurationSec: t,
		}
	}

	return CallbackAction{
		Type:     cbQuickAdd,
		Exercise: ex,
		Count:    cnt,
		Params:   params,
	}, nil
}

func parseShowExerciseCallback(parts []string) (CallbackAction, error) {
	if len(parts) < 2 {
		return CallbackAction{}, fmt.Errorf("%w: missing exercise", errInvalidCallbackData)
	}
	return CallbackAction{Type: cbShowExercise, Exercise: Exercise(parts[1])}, nil
}

func parseShowPeriodCallback(parts []string) (CallbackAction, error) {
	if len(parts) < 2 {
		return CallbackAction{}, fmt.Errorf("%w: missing period", errInvalidCallbackData)
	}
	return CallbackAction{Type: cbShowPeriod, Period: textPeriod(parts[1])}, nil
}

func parsePageCallback(parts []string) (CallbackAction, error) {
	if len(parts) < 2 {
		return CallbackAction{}, fmt.Errorf("%w: missing page", errInvalidCallbackData)
	}
	page, err := strconv.Atoi(parts[1])
	if err != nil {
		return CallbackAction{}, fmt.Errorf("%w: parse page: %w", errInvalidCallbackData, err)
	}
	cbCtx := ""
	if len(parts) > 2 {
		cbCtx = parts[2]
	}
	return CallbackAction{Type: cbExercisePage, Page: page, Context: cbCtx}, nil
}

func parseOptionalFloat(s string) *float64 {
	if s == "" {
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &v
}

func formatOptionalFloat(v *float64) string {
	if v == nil {
		return ""
	}
	return strconv.FormatFloat(*v, 'f', -1, 64)
}

func encodeQuickAddData(ex Exercise, cnt float64, params *db.StatisticParams) string {
	var w, d, t *float64
	if params != nil {
		w = params.WeightKg
		d = params.DistanceM
		t = params.DurationSec
	}
	return encodeCallbackData(
		"qa",
		string(ex),
		strconv.FormatFloat(cnt, 'f', -1, 64),
		formatOptionalFloat(w),
		formatOptionalFloat(d),
		formatOptionalFloat(t),
	)
}
