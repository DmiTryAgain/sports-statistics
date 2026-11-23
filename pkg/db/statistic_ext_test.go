package db_test

import (
	"slices"
	"testing"
	"time"

	"github.com/DmiTryAgain/sports-statistics/pkg/db"
	dbtest "github.com/DmiTryAgain/sports-statistics/pkg/db/test"
)

func TestFrequentExercisesByUser(t *testing.T) {
	ctx := t.Context()
	dbc, _ := dbtest.Setup(t)
	userID := dbtest.NextStringID()
	repo := db.NewStatisticRepo(dbc)

	// Создаём 5 записей жим 80кг x10
	for range 5 {
		_, cl := dbtest.Statistic(t, dbc, &db.Statistic{
			TgUserID: userID,
			Exercise: "benchPress",
			Count:    10,
			Params:   &db.StatisticParams{WeightKg: ptr[float64](80)},
			StatusID: 1,
		}, dbtest.WithFakeStatistic)
		defer cl()
	}

	// Создаём 3 записи подтягивания x12
	for range 3 {
		_, cl := dbtest.Statistic(t, dbc, &db.Statistic{
			TgUserID: userID,
			Exercise: "pullUp",
			Count:    12,
			StatusID: 1,
		}, dbtest.WithFakeStatistic)
		defer cl()
	}

	// Создаём 1 запись от другого пользователя (не должна попасть)
	_, cl := dbtest.Statistic(t, dbc, &db.Statistic{
		TgUserID: dbtest.NextStringID(),
		Exercise: "pullUp",
		Count:    5,
		StatusID: 1,
	}, dbtest.WithFakeStatistic)
	defer cl()

	result, err := repo.FrequentExercisesByUser(ctx, userID, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 frequent exercises, got %d", len(result))
	}

	// Первый — жим (частота 5)
	if result[0].Exercise != "benchPress" {
		t.Errorf("expected first exercise benchPress, got %s", result[0].Exercise)
	}
	if result[0].Frequency != 5 {
		t.Errorf("expected frequency 5, got %d", result[0].Frequency)
	}
	if result[0].Count != 10 {
		t.Errorf("expected count 10, got %f", result[0].Count)
	}
	if result[0].Params == nil || result[0].Params.WeightKg == nil || *result[0].Params.WeightKg != 80 {
		t.Errorf("expected weight 80, got %+v", result[0].Params)
	}

	// Второй — подтягивания (частота 3)
	if result[1].Exercise != "pullUp" {
		t.Errorf("expected second exercise pullUp, got %s", result[1].Exercise)
	}
	if result[1].Frequency != 3 {
		t.Errorf("expected frequency 3, got %d", result[1].Frequency)
	}
}

func TestFrequentExercisesByUser_Limit(t *testing.T) {
	ctx := t.Context()
	dbc, _ := dbtest.Setup(t)
	userID := dbtest.NextStringID()
	repo := db.NewStatisticRepo(dbc)

	// Создаём 3 разных упражнения
	for _, ex := range []string{"pullUp", "pushUp", "abs"} {
		_, cl := dbtest.Statistic(t, dbc, &db.Statistic{
			TgUserID: userID,
			Exercise: ex,
			Count:    10,
			StatusID: 1,
		}, dbtest.WithFakeStatistic)
		defer cl()
	}

	result, err := repo.FrequentExercisesByUser(ctx, userID, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 results with limit 2, got %d", len(result))
	}
}

func TestFrequentExercisesByUser_OldRecordsExcluded(t *testing.T) {
	ctx := t.Context()
	dbc, _ := dbtest.Setup(t)
	userID := dbtest.NextStringID()
	repo := db.NewStatisticRepo(dbc)

	// Создаём запись старше 30 дней
	_, cl := dbtest.Statistic(t, dbc, &db.Statistic{
		TgUserID:  userID,
		Exercise:  "pullUp",
		Count:     10,
		StatusID:  1,
		CreatedAt: time.Now().AddDate(0, 0, -31),
	}, dbtest.WithFakeStatistic)
	defer cl()

	result, err := repo.FrequentExercisesByUser(ctx, userID, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Fatalf("expected 0 results for old records, got %d", len(result))
	}
}

func TestUniqueExercisesByUser(t *testing.T) {
	ctx := t.Context()
	dbc, _ := dbtest.Setup(t)
	userID := dbtest.NextStringID()
	repo := db.NewStatisticRepo(dbc)

	exercises := []string{"abs", "benchPress", "pullUp"}
	for _, ex := range exercises {
		// Добавляем по 2 записи каждого упражнения
		for range 2 {
			_, cl := dbtest.Statistic(t, dbc, &db.Statistic{
				TgUserID: userID,
				Exercise: ex,
				Count:    10,
				StatusID: 1,
			}, dbtest.WithFakeStatistic)
			defer cl()
		}
	}

	result, err := repo.UniqueExercisesByUser(ctx, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 unique exercises, got %d: %v", len(result), result)
	}

	// Результат отсортирован по алфавиту
	if result[0] != "abs" {
		t.Errorf("expected first exercise abs, got %s", result[0])
	}
	if result[1] != "benchPress" {
		t.Errorf("expected second exercise benchPress, got %s", result[1])
	}
	if result[2] != "pullUp" {
		t.Errorf("expected third exercise pullUp, got %s", result[2])
	}
}

func TestUniqueExercisesByUser_Empty(t *testing.T) {
	ctx := t.Context()
	dbc, _ := dbtest.Setup(t)
	userID := dbtest.NextStringID()
	repo := db.NewStatisticRepo(dbc)

	result, err := repo.UniqueExercisesByUser(ctx, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Fatalf("expected 0 exercises for new user, got %d", len(result))
	}
}

func TestUniqueWeightsByExercise(t *testing.T) {
	ctx := t.Context()
	dbc, _ := dbtest.Setup(t)
	userID := dbtest.NextStringID()
	repo := db.NewStatisticRepo(dbc)

	weights := []float64{60, 80, 80, 100}
	for _, w := range weights {
		_, cl := dbtest.Statistic(t, dbc, &db.Statistic{
			TgUserID: userID,
			Exercise: "benchPress",
			Count:    10,
			Params:   &db.StatisticParams{WeightKg: ptr(w)},
			StatusID: 1,
		}, dbtest.WithFakeStatistic)
		defer cl()
	}

	// Добавляем запись без веса (не должна попасть)
	_, cl := dbtest.Statistic(t, dbc, &db.Statistic{
		TgUserID: userID,
		Exercise: "benchPress",
		Count:    10,
		StatusID: 1,
	}, dbtest.WithFakeStatistic)
	defer cl()

	result, err := repo.UniqueWeightsByExercise(ctx, userID, "benchPress", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 3 уникальных веса: 100, 80, 60 (DESC)
	if len(result) != 3 {
		t.Fatalf("expected 3 unique weights, got %d: %v", len(result), result)
	}

	if result[0] != 100 {
		t.Errorf("expected first weight 100, got %f", result[0])
	}
	if result[1] != 80 {
		t.Errorf("expected second weight 80, got %f", result[1])
	}
	if result[2] != 60 {
		t.Errorf("expected third weight 60, got %f", result[2])
	}
}

func TestUniqueWeightsByExercise_Limit(t *testing.T) {
	ctx := t.Context()
	dbc, _ := dbtest.Setup(t)
	userID := dbtest.NextStringID()
	repo := db.NewStatisticRepo(dbc)

	for _, w := range []float64{40, 60, 80, 100} {
		_, cl := dbtest.Statistic(t, dbc, &db.Statistic{
			TgUserID: userID,
			Exercise: "deadlift",
			Count:    5,
			Params:   &db.StatisticParams{WeightKg: ptr(w)},
			StatusID: 1,
		}, dbtest.WithFakeStatistic)
		defer cl()
	}

	result, err := repo.UniqueWeightsByExercise(ctx, userID, "deadlift", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 weights with limit 2, got %d", len(result))
	}

	// Должны быть 2 самых тяжёлых: 100, 80
	if result[0] != 100 {
		t.Errorf("expected first weight 100, got %f", result[0])
	}
	if result[1] != 80 {
		t.Errorf("expected second weight 80, got %f", result[1])
	}
}

func TestUniqueWeightsByExercise_DifferentExercise(t *testing.T) {
	ctx := t.Context()
	dbc, _ := dbtest.Setup(t)
	userID := dbtest.NextStringID()
	repo := db.NewStatisticRepo(dbc)

	// Добавляем вес для benchPress
	_, cl := dbtest.Statistic(t, dbc, &db.Statistic{
		TgUserID: userID,
		Exercise: "benchPress",
		Count:    10,
		Params:   &db.StatisticParams{WeightKg: ptr[float64](80)},
		StatusID: 1,
	}, dbtest.WithFakeStatistic)
	defer cl()

	// Ищем для deadlift — должно быть пусто
	result, err := repo.UniqueWeightsByExercise(ctx, userID, "deadlift", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Fatalf("expected 0 weights for different exercise, got %d", len(result))
	}
}

func TestStatisticsByPeriods(t *testing.T) {
	ctx := t.Context()
	dbc, _ := dbtest.Setup(t)
	userID := dbtest.NextStringID()
	repo := db.NewStatisticRepo(dbc)

	// Очистка через defer, а не t.Cleanup: фабрика удаляет запись с t.Context(), который к моменту t.Cleanup уже отменён
	var cleaners []dbtest.Cleaner
	defer func() {
		for _, cl := range cleaners {
			cl()
		}
	}()

	newest := time.Date(2025, 10, 15, 12, 0, 0, 0, time.UTC)
	add := func(in db.Statistic) int {
		in.Count = 10
		st, cl := dbtest.Statistic(t, dbc, &in, dbtest.WithFakeStatistic)
		cleaners = append(cleaners, cl)
		return st.ID
	}

	pullUpNewest := add(db.Statistic{TgUserID: userID, Exercise: "pullUp", CreatedAt: newest})
	benchPress := add(db.Statistic{TgUserID: userID, Exercise: "benchPress", CreatedAt: newest.AddDate(0, 0, -1), Params: &db.StatisticParams{WeightKg: ptr[float64](80)}})
	pullUpOldest := add(db.Statistic{TgUserID: userID, Exercise: "pullUp", CreatedAt: newest.AddDate(0, 0, -10)})
	// Отключённая запись и запись другого пользователя не должны попадать в выборку
	add(db.Statistic{TgUserID: userID, Exercise: "pullUp", CreatedAt: newest, StatusID: db.StatusDisabled})
	add(db.Statistic{TgUserID: dbtest.NextStringID(), Exercise: "pullUp", CreatedAt: newest})

	tests := []struct {
		name      string
		exercises []string
		periods   []db.Period
		limit     int
		wantIDs   []int
		wantTotal int
	}{
		{
			name:      "all time newest first",
			limit:     10,
			wantIDs:   []int{pullUpNewest, benchPress, pullUpOldest},
			wantTotal: 3,
		},
		{
			name:      "filter by exercise",
			exercises: []string{"pullUp"},
			limit:     10,
			wantIDs:   []int{pullUpNewest, pullUpOldest},
			wantTotal: 2,
		},
		{
			name:      "period includes from and excludes to",
			periods:   []db.Period{{From: newest.AddDate(0, 0, -1), To: newest}},
			limit:     10,
			wantIDs:   []int{benchPress},
			wantTotal: 1,
		},
		{
			name: "several periods",
			periods: []db.Period{
				{From: newest.AddDate(0, 0, -10), To: newest.AddDate(0, 0, -9)},
				{From: newest, To: newest.Add(time.Hour)},
			},
			limit:     10,
			wantIDs:   []int{pullUpNewest, pullUpOldest},
			wantTotal: 2,
		},
		{
			name:      "limit keeps newest and total counts all",
			limit:     2,
			wantIDs:   []int{pullUpNewest, benchPress},
			wantTotal: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			search := &db.StatisticSearch{TgUserID: &userID, Exercises: tt.exercises}
			got, total, err := repo.StatisticsByPeriods(ctx, search, tt.periods, tt.limit)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			gotIDs := make([]int, len(got))
			for i := range got {
				gotIDs[i] = got[i].ID
			}
			if !slices.Equal(gotIDs, tt.wantIDs) {
				t.Errorf("got IDs %v, want %v", gotIDs, tt.wantIDs)
			}
			if total != tt.wantTotal {
				t.Errorf("got total %d, want %d", total, tt.wantTotal)
			}
		})
	}
}

func ptr[T any](t T) *T { return &t }
