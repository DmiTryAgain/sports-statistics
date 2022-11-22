package user

type Id struct {
	value int
}

func (id *Id) Construct(value int) *Id {
	id.value = value

	return id
}

func (id *Id) GetValue() int {
	return id.value
}
