package tg

import (
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/vmkteam/embedlog"
)

func TestBotTimeZone(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		want    *time.Location
		wantErr bool
	}{
		{
			name:   "time zone by name",
			config: "[Bot]\nTimeZone = \"UTC\"",
			want:   time.UTC,
		},
		{
			name:   "time zone is not set",
			config: "[Bot]\nName = \"bot\"",
			want:   nil,
		},
		{
			name:    "unknown time zone",
			config:  "[Bot]\nTimeZone = \"Mars/Phobos\"",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg Config
			_, err := toml.Decode(tt.config, &cfg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("decode config error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := cfg.Bot.TimeZone.Location; got != tt.want {
				t.Errorf("TimeZone = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("handler falls back to the server time zone", func(t *testing.T) {
		h := New(t.Context(), embedlog.NewDevLogger(), dbc, mh.statRepo, nil, Bot{})
		if h.loc != time.Local {
			t.Errorf("handler time zone = %v, want %v", h.loc, time.Local)
		}
	})
}
