package statistic

const tableName = "statistic"
const ColumnTelegramUserId = "telegram_user_id"
const ColumnTrainingId = "training_id"
const ColumnCount = "count"

type Statistic struct {
	Id, TelegramUserId, TrainingId, count int
	Created                               string
}

func (t *Statistic) GetTableName() string {
	return tableName
}

func (t *Statistic) GetTelegramUserIdColumnName() string {
	return ColumnTelegramUserId
}

func (t *Statistic) GetTrainingIdColumnName() string {
	return ColumnTrainingId
}

func (t *Statistic) GetCountColumnName() string {
	return ColumnCount
}
