package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
)

type GroupedStatistic struct {
	TgUserID       string   `pg:"tgUserId,use_zero"`
	Exercise       string   `pg:"exercise,use_zero"`
	SumCount       float64  `pg:"sumCount,use_zero"`
	Sets           int      `pg:"sets,use_zero"`
	WeightKg       *float64 `pg:"weightKg"`
	DistanceM      *float64 `pg:"distanceM"`
	SumDurationSec *float64 `pg:"sumDurationSec"`
}

type Period struct {
	From, To time.Time
}

type GroupedStatisticSearch struct {
	StatisticSearch
	Periods []Period
}

func (sr StatisticRepo) GroupedStatisticByFilters(ctx context.Context, search GroupedStatisticSearch) (gs []GroupedStatistic, err error) {
	b := strings.Builder{}
	b.WriteString(`
		SELECT t."tgUserId", t."exercise",
			sum(t.count) as "sumCount",
			count(*) as "sets",
			(t."params"->>'weightKg')::float8 as "weightKg",
			(t."params"->>'distanceM')::float8 as "distanceM",
			sum((t."params"->>'durationSec')::float8) as "sumDurationSec"
		FROM statistics t
		WHERE true
	`)

	enabledFilter := string(formatter.FormatQuery([]byte{}, ` AND t."statusId" = ? `, StatusEnabled))
	b.WriteString(enabledFilter)
	tgUserIDFilter := string(formatter.FormatQuery([]byte{}, ` AND t."tgUserId" = ? `, search.TgUserID))
	b.WriteString(tgUserIDFilter)

	if len(search.Exercises) != 0 {
		exFilter := string(formatter.FormatQuery([]byte{}, ` AND t."exercise" in (?) `, pg.In(search.Exercises)))
		b.WriteString(exFilter)
	}

	if len(search.Periods) != 0 {
		b.WriteString(`AND (false`)
		for i := range search.Periods {
			periodFilter := string(formatter.FormatQuery([]byte{}, ` OR (t."createdAt" >= ? AND t."createdAt" < ?) `, search.Periods[i].From, search.Periods[i].To))
			b.WriteString(periodFilter)
		}
		b.WriteString(`)`)
	}

	b.WriteString(`
		GROUP BY t."tgUserId", t."exercise", t."params"->>'weightKg', t."params"->>'distanceM'
		ORDER BY t."exercise", "weightKg" DESC NULLS LAST, "distanceM" DESC NULLS LAST, "sumCount" DESC
	`)

	if _, err = sr.db.QueryContext(ctx, &gs, b.String()); err != nil {
		return nil, fmt.Errorf("grouped statistic by filter=%+v, err=%w", search, err)
	}

	return gs, nil
}

// FrequentExercise описывает часто используемую комбинацию упражнения с параметрами
type FrequentExercise struct {
	Exercise  string           `pg:"exercise,use_zero"`
	Count     float64          `pg:"count,use_zero"`
	Params    *StatisticParams `pg:"params"`
	Frequency int              `pg:"frequency,use_zero"`
}

// FrequentExercisesByUser возвращает топ частых комбинаций (упражнение + count + params) за последние 30 дней.
func (sr StatisticRepo) FrequentExercisesByUser(ctx context.Context, tgUserID string, limit int) ([]FrequentExercise, error) {
	var result []FrequentExercise
	_, err := sr.db.QueryContext(ctx, &result, `
		SELECT t."exercise", t."count", t."params",
			   COUNT(*) as "frequency"
		FROM statistics t
		WHERE t."tgUserId" = ?
		  AND t."statusId" = ?
		  AND t."createdAt" >= NOW() - INTERVAL '30 days'
		GROUP BY t."exercise", t."count", t."params"
		ORDER BY "frequency" DESC, MAX(t."createdAt") DESC
		LIMIT ?
	`, tgUserID, StatusEnabled, limit)
	if err != nil {
		return nil, fmt.Errorf("frequent exercises for user=%s, err=%w", tgUserID, err)
	}
	return result, nil
}

// UniqueExercisesByUser возвращает все уникальные упражнения пользователя.
func (sr StatisticRepo) UniqueExercisesByUser(ctx context.Context, tgUserID string) ([]string, error) {
	var result []struct {
		Exercise string `pg:"exercise"`
	}
	_, err := sr.db.QueryContext(ctx, &result, `
		SELECT DISTINCT t."exercise"
		FROM statistics t
		WHERE t."tgUserId" = ?
		  AND t."statusId" = ?
		ORDER BY t."exercise"
	`, tgUserID, StatusEnabled)
	if err != nil {
		return nil, fmt.Errorf("unique exercises for user=%s, err=%w", tgUserID, err)
	}

	exercises := make([]string, len(result))
	for i := range result {
		exercises[i] = result[i].Exercise
	}
	return exercises, nil
}

// UniqueWeightsByExercise возвращает уникальные веса для упражнения пользователя.
func (sr StatisticRepo) UniqueWeightsByExercise(ctx context.Context, tgUserID, exercise string, limit int) ([]float64, error) {
	var result []struct {
		Weight float64 `pg:"weight"`
	}
	_, err := sr.db.QueryContext(ctx, &result, `
		SELECT DISTINCT (t."params"->>'weightKg')::float8 as "weight"
		FROM statistics t
		WHERE t."tgUserId" = ?
		  AND t."exercise" = ?
		  AND t."statusId" = ?
		  AND t."params"->>'weightKg' IS NOT NULL
		ORDER BY "weight" DESC
		LIMIT ?
	`, tgUserID, exercise, StatusEnabled, limit)
	if err != nil {
		return nil, fmt.Errorf("unique weights for user=%s exercise=%s, err=%w", tgUserID, exercise, err)
	}

	weights := make([]float64, len(result))
	for i := range result {
		weights[i] = result[i].Weight
	}
	return weights, nil
}
