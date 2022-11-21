package helpers

import (
	"regexp"
	"sports-statistics/internal/service/entity/statistic"
	"strings"
)

const firstSliceElemIndex = 0
const secondSliceElemIndex = 1

type SliceHelper struct{}

func (h SliceHelper) SplitStringToSlice(text string, separator string) []string {
	needle := regexp.MustCompile(`[[:punct:]]`)
	replacePuncts := needle.ReplaceAllString(strings.ToLower(text), "")
	replaceSpaces := regexp.MustCompile("\\s+")
	replace := replaceSpaces.ReplaceAllString(strings.TrimSpace(replacePuncts), " ")

	return strings.Split(replace, separator)
}

func (h SliceHelper) SplitStringDatesToSlice(interval string) []string {
	replaceSpaces := regexp.MustCompile("\\s+")
	replace := replaceSpaces.ReplaceAllString(strings.TrimSpace(interval), " ")

	return strings.Split(replace, " ")
}

func (h SliceHelper) DeleteElemFromSlice(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

func (h SliceHelper) CheckLenSlice(slice []string, length int) bool {
	return len(slice) == length
}

func (h SliceHelper) IsEmptySlice(slice []string) bool {
	return h.CheckLenSlice(slice, 0)
}

func (h SliceHelper) FirstSliceElemIndex() int {
	return firstSliceElemIndex
}

func (h SliceHelper) SecondSliceElemIndex() int {
	return secondSliceElemIndex
}

func (h SliceHelper) ConvertFromStringToAnyElems(slice []string) []any {
	result := make([]any, len(slice))

	for i, val := range slice {
		result[i] = val
	}

	return result
}

func (h SliceHelper) IsEmptySliceStatisticEntity(slice []*statistic.Statistic) bool {
	return h.CheckLenSliceStatisticEntity(slice, 0)
}

func (h SliceHelper) CheckLenSliceStatisticEntity(slice []*statistic.Statistic, length int) bool {
	return len(slice) == length
}
