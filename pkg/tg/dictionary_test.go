package tg

import (
	"slices"
	"testing"
)

func TestExerciseRegistrationCompleteness(t *testing.T) {
	orderSet := make(map[Exercise]struct{}, len(exerciseOrder))
	for _, ex := range exerciseOrder {
		orderSet[ex] = struct{}{}
	}

	// Reverse index: which exercises have synonyms
	synonymsRU := make(map[Exercise]struct{})
	for _, ex := range exerciseByLang[langRU] {
		synonymsRU[ex] = struct{}{}
	}
	synonymsEN := make(map[Exercise]struct{})
	for _, ex := range exerciseByLang[langEN] {
		synonymsEN[ex] = struct{}{}
	}

	// Forward: every exercise in exerciseOrder must be in mandatory maps
	for _, ex := range exerciseOrder {
		if exTextByLang[langRU][ex] == "" {
			t.Errorf("exercise %q: missing from exTextByLang[langRU]", ex)
		}
		if exTextByLang[langEN][ex] == "" {
			t.Errorf("exercise %q: missing from exTextByLang[langEN]", ex)
		}
		if _, ok := synonymsRU[ex]; !ok {
			t.Errorf("exercise %q: no synonym in exerciseByLang[langRU]", ex)
		}
		if _, ok := synonymsEN[ex]; !ok {
			t.Errorf("exercise %q: no synonym in exerciseByLang[langEN]", ex)
		}
	}

	// Reverse: no orphaned entries in maps
	for ex := range exTextByLang[langRU] {
		if _, ok := orderSet[ex]; !ok {
			t.Errorf("exercise %q: in exTextByLang[langRU] but not in exerciseOrder", ex)
		}
	}
	for ex := range exTextByLang[langEN] {
		if _, ok := orderSet[ex]; !ok {
			t.Errorf("exercise %q: in exTextByLang[langEN] but not in exerciseOrder", ex)
		}
	}
	for ex := range exerciseCategoryMap {
		if _, ok := orderSet[ex]; !ok {
			t.Errorf("exercise %q: in exerciseCategoryMap but not in exerciseOrder", ex)
		}
	}
	for ex := range exerciseOptionalParamsMap {
		if _, ok := orderSet[ex]; !ok {
			t.Errorf("exercise %q: in exerciseOptionalParamsMap but not in exerciseOrder", ex)
		}
	}
}

func TestExerciseOrderNoDuplicates(t *testing.T) {
	seen := make(map[Exercise]struct{}, len(exerciseOrder))
	for _, ex := range exerciseOrder {
		if _, ok := seen[ex]; ok {
			t.Errorf("duplicate exercise in exerciseOrder: %q", ex)
		}
		seen[ex] = struct{}{}
	}
}

func TestExerciseByLangSynonymConsistency(t *testing.T) {
	orderSet := make(map[Exercise]struct{}, len(exerciseOrder))
	for _, ex := range exerciseOrder {
		orderSet[ex] = struct{}{}
	}

	for lang, synonyms := range exerciseByLang {
		for text, ex := range synonyms {
			if ex == allEx {
				continue
			}
			if _, ok := orderSet[ex]; !ok {
				t.Errorf("exerciseByLang[%s][%q] maps to %q which is not in exerciseOrder", lang, text, ex)
			}
		}
	}
}

func TestGroupingRegistrationCompleteness(t *testing.T) {
	for _, lang := range []language{langRU, langEN} {
		for _, g := range groupingOrder {
			if groupingTextByLang[lang][g] == "" {
				t.Errorf("grouping %q: missing from groupingTextByLang[%s]", g, lang)
			}
		}

		synonyms := make(map[statGrouping]struct{})
		for text, g := range groupingByLang[lang] {
			synonyms[g] = struct{}{}
			if !slices.Contains(groupingOrder, g) {
				t.Errorf("groupingByLang[%s][%q] maps to %q which is not in groupingOrder", lang, text, g)
			}
		}
		for _, g := range groupingOrder {
			if _, ok := synonyms[g]; !ok && g != noGrouping {
				t.Errorf("grouping %q: no phrase in groupingByLang[%s]", g, lang)
			}
		}
	}
}

// Фразы группировки отрезаются от конца запроса статистики до разбора упражнений и периодов,
// поэтому название упражнения или периода не должно заканчиваться такой фразой
func TestGroupingPhrasesDoNotCutExercisesAndPeriods(t *testing.T) {
	for _, lang := range []language{langRU, langEN} {
		for text := range exerciseByLang[lang] {
			if cleaned, g := extractGrouping(text, lang); cleaned != text || g != noGrouping {
				t.Errorf("exerciseByLang[%s][%q] is cut to %q with grouping %q", lang, text, cleaned, g)
			}
		}
		for text := range periodByLang[lang] {
			if cleaned, g := extractGrouping(text, lang); cleaned != text || g != noGrouping {
				t.Errorf("periodByLang[%s][%q] is cut to %q with grouping %q", lang, text, cleaned, g)
			}
		}
	}
}
