package training

const tableName = "training"
const nameColumn = "name"

type Training struct {
	Id          int
	Alias, Name string
}

func (t *Training) GetTableName() string {
	return tableName
}

func (t *Training) GetNameColumn() string {
	return nameColumn
}
