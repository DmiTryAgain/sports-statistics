package periods

type RepositoryInterface interface {
	Construct() RepositoryInterface
	GetConditionsByPeriod(period string) (string, bool)
	GetConditionsByDate(period string) string
	GetConditionsByDateInterval(from string, to string) string
	GetAllowTextPeriods() []string
}
