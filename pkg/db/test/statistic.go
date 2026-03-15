//nolint:dupl,funlen
package test

import (
	"testing"
	"time"

	"github.com/DmiTryAgain/sports-statistics/pkg/db"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-pg/pg/v10/orm"
)

type StatisticOpFunc func(t *testing.T, dbo orm.DB, in *db.Statistic) Cleaner

func Statistic(t *testing.T, dbo orm.DB, in *db.Statistic, ops ...StatisticOpFunc) (*db.Statistic, Cleaner) {
	repo := db.NewStatisticRepo(dbo)
	var cleaners []Cleaner

	// Fill the incoming entity
	if in == nil {
		in = &db.Statistic{}
	}

	// Check if PKs are provided
	if in.ID != 0 {
		// Fetch the entity by PK
		statistic, err := repo.StatisticByID(t.Context(), in.ID, repo.FullStatistic())
		if err != nil {
			t.Fatal(err)
		}
		// Return if found without real cleanup
		if statistic != nil {
			return statistic, emptyClean
		}

		// If we're here, we don't find the entity by PKs. Just try to add the entity by provided PK
		t.Logf("the entity Statistic is not found by provided PKs:ID=%v. Trying to create one", in.ID)
	}

	for _, op := range ops {
		if cl := op(t, dbo, in); cl != nil {
			cleaners = append(cleaners, cl)
		}
	}

	// Create the main entity
	statistic, err := repo.AddStatistic(t.Context(), in, db.WithoutColumns())
	if err != nil {
		t.Fatal(err)
	}

	return statistic, func() {
		if _, err := dbo.ModelContext(t.Context(), &db.Statistic{ID: statistic.ID}).WherePK().Delete(); err != nil {
			t.Fatal(err)
		}
		// Clean up related entities from the last to the first
		for i := len(cleaners) - 1; i >= 0; i-- {
			cleaners[i]()
		}
	}
}

func WithFakeStatistic(t *testing.T, dbo orm.DB, in *db.Statistic) Cleaner {
	if in.TgUserID == "" {
		in.TgUserID = cutS(gofakeit.Sentence(10), 255)
	}

	if in.Exercise == "" {
		in.Exercise = cutS(gofakeit.Sentence(10), 255)
	}

	if in.CreatedAt.IsZero() {
		in.CreatedAt = time.Now()
	}

	if in.StatusID == 0 {
		in.StatusID = 1
	}

	if in.Count == 0 {
		in.Count = gofakeit.Float64Range(1, 10)
	}

	return emptyClean
}
