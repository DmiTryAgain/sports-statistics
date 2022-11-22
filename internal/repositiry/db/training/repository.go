package training

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sports-statistics/internal/config"
	mt "sports-statistics/internal/models/training"
	statisticEntity "sports-statistics/internal/service/entity/statistic"
	trainingEntity "sports-statistics/internal/service/entity/training"
	"sports-statistics/internal/service/repository/training"
)

type Repository struct {
	db    *sql.DB
	err   error
	model mt.Training
}

func (r *Repository) Construct() training.RepositoryInterface {
	r.db, r.err = sql.Open(config.Configs.GetDbType(), config.Configs.GetDbDsn())
	r.model = mt.Training{}

	return r
}

func (r *Repository) Destruct() error {
	return r.db.Close()
}

func (r *Repository) GetError() error {
	return r.err
}

func (r *Repository) GetTrainingByName(trainingName string) *statisticEntity.Training {
	train := r.db.QueryRow("SELECT * from `"+r.model.GetTableName()+"` where `name` = ? LIMIT 1", trainingName)
	err := train.Scan(&r.model.Id, &r.model.Alias, &r.model.Name)

	if err != nil {
		r.err = err
	}

	return r.modelToEntity(&r.model)
}

func (r *Repository) GetTrainingNames() []*statisticEntity.Training {
	trains, err := r.db.Query("SELECT " + r.model.GetNameColumn() + " from `" + r.model.GetTableName() + "`")

	var results []*statisticEntity.Training
	for trains.Next() {
		var trainName string
		err = trains.Scan(&trainName)

		if err != nil {
			r.err = err
		}

		results = append(
			results,
			new(statisticEntity.Training).Construct(
				nil,
				nil,
				new(trainingEntity.Name).Construct(trainName),
			),
		)
	}

	return results
}

func (r *Repository) modelToEntity(model *mt.Training) *statisticEntity.Training {
	return new(statisticEntity.Training).Construct(
		new(trainingEntity.Id).Construct(model.Id),
		new(trainingEntity.Alias).Construct(model.Alias),
		new(trainingEntity.Name).Construct(model.Name),
	)
}
