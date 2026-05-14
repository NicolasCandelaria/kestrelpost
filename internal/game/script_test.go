package game

import "testing"

func TestActForNight(t *testing.T) {
	tests := []struct {
		night, want int
	}{
		{1, 1}, {3, 1}, {4, 2}, {6, 2}, {7, 3}, {9, 3},
	}
	for _, tt := range tests {
		if g := ActForNight(tt.night); g != tt.want {
			t.Fatalf("ActForNight(%d)=%d want %d", tt.night, g, tt.want)
		}
	}
}

func TestNightScript_coversNineNights(t *testing.T) {
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
