package statistic

import (
	"sports-statistics/internal/service/entity/statistic"
)

type RepositoryInterface interface {
	Construct() RepositoryInterface
	Destruct() error
	GetError() error
	AddStatistic(statistic *statistic.Statistic)
	GetByConditions(trainings []any, periods []string, userId int) []*statistic.Statistic
}
