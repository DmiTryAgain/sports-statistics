package helpers

import "strings"

type StringHelper struct{}

func (s *StringHelper) KeyMapToString(m map[string]any, separator string) string {
	sliceKeys := make([]string, len(m))
	i := firstSliceElemIndex

	for key, _ := range m {
		sliceKeys[i] = key
		i++
	}

	return strings.Join(sliceKeys, separator)
}
