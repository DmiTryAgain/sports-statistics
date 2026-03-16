package db

type StatisticParams struct {
	WeightKg    *float64 `json:"weightKg,omitempty"`
	DistanceM   *float64 `json:"distanceM,omitempty"`
	DurationSec *float64 `json:"durationSec,omitempty"`
}

func (sp *StatisticParams) IsEmpty() bool {
	if sp == nil {
		return true
	}
	return sp.WeightKg == nil && sp.DistanceM == nil && sp.DurationSec == nil
}
