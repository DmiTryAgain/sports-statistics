package tg

import (
	"fmt"
	"strings"
	"time"

	"github.com/DmiTryAgain/sports-statistics/pkg/db"
)

// ExerciseCategory определяет, какие параметры ожидаются у упражнения
type ExerciseCategory int

const (
	CategoryReps           ExerciseCategory = iota // только повторения (текущие упражнения)
	CategoryRepsWeight                             // повторения + вес
	CategoryDistTime                               // дистанция и/или время (хотя бы одно)
	CategoryDuration                               // только время (длительность)
	CategoryDurationWeight                         // время + вес
)

// RequiredParams возвращает список безусловно обязательных типов параметров для категории
func (c ExerciseCategory) RequiredParams() []ParamType {
	switch c {
	case CategoryReps:
		return []ParamType{ParamCount}
	case CategoryRepsWeight:
		return []ParamType{ParamWeight, ParamCount}
	case CategoryDistTime:
		return nil // хотя бы одно из (distance, duration), проверяется в ValidateParams
	case CategoryDuration:
		return []ParamType{ParamDuration}
	case CategoryDurationWeight:
		return []ParamType{ParamWeight, ParamDuration}
	default:
		return nil
	}
}

// SoftRequiredParams возвращает параметры, из которых хотя бы один должен быть заполнен.
// Для CategoryDistTime — дистанция или время.
func (c ExerciseCategory) SoftRequiredParams() []ParamType {
	switch c {
	case CategoryDistTime:
		return []ParamType{ParamDistance, ParamDuration}
	default:
		return nil
	}
}

// ValidateParams проверяет наличие всех обязательных параметров.
// Для CategoryDistTime — хотя бы одно из (distance, duration).
func (c ExerciseCategory) ValidateParams(pp *ParsedParams) error {
	// Проверяем безусловно обязательные
	for _, rp := range c.RequiredParams() {
		if pp.GetParam(rp) == nil {
			return rp.requiredError()
		}
	}

	// Проверяем мягкие требования (хотя бы одно)
	soft := c.SoftRequiredParams()
	if len(soft) > 0 {
		hasSome := false
		for _, sp := range soft {
			if pp.GetParam(sp) != nil {
				hasSome = true
			}
		}
		if !hasSome {
			return errDistOrTimeRequired
		}
	}

	return nil
}

// ParamType — тип параметра упражнения
type ParamType int

const (
	ParamCount    ParamType = iota // количество повторений (голое число без суффикса)
	ParamWeight                    // вес: кг, г
	ParamDistance                  // дистанция: км, м
	ParamDuration                  // длительность: ч, мин, сек
)

// requiredError возвращает ошибку "параметр обязателен" для данного типа
func (pt ParamType) requiredError() error {
	switch pt {
	case ParamCount:
		return errCountRequired
	case ParamWeight:
		return errWeightRequired
	case ParamDistance:
		return errDistanceRequired
	case ParamDuration:
		return errDurationRequired
	}
	return nil
}

// UnitDef описывает единицу измерения и коэффициент нормализации к базовой единице
type UnitDef struct {
	ParamType  ParamType
	Multiplier float64
}

// ParsedValue — результат парсинга числа с опциональной единицей измерения.
type ParsedValue struct {
	Value float64  // Нормализованное значение (с учётом множителя единицы, если есть)
	Unit  *UnitDef // nil, если единица не распознана (голое число)
}

// HasUnit возвращает true, если единица измерения была распознана.
func (p ParsedValue) HasUnit() bool {
	return p.Unit != nil
}

// ParsedParams — результат парсинга текстовых параметров упражнения
type ParsedParams struct {
	Count       *float64 // повторения
	WeightKg    *float64 // вес в кг (нормализован)
	DistanceM   *float64 // дистанция в метрах (нормализована)
	DurationSec *float64 // длительность в секундах (нормализована)
}

// GetParam возвращает указатель на значение параметра по типу
func (pp *ParsedParams) GetParam(pt ParamType) *float64 {
	switch pt {
	case ParamWeight:
		return pp.WeightKg
	case ParamCount:
		return pp.Count
	case ParamDistance:
		return pp.DistanceM
	case ParamDuration:
		return pp.DurationSec
	}
	return nil
}

// SetParam устанавливает значение параметра по типу
func (pp *ParsedParams) SetParam(pt ParamType, val float64) {
	switch pt {
	case ParamWeight:
		pp.WeightKg = ptr(val)
	case ParamCount:
		pp.Count = ptr(val)
	case ParamDistance:
		pp.DistanceM = ptr(val)
	case ParamDuration:
		pp.DurationSec = ptr(val)
	}
}

func (pp *ParsedParams) String() string {
	var s []string
	if pp.Count != nil {
		s = append(s, fmt.Sprintf("count: %f", *pp.Count))
	}
	if pp.WeightKg != nil {
		s = append(s, fmt.Sprintf("weight: %f", *pp.WeightKg))
	}
	if pp.DistanceM != nil {
		s = append(s, fmt.Sprintf("distance: %f", *pp.DistanceM))
	}
	if pp.DurationSec != nil {
		s = append(s, fmt.Sprintf("duration: %f", *pp.DurationSec))
	}

	if len(s) == 0 {
		return "empty params"
	}

	return strings.Join(s, " ")
}

// ToDBParams конвертирует в структуру для БД
func (pp *ParsedParams) ToDBParams() *db.StatisticParams {
	p := &db.StatisticParams{
		WeightKg:    pp.WeightKg,
		DistanceM:   pp.DistanceM,
		DurationSec: pp.DurationSec,
	}
	if p.IsEmpty() {
		return nil
	}
	return p
}

// CountOrDefault возвращает count или 1, если не задан
func (pp *ParsedParams) CountOrDefault() float64 {
	if pp.Count != nil {
		return *pp.Count
	}
	return 1
}

type allChecker interface {
	isAll() bool
}

type language string

type cmd int

type Exercise string

func (e Exercise) String() string { return string(e) }

func (e Exercise) isAll() bool {
	return e == allEx
}

func (e Exercise) isZero() bool {
	var def Exercise
	return e == def
}

func (e Exercise) Category() ExerciseCategory {
	cat, ok := exerciseCategoryMap[e]
	if !ok {
		return CategoryReps
	}
	return cat
}

// OptionalParams возвращает список необязательных параметров для конкретного упражнения
func (e Exercise) OptionalParams() []ParamType {
	if params, ok := exerciseOptionalParamsMap[e]; ok {
		return params
	}
	return nil
}

type Exercises []Exercise

func (e Exercises) String() string {
	b := new(strings.Builder)
	for i := range e {
		b.WriteString("exercise: ")
		b.WriteString(e[i].String())
	}

	return b.String()
}

func (e Exercises) StringSlice() []string {
	res := make([]string, len(e))
	for i := range e {
		res[i] = e[i].String()
	}

	return res
}

func textContainsAllExerciseWords(text string, lang language) bool {
	return textContainsSubstringInMapInAllValByLang(text, exerciseByLang[lang])
}

type textPeriod string

func (tp textPeriod) isAll() bool {
	return tp == allPeriod
}

func textContainsAllPeriodWords(text string, lang language) bool {
	return textContainsSubstringInMapInAllValByLang(text, periodByLang[lang])
}

func textContainsSubstringInMapInAllValByLang[T allChecker](text string, m map[string]T) bool {
	for i, v := range m {
		if strings.Contains(text, i) && v.isAll() {
			return true
		}
	}

	return false
}

type period struct {
	from, to time.Time
}

func (p period) IsZero() bool {
	return p.from.IsZero() && p.to.IsZero()
}

func (p period) ToDB() db.Period {
	return db.Period{
		From: p.from,
		To:   p.to,
	}
}

type periods []period

func (ps periods) ToDB() []db.Period {
	res := make([]db.Period, len(ps))
	for i := range ps {
		res[i] = ps[i].ToDB()
	}

	return res
}

type GroupedStatistic struct {
	db.GroupedStatistic
	TranslatedExercise string
	Category           ExerciseCategory
}

func NewGroupedStatistic(in db.GroupedStatistic, lang language) GroupedStatistic {
	ex := Exercise(in.Exercise)
	return GroupedStatistic{
		GroupedStatistic:   in,
		TranslatedExercise: exTextByLang[lang][ex],
		Category:           ex.Category(),
	}
}

func NewGroupedStatisticList(in []db.GroupedStatistic, lang language) []GroupedStatistic {
	res := make([]GroupedStatistic, len(in))
	for i := range in {
		res[i] = NewGroupedStatistic(in[i], lang)
	}

	return res
}

func anyHasWeight(stats []GroupedStatistic) bool {
	for _, s := range stats {
		if s.WeightKg != nil {
			return true
		}
	}
	return false
}

func anyHasDistance(stats []GroupedStatistic) bool {
	for _, s := range stats {
		if s.DistanceM != nil {
			return true
		}
	}
	return false
}

func anyHasDuration(stats []GroupedStatistic) bool {
	for _, s := range stats {
		if s.SumDurationSec != nil {
			return true
		}
	}
	return false
}

func formatWeight(kg float64, lang language) string {
	suffixKg, suffixG := "кг", "г"
	if lang == langEN {
		suffixKg, suffixG = "kg", "g"
	}
	if kg < 1 {
		return fmt.Sprintf("%.0f%s", kg*1000, suffixG)
	}
	if kg == float64(int(kg)) {
		return fmt.Sprintf("%.0f%s", kg, suffixKg)
	}
	return fmt.Sprintf("%.1f%s", kg, suffixKg)
}

func formatDistance(meters float64, lang language) string {
	suffixKm, suffixM := "км", "м"
	if lang == langEN {
		suffixKm, suffixM = "km", "m"
	}
	if meters >= 1000 {
		km := meters / 1000
		if km == float64(int(km)) {
			return fmt.Sprintf("%.0f%s", km, suffixKm)
		}
		return fmt.Sprintf("%.1f%s", km, suffixKm)
	}
	return fmt.Sprintf("%.0f%s", meters, suffixM)
}

func formatDuration(sec float64, lang language) string {
	suffixH, suffixMin, suffixSec := "ч", "мин", "сек"
	if lang == langEN {
		suffixH, suffixMin, suffixSec = "h", "min", "sec"
	}
	totalSec := int(sec)
	if totalSec < 60 {
		return fmt.Sprintf("%d%s", totalSec, suffixSec)
	}
	if totalSec < 3600 {
		mins := totalSec / 60
		remSec := totalSec % 60
		if remSec == 0 {
			return fmt.Sprintf("%d%s", mins, suffixMin)
		}
		return fmt.Sprintf("%d%s %d%s", mins, suffixMin, remSec, suffixSec)
	}
	hours := totalSec / 3600
	remMin := (totalSec % 3600) / 60
	if remMin == 0 {
		return fmt.Sprintf("%d%s", hours, suffixH)
	}
	return fmt.Sprintf("%d%s %d%s", hours, suffixH, remMin, suffixMin)
}
