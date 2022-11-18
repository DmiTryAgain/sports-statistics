package training

import (
	"database/sql"
	"sports-statistics/internal/app"
	mt "sports-statistics/internal/models/training"
	statisticEntity "sports-statistics/internal/service/entity/statistic"
	trainingEntity "sports-statistics/internal/service/entity/training"
	"sports-statistics/internal/service/repository/training"
)

type Repository struct {
	db    *sql.DB
	err   error
	model *mt.Training
}

func (r *Repository) Construct() training.RepositoryInterface {
	r.db, r.err = sql.Open(app.Config.GetDbType(), app.Config.GetDbDsn())

	return r
}

func (r *Repository) Destruct() error {
	return r.db.Close()
}

func (r *Repository) GetError() error {
	return r.err
}

func (r *Repository) GetTrainingByName(trainingName string) *statisticEntity.Training {
	train := r.db.QueryRow("SELECT * from `?` where `Name` = ? LIMIT 1", r.model.GetTableName(), trainingName)
	err := train.Scan(&r.model.Id, &r.model.Alias, &r.model.Name)

	if err != nil {
		r.err = err
	}

	return r.modelToEntity(r.model)
}

func (r *Repository) modelToEntity(model *mt.Training) *statisticEntity.Training {
	return new(statisticEntity.Training).Construct(
		new(trainingEntity.Id).Construct(model.Id),
		new(trainingEntity.Alias).Construct(model.Alias),
		new(trainingEntity.Name).Construct(model.Name),
	)
}
