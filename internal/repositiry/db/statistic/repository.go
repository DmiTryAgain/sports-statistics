package statistic

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sports-statistics/internal/config"
	ms "sports-statistics/internal/models/statistic"
	statisticEntity "sports-statistics/internal/service/entity/statistic"
	trainingEntity "sports-statistics/internal/service/entity/training"
	"sports-statistics/internal/service/entity/user"
	"sports-statistics/internal/service/repository/statistic"
	"strings"
)

type Repository struct {
	db    *sql.DB
	err   error
	model ms.Statistic
}

func (r *Repository) Construct() statistic.RepositoryInterface {
	r.db, r.err = sql.Open(config.Configs.GetDbType(), config.Configs.GetDbDsn())
	r.model = ms.Statistic{}

	return r
}

func (r *Repository) Destruct() error {
	return r.db.Close()
}

func (r *Repository) GetError() error {
	return r.err
}

func (r *Repository) AddStatistic(statistic *statisticEntity.Statistic) {
	query := `INSERT INTO ` + r.model.GetTableName() + ` (
		` + r.model.GetTelegramUserIdColumnName() + `, 
		` + r.model.GetTrainingIdColumnName() + `, 
		` + r.model.GetCountColumnName() +
		`) 
	VALUES(?, ?, ?)`

	insert, err := r.db.Query(
		query,
		statistic.GetUser().GetId().GetValue(),
		statistic.GetTraining().GetId().GetValue(),
		statistic.GetCount().GetValue(),
	)

	if err != nil {
		r.err = err
	}

	insert.Close()
}

func (r *Repository) GetByConditions(trainings []any, periods []string, userId int) []*statisticEntity.Statistic {
	periodsSQL := strings.Join(periods, " OR ")
	query := `
				SELECT t.name as train, sum(count) as total, count(count) as sets 
				FROM statistic 
				JOIN training t on t.id = statistic.training_id 
				WHERE telegram_user_id = ? 
				AND ` + periodsSQL + ` 
 				AND t.name in (?` + strings.Repeat(",?", len(trainings)-1) + `)
				GROUP BY training_id;
			`
	//fmt.Println(query)

	result, err := r.db.Query(
		query,
		append(append([]any{}, userId), trainings...)...,
	)

	//log.Printf("Search result:  %v \n", result)
	if err != nil {
		r.err = err
	}

	var results []*statisticEntity.Statistic
	for result.Next() {
		var train string
		var total int
		var sets int
		err = result.Scan(&train, &total, &sets)

		if err != nil {
			r.err = err
		}

		results = append(
			results,
			new(statisticEntity.Statistic).Construct(
				nil,
				new(statisticEntity.Training).Construct(nil, nil, new(trainingEntity.Name).Construct(train)),
				new(statisticEntity.User).Construct(new(user.Id).Construct(userId)),
				nil,
				new(statisticEntity.Count).Construct(total),
				new(statisticEntity.Sets).Construct(sets),
			),
		)
	}

	return results

}
