package statistic

import st "sports-statistics/internal/service/entity/training"

type Training struct {
	id    *st.Id
	alias *st.Alias
	name  *st.Name
}

func (t *Training) Construct(
	id *st.Id,
	alias *st.Alias,
	name *st.Name,
) *Training {
	t.id = id
	t.alias = alias
	t.name = name

	return t
}

func (t *Training) GetId() *st.Id {
	return t.id
}

func (t *Training) GetAlias() *st.Id {
	return t.id
}

func (t *Training) GetName() *st.Id {
	return t.id
}
