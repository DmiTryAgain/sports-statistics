package training

const tableName = "training"

type Training struct {
	Id          int
	Alias, Name string
}

func (t *Training) GetTableName() string {
	return tableName
}
