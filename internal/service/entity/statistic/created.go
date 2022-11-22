package statistic

type Created struct {
	value string
}

func (c *Created) Construct(value string) *Created {
	c.value = value

	return c
}

func (c *Created) GetValue() string {
	return c.value
}
