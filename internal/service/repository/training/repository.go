package training

import (
	"sports-statistics/internal/service/entity/statistic"
)

type RepositoryInterface interface {
	Construct() RepositoryInterface
	Destruct() error
	GetError() error
	GetTrainingByName(trainingName string) *statistic.Training
	GetTrainingNames() []*statistic.Training
}
