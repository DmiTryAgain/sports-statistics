package tg

import (
	"context"
	"testing"
	"time"
)

func TestSessionStore_GetSetDelete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ss := NewSessionStore(ctx)

	// Get non-existent session
	if s := ss.Get("user1"); s != nil {
		t.Fatalf("expected nil, got %+v", s)
	}

	// Set and Get
	session := &UserSession{State: StateAwaitExercise, Lang: langRU}
	ss.Set("user1", session)

	got := ss.Get("user1")
	if got == nil {
		t.Fatal("expected session, got nil")
	}
	if got.State != StateAwaitExercise {
		t.Fatalf("expected StateAwaitExercise, got %d", got.State)
	}
	if got.Lang != langRU {
		t.Fatalf("expected langRU, got %s", got.Lang)
	}

	// Delete
	ss.Delete("user1")
	if s := ss.Get("user1"); s != nil {
		t.Fatalf("expected nil after delete, got %+v", s)
	}
}

func TestSessionStore_Expiry(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ss := NewSessionStore(ctx)

	session := &UserSession{
		State:     StateAwaitWeight,
		Lang:      langEN,
		UpdatedAt: time.Now().Add(-sessionTTL - time.Minute),
	}

	ss.mu.Lock()
	ss.sessions["user1"] = session
	ss.mu.Unlock()

	// Should be nil because expired
	if s := ss.Get("user1"); s != nil {
		t.Fatalf("expected nil for expired session, got %+v", s)
	}
}

func TestUserSession_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		updatedAt time.Time
		want      bool
	}{
		{
			name:      "not expired",
			updatedAt: time.Now(),
			want:      false,
		},
		{
			name:      "expired",
			updatedAt: time.Now().Add(-sessionTTL - time.Second),
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserSession{UpdatedAt: tt.updatedAt}
			if got := s.IsExpired(sessionTTL); got != tt.want {
				t.Errorf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserSession_Reset(t *testing.T) {
	s := &UserSession{
		State:    StateAwaitWeight,
		Lang:     langRU,
		Exercise: benchPressEx,
		Params:   ParsedParams{WeightKg: ptr[float64](80)},
	}

	s.Reset()

	if s.State != StateIdle {
		t.Errorf("expected StateIdle, got %d", s.State)
	}
	if s.Lang != langRU {
		t.Errorf("expected lang preserved as langRU, got %s", s.Lang)
	}
	if !s.Exercise.isZero() {
		t.Errorf("expected empty exercise, got %s", s.Exercise)
	}
	if s.Params.WeightKg != nil {
		t.Errorf("expected nil WeightKg, got %v", *s.Params.WeightKg)
	}
}
