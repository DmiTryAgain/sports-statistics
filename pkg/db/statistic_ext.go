package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
)

type GroupedStatistic struct {
	TgUserID string  `pg:"tgUserId,use_zero"`
	Exercise string  `pg:"exercise,use_zero"`
	SumCount float64 `pg:"sumCount,use_zero"`
	Sets     int     `pg:"sets,use_zero"`
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
		SELECT t."tgUserId", t."exercise", sum(t.count) as "sumCount", count(t.count) as "sets"
		FROM statistics t
		WHERE true
	`)

	enabledFilter := string(formatter.FormatQuery([]byte{}, ` AND t."statusId" = ? `, StatusEnabled))
	b.WriteString(enabledFilter)
	tgUserIDFilter := string(formatter.FormatQuery([]byte{}, ` AND t."tgUserId" = ? `, search.TgUserID))
	b.WriteString(tgUserIDFilter)

	if len(search.Exercises) != 0 {
		exFilter := string(formatter.FormatQuery([]byte{}, ` AND t."exercise" = ? `, pg.In(search.Exercises)))
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
		GROUP BY 1, 2
		ORDER BY 3 DESC
	`)

	if _, err = sr.db.QueryContext(ctx, &gs, b.String()); err != nil {
		return nil, fmt.Errorf("grouped statistic by filter=%+v, err=%w", search, err)
	}

	return gs, nil
}
