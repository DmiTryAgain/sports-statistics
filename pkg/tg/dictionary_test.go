package tg

import "testing"

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
