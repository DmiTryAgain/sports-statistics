package statistic

type Alias struct {
	value string
}

func (a *Alias) Construct(value string) *Alias {
	a.value = value

	return a
}

func (a *Alias) GetValue() string {
	return a.value
}
