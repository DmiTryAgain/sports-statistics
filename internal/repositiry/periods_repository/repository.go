package periods_repository

import (
	"fmt"
	"sports-statistics/internal/models/statistic"
	"sports-statistics/internal/service/repository/periods"
)

type Repository struct {
	conditions map[string]string
	colCreated string
}

func (r *Repository) Construct() periods.RepositoryInterface {
	r.colCreated = statistic.Statistic{}.GetCountColumnCreated()
	r.conditions = map[string]string{
		"сегодня":   fmt.Sprintf("YEAR(`%s`) = YEAR(NOW()) AND WEEK(`%s`, 1) = WEEK(NOW(), 1) AND DAY(`%s`) = DAY(NOW())", r.colCreated, r.colCreated, r.colCreated),
		"вчера":     fmt.Sprintf("MONTH(`%s`) = MONTH(DATE_ADD(NOW(), INTERVAL -1 DAY)) and YEAR(`%s`) = YEAR(DATE_ADD(NOW(), INTERVAL -1 DAY))", r.colCreated, r.colCreated),
		"позавчера": fmt.Sprintf("MONTH(`%s`) = MONTH(DATE_ADD(NOW(), INTERVAL -2 DAY)) and YEAR(`%s`) = YEAR(DATE_ADD(NOW(), INTERVAL -2 DAY))", r.colCreated, r.colCreated),
		"неделю":    fmt.Sprintf("YEAR(`%s`) = YEAR(NOW()) AND WEEK(`%s`, 1) = WEEK(NOW(), 1)", r.colCreated, r.colCreated),
		"месяц":     fmt.Sprintf("MONTH(`%s`) = MONTH(NOW()) AND YEAR(`%s`) = YEAR(NOW())", r.colCreated, r.colCreated),
		"год":       fmt.Sprintf("YEAR(`%s`) = YEAR(NOW())", r.colCreated),
	}

	return r
}
func (r *Repository) GetConditionsByPeriod(period string) (string, bool) {
	res, ok := r.conditions[period]

	return res, ok
}

func (r *Repository) GetConditionsByDate(period string) string {
	return fmt.Sprintf("DATE(`%s`) = DATE('%s')", r.colCreated, period)
}

func (r *Repository) GetConditionsByDateInterval(from string, to string) string {
	return fmt.Sprintf(
		"DATE(`%s`) >= DATE('%s') AND DATE(`%s`) <= DATE('%s')",
		r.colCreated,
		from,
		r.colCreated,
		to,
	)
}
