package training

import mt "sports-statistics/internal/models/training"

type RepositoryInterface interface {
	Construct() RepositoryInterface
	Destruct() error
	GetError() error
	GetTrainingByName(trainingName string) mt.Training
}
