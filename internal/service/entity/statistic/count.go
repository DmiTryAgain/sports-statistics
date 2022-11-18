package statistic

type Count struct {
	value int
}

func (c *Count) Construct(value int) *Count {
	c.value = value

	return c
}

func (c *Count) GetValue() int {
	return c.value
}
