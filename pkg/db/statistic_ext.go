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
