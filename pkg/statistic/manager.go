package statistic

import (
	"context"

	"github.com/DmiTryAgain/sports-statistics/pkg/db"
	"github.com/DmiTryAgain/sports-statistics/pkg/embedlog"

	"github.com/pkg/errors"
)

type Manager struct {
	dbo           db.DB
	statisticRepo db.StatisticRepo
	logger        embedlog.Logger
}

func NewManager(dbo db.DB, logger embedlog.Logger) Manager {
	return Manager{
		dbo:           dbo,
		logger:        logger,
		statisticRepo: db.NewStatisticRepo(dbo),
	}
}

func (m Manager) AddStatistic(ctx context.Context, statistic *Statistic) (*Statistic, error) {
	s, err := m.statisticRepo.AddStatistic(ctx, statistic.ToDB())
	if err != nil {
		return nil, errors.Wrap(err, "failed to add new statistic record")
	}

	return NewStatistic(s), nil
}

func (m Manager) DeleteStatistic(ctx context.Context, id int) (bool, error) {
	if _, err := m.byID(ctx, id); err != nil {
		return false, err
	}

	ok, err := m.statisticRepo.DeleteStatistic(ctx, id)
	if err != nil {
		return false, errors.Wrapf(err, "failed to delete statistic record: %d", id)
	}

	return ok, nil
}

func (m Manager) byID(ctx context.Context, id int) (*Statistic, error) {
	s, err := m.statisticRepo.StatisticByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get statistic by id")
	}

	return NewStatistic(s), nil
}
