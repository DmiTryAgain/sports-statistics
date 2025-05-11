package statistic

import (
	"github.com/DmiTryAgain/sports-statistics/pkg/db"
)

type Statistic struct {
	db.Statistic
}

func (s *Statistic) ToDB() *db.Statistic {
	if s == nil {
		return nil
	}

	return &db.Statistic{
		TgUserID: s.TgUserID,
		Exercise: s.Exercise,
		Value:    s.Value,
		StatusID: s.StatusID,
	}
}

func NewStatistic(in *db.Statistic) *Statistic {
	if in == nil {
		return nil
	}

	return &Statistic{
		Statistic: *in,
	}
}
