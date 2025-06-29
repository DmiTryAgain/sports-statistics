package db

import (
	"context"
	"errors"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type StatisticRepo struct {
	db      orm.DB
	filters map[string][]Filter
	sort    map[string][]SortField
	join    map[string][]string
}

// NewStatisticRepo returns new repository
func NewStatisticRepo(db orm.DB) StatisticRepo {
	return StatisticRepo{
		db: db,
		filters: map[string][]Filter{
			Tables.Statistic.Name: {StatusFilter},
		},
		sort: map[string][]SortField{
			Tables.Statistic.Name: {{Column: Columns.Statistic.CreatedAt, Direction: SortDesc}},
		},
		join: map[string][]string{
			Tables.Statistic.Name: {TableColumns},
		},
	}
}

// WithTransaction is a function that wraps StatisticRepo with pg.Tx transaction.
func (sr StatisticRepo) WithTransaction(tx *pg.Tx) StatisticRepo {
	sr.db = tx
	return sr
}

// WithEnabledOnly is a function that adds "statusId"=1 as base filter.
func (sr StatisticRepo) WithEnabledOnly() StatisticRepo {
	f := make(map[string][]Filter, len(sr.filters))
	for i := range sr.filters {
		f[i] = make([]Filter, len(sr.filters[i]))
		copy(f[i], sr.filters[i])
		f[i] = append(f[i], StatusEnabledFilter)
	}
	sr.filters = f

	return sr
}

/*** Statistic ***/

// FullStatistic returns full joins with all columns
func (sr StatisticRepo) FullStatistic() OpFunc {
	return WithColumns(sr.join[Tables.Statistic.Name]...)
}

// DefaultStatisticSort returns default sort.
func (sr StatisticRepo) DefaultStatisticSort() OpFunc {
	return WithSort(sr.sort[Tables.Statistic.Name]...)
}

// StatisticByID is a function that returns Statistic by ID(s) or nil.
func (sr StatisticRepo) StatisticByID(ctx context.Context, id int, ops ...OpFunc) (*Statistic, error) {
	return sr.OneStatistic(ctx, &StatisticSearch{ID: &id}, ops...)
}

// OneStatistic is a function that returns one Statistic by filters. It could return pg.ErrMultiRows.
func (sr StatisticRepo) OneStatistic(ctx context.Context, search *StatisticSearch, ops ...OpFunc) (*Statistic, error) {
	obj := &Statistic{}
	err := buildQuery(ctx, sr.db, obj, search, sr.filters[Tables.Statistic.Name], PagerTwo, ops...).Select()

	if errors.Is(err, pg.ErrMultiRows) {
		return nil, err
	} else if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}

	return obj, err
}

// StatisticsByFilters returns Statistic list.
func (sr StatisticRepo) StatisticsByFilters(ctx context.Context, search *StatisticSearch, pager Pager, ops ...OpFunc) (statistics []Statistic, err error) {
	err = buildQuery(ctx, sr.db, &statistics, search, sr.filters[Tables.Statistic.Name], pager, ops...).Select()
	return
}

// CountStatistics returns count
func (sr StatisticRepo) CountStatistics(ctx context.Context, search *StatisticSearch, ops ...OpFunc) (int, error) {
	return buildQuery(ctx, sr.db, &Statistic{}, search, sr.filters[Tables.Statistic.Name], PagerOne, ops...).Count()
}

// AddStatistic adds Statistic to DB.
func (sr StatisticRepo) AddStatistic(ctx context.Context, statistic *Statistic, ops ...OpFunc) (*Statistic, error) {
	q := sr.db.ModelContext(ctx, statistic)
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.Statistic.CreatedAt)
	}
	applyOps(q, ops...)
	_, err := q.Insert()

	return statistic, err
}

// UpdateStatistic updates Statistic in DB.
func (sr StatisticRepo) UpdateStatistic(ctx context.Context, statistic *Statistic, ops ...OpFunc) (bool, error) {
	q := sr.db.ModelContext(ctx, statistic).WherePK()
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.Statistic.CreatedAt)
	}
	applyOps(q, ops...)
	res, err := q.Update()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}

// DeleteStatistic set statusId to deleted in DB.
func (sr StatisticRepo) DeleteStatistic(ctx context.Context, id int) (deleted bool, err error) {
	statistic := &Statistic{ID: id, StatusID: StatusDeleted}

	return sr.UpdateStatistic(ctx, statistic, WithColumns(Columns.Statistic.StatusID))
}
