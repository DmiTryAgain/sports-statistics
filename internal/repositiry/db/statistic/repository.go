package statistic

import (
	"database/sql"
	"sports-statistics/internal/app"
	ms "sports-statistics/internal/models/statistic"
	"sports-statistics/internal/service/repository/statistic"
)

type Repository struct {
	db    *sql.DB
	err   error
	model ms.Statistic
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

func (r *Repository) AddStatistic(trainingId int, count int, userId int) {
	insert, err := r.db.Query(
		"INSERT INTO `?` (`?`, `?`, `?`) VALUES(?, ?, ?)",
		r.model.GetTableName(),
		r.model.GetTelegramUserIdColumnName(),
		r.model.GetTrainingIdColumnName(),
		r.model.GetCountColumnName(),
		userId,
		trainingId,
		count,
	)

	if err != nil {
		r.err = err
	}

	err = insert.Scan()

	if err != nil {
		r.err = err
	}
}
func (r *Repository) GetByConditions(trainings []any, periods []string) []ms.Result {

}
