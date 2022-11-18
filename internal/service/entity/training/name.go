package statistic

type Name struct {
	value string
}

func (n *Name) Construct(value string) *Name {
	n.value = value

	return n
}

func (n *Name) GetValue() string {
	return n.value
}
