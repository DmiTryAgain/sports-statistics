package db

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/vmkteam/embedlog"
)

type QueryLogger struct {
	embedlog.Logger
}

func NewQueryLogger(logger embedlog.Logger) QueryLogger {
	return QueryLogger{Logger: logger}
}

func (ql QueryLogger) BeforeQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	if event.Stash == nil {
		event.Stash = make(map[interface{}]interface{})
	}

	event.Stash["startedAt"] = time.Now()
	return ctx, nil
}

func (ql QueryLogger) AfterQuery(ctx context.Context, event *pg.QueryEvent) error {
	query, err := event.FormattedQuery()
	if err != nil {
		ql.Error(ctx, string(query), "err", err)
	}

	var since time.Duration
	if event.Stash != nil {
		if v, ok := event.Stash["startedAt"]; ok {
			if startAt, ok := v.(time.Time); ok {
				since = time.Since(startAt)
			}
		}
	}

	ql.Log().DebugContext(ctx, string(query), "duration", since)
	return nil
}
