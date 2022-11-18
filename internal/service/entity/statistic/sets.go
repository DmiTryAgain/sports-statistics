package statistic

type Sets struct {
	value int
}

func (s *Sets) Construct(value int) *Sets {
	s.value = value

	return s
}

func (s *Sets) GetValue() int {
	return s.value
}
