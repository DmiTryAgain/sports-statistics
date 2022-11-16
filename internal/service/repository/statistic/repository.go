package statistic

import "sports-statistics/internal/models/statistic"

type RepositoryInterface interface {
	Construct() RepositoryInterface
	Destruct() error
	GetError() error
	AddStatistic(trainingId int, count int, userId int)
	GetByConditions(trainings []any, periods []string) []statistic.Result
}
