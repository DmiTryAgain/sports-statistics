package tg

import (
	"errors"
	"testing"

	"github.com/DmiTryAgain/sports-statistics/pkg/db"
)

func TestParseCallbackData(t *testing.T) { //nolint:gocognit
	tests := []struct {
		name    string
		data    string
		want    CallbackAction
		wantErr bool
	}{
		{
			name: "select exercise",
			data: "ex|pullUp",
			want: CallbackAction{Type: cbSelectExercise, Exercise: pullUpEx},
		},
		{
			name: "select weight",
			data: "w|80",
			want: CallbackAction{Type: cbSelectWeight, Value: 80},
		},
		{
			name: "select count",
			data: "c|10",
			want: CallbackAction{Type: cbSelectCount, Value: 10},
		},
		{
			name: "select distance",
			data: "d|5000",
			want: CallbackAction{Type: cbSelectDistance, Value: 5000},
		},
		{
			name: "select duration",
			data: "t|1500",
			want: CallbackAction{Type: cbSelectDuration, Value: 1500},
		},
		{
			name: "custom input weight",
			data: "cu|w",
			want: CallbackAction{Type: cbCustomInput, CustomTarget: TargetWeight},
		},
		{
			name: "custom input duration",
			data: "cu|t",
			want: CallbackAction{Type: cbCustomInput, CustomTarget: TargetDuration},
		},
		{
			name: "quick add with weight",
			data: "qa|benchPress|10|80||",
			want: CallbackAction{
				Type:     cbQuickAdd,
				Exercise: benchPressEx,
				Count:    10,
				Params:   &db.StatisticParams{WeightKg: ptr[float64](80)},
			},
		},
		{
			name: "quick add reps only",
			data: "qa|pullUp|12|||",
			want: CallbackAction{
				Type:     cbQuickAdd,
				Exercise: pullUpEx,
				Count:    12,
				Params:   nil,
			},
		},
		{
			name: "quick add dist+time",
			data: "qa|jogging|1|0|5000|1500",
			want: CallbackAction{
				Type:     cbQuickAdd,
				Exercise: joggingEx,
				Count:    1,
				Params:   &db.StatisticParams{WeightKg: ptr[float64](0), DistanceM: ptr[float64](5000), DurationSec: ptr[float64](1500)},
			},
		},
		{
			name: "show exercise",
			data: "se|pullUp",
			want: CallbackAction{Type: cbShowExercise, Exercise: pullUpEx},
		},
		{
			name: "show all",
			data: "sa",
			want: CallbackAction{Type: cbShowAll},
		},
		{
			name: "show period",
			data: "sp|today",
			want: CallbackAction{Type: cbShowPeriod, Period: todayPeriod},
		},
		{
			name: "page",
			data: "pg|2|add",
			want: CallbackAction{Type: cbExercisePage, Page: 2, Context: "add"},
		},
		{
			name: "cancel",
			data: "x",
			want: CallbackAction{Type: cbCancel},
		},
		{
			name:    "empty data",
			data:    "",
			wantErr: true,
		},
		{
			name:    "unknown type",
			data:    "zzz|123",
			wantErr: true,
		},
		{
			name:    "weight missing value",
			data:    "w",
			wantErr: true,
		},
		{
			name:    "weight invalid value",
			data:    "w|abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCallbackData(tt.data)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !errors.Is(err, errInvalidCallbackData) {
					t.Fatalf("expected errInvalidCallbackData, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Type != tt.want.Type {
				t.Errorf("Type = %d, want %d", got.Type, tt.want.Type)
			}
			if got.Exercise != tt.want.Exercise {
				t.Errorf("Exercise = %q, want %q", got.Exercise, tt.want.Exercise)
			}
			if got.Value != tt.want.Value {
				t.Errorf("Value = %f, want %f", got.Value, tt.want.Value)
			}
			if got.Count != tt.want.Count {
				t.Errorf("Count = %f, want %f", got.Count, tt.want.Count)
			}
			if got.Period != tt.want.Period {
				t.Errorf("Period = %q, want %q", got.Period, tt.want.Period)
			}
			if got.Page != tt.want.Page {
				t.Errorf("Page = %d, want %d", got.Page, tt.want.Page)
			}
			if got.Context != tt.want.Context {
				t.Errorf("Context = %q, want %q", got.Context, tt.want.Context)
			}
			if got.CustomTarget != tt.want.CustomTarget {
				t.Errorf("CustomTarget = %d, want %d", got.CustomTarget, tt.want.CustomTarget)
			}
			// Compare params
			if !paramsEqual(got.Params, tt.want.Params) {
				t.Errorf("Params = %+v, want %+v", got.Params, tt.want.Params)
			}
		})
	}
}

func TestEncodeCallbackData(t *testing.T) {
	tests := []struct {
		name  string
		parts []string
		want  string
	}{
		{
			name:  "exercise",
			parts: []string{"ex", "pullUp"},
			want:  "ex|pullUp",
		},
		{
			name:  "weight",
			parts: []string{"w", "80"},
			want:  "w|80",
		},
		{
			name:  "cancel",
			parts: []string{"x"},
			want:  "x",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodeCallbackData(tt.parts...)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEncodeQuickAddData(t *testing.T) {
	tests := []struct {
		name   string
		ex     Exercise
		cnt    float64
		params *db.StatisticParams
		want   string
	}{
		{
			name:   "reps only",
			ex:     pullUpEx,
			cnt:    12,
			params: nil,
			want:   "qa|pullUp|12|||",
		},
		{
			name:   "with weight",
			ex:     benchPressEx,
			cnt:    10,
			params: &db.StatisticParams{WeightKg: ptr[float64](80)},
			want:   "qa|benchPress|10|80||",
		},
		{
			name:   "dist and time",
			ex:     joggingEx,
			cnt:    1,
			params: &db.StatisticParams{DistanceM: ptr[float64](5000), DurationSec: ptr[float64](1500)},
			want:   "qa|jogging|1||5000|1500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodeQuickAddData(tt.ex, tt.cnt, tt.params)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCallbackDataLength(t *testing.T) {
	// Verify all possible exercise Quick-Add callbacks fit in 64 bytes
	for _, ex := range exerciseOrder {
		data := encodeQuickAddData(ex, 999, &db.StatisticParams{
			WeightKg:    ptr[float64](999),
			DistanceM:   ptr[float64](99999),
			DurationSec: ptr[float64](99999),
		})
		if len(data) > 64 {
			t.Errorf("callback data for %s is %d bytes (>64): %q", ex, len(data), data)
		}
	}
}

func paramsEqual(a, b *db.StatisticParams) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return floatPtrEqual(a.WeightKg, b.WeightKg) &&
		floatPtrEqual(a.DistanceM, b.DistanceM) &&
		floatPtrEqual(a.DurationSec, b.DurationSec)
}
