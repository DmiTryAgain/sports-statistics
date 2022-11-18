package statistic

import (
	"database/sql"
	"sports-statistics/internal/app"
	ms "sports-statistics/internal/models/statistic"
	statisticEntity "sports-statistics/internal/service/entity/statistic"
	trainingEntity "sports-statistics/internal/service/entity/training"
	"sports-statistics/internal/service/entity/user"
	"sports-statistics/internal/service/repository/statistic"
)

type Repository struct {
	db    *sql.DB
	err   error
	model *ms.Statistic
}

func (r *Repository) Construct() statistic.RepositoryInterface {
	r.db, r.err = sql.Open(app.Config.GetDbType(), app.Config.GetDbDsn())

	return r
}

func (r *Repository) Destruct() error {
	return r.db.Close()
}

func (r *Repository) GetError() error {
	return r.err
}

func (r *Repository) AddStatistic(trainingId *trainingEntity.Id, count *statisticEntity.Count, userId *user.Id) {
	insert, err := r.db.Query(
		"INSERT INTO `?` (`?`, `?`, `?`) VALUES(?, ?, ?)",
		r.model.GetTableName(),
		r.model.GetTelegramUserIdColumnName(),
		r.model.GetTrainingIdColumnName(),
		r.model.GetCountColumnName(),
		userId.GetValue(),
		trainingId.GetValue(),
		count.GetValue(),
	)

	if err != nil {
		r.err = err
	}

	err = insert.Scan()

	if err != nil {
		r.err = err
	}
}

func (r *Repository) GetByConditions(trainings []any, periods []string) []statisticEntity.Statistic {

}
