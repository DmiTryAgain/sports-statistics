package statistic

import (
	"sports-statistics/internal/service/entity/statistic"
	trainingEntity "sports-statistics/internal/service/entity/training"
	"sports-statistics/internal/service/entity/user"
)

type RepositoryInterface interface {
	Construct() RepositoryInterface
	Destruct() error
	GetError() error
	AddStatistic(trainingId *trainingEntity.Id, count *statistic.Count, userId *user.Id)
	GetByConditions(trainings []any, periods []string) []statistic.Statistic
}
