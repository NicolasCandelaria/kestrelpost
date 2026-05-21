package game

import (
	"strings"
	"testing"
)

func TestActForNight(t *testing.T) {
	tests := []struct {
		night, want int
	}{
		{1, 1}, {5, 1}, {6, 2}, {10, 2}, {11, 3}, {15, 3}, {16, 4}, {20, 4},
	}
	for _, tt := range tests {
		if g := ActForNight(tt.night); g != tt.want {
			t.Fatalf("ActForNight(%d)=%d want %d", tt.night, g, tt.want)
		}
	}
}

func TestNightScript_coversLegacyNineNights(t *testing.T) {
	for n := 1; n <= 9; n++ {
		c := NightScript(n)
		if c.Source == "" {
			t.Fatalf("night %d empty source", n)
		}
		for i := range c.Choices {
			if c.Choices[i].Reply == "" || c.Choices[i].Fuel <= 0 {
				t.Fatalf("night %d choice %d invalid", n, i+1)
			}
		}
	}
}

func TestNightScript_dialogueAvoidsVaguePhrases(t *testing.T) {
	banned := []string{
		"diesel on the wind",
		"hungry ear",
		"sky grammar",
		"spend yourself",
		"air you refuse to own",
		"corridor truth",
		"lighthouse",
		"charity scales",
	}

	for n := 1; n <= 9; n++ {
		c := NightScript(n)
		text := strings.ToLower(c.Quote)
		for i := range c.Choices {
			text += "\n" + strings.ToLower(c.Choices[i].Reply)
		}
		for _, phrase := range banned {
			if strings.Contains(text, phrase) {
				t.Fatalf("night %d contains vague phrase %q", n, phrase)
			}
		}
	}
}
