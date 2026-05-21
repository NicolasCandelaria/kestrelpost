package rig

import (
	"fmt"
	"math"
	"strings"
)

type Band string

const (
	BandLow     Band = "LOW"
	BandMid     Band = "MID"
	BandHigh    Band = "HIGH"
	BandUtility Band = "UTILITY"
)

type Signal struct {
	ID        string
	Band      Band
	Frequency float64
	MinNight  int
	MaxNight  int
	FromHour  int
	ToHour    int
	Content   string
}

type Tuner struct {
	Band      Band
	Frequency float64
}

func NewTuner() Tuner {
	return Tuner{Band: BandMid, Frequency: 7.24}
}

func (t *Tuner) Tune(delta float64) {
	t.Frequency += delta
	if t.Frequency < 1.8 {
		t.Frequency = 1.8
	}
	if t.Frequency > 30.0 {
		t.Frequency = 30.0
	}
}

func (t *Tuner) SetBand(b Band) {
	t.Band = b
	switch b {
	case BandLow:
		t.Frequency = 2.5
	case BandMid:
		t.Frequency = 7.0
	case BandHigh:
		t.Frequency = 12.5
	case BandUtility:
		t.Frequency = 8.4
	}
}

func DiscoverableSignals() []Signal {
	return []Signal{
		{ID: "numbers_station", Band: BandHigh, Frequency: 14.22, MinNight: 1, MaxNight: 20, FromHour: 2, ToHour: 5, Content: "Numbers station: 7 1 9 3 ..."},
		{ID: "storybook", Band: BandLow, Frequency: 2.31, MinNight: 1, MaxNight: 20, FromHour: 0, ToHour: 23, Content: "A woman reads from a children's storybook."},
		{ID: "coles_channel", Band: BandMid, Frequency: 6.66, MinNight: 6, MaxNight: 20, FromHour: 0, ToHour: 23, Content: "Cole's crew check-in chatter."},
		{ID: "harrow_secondary", Band: BandMid, Frequency: 7.77, MinNight: 7, MaxNight: 20, FromHour: 0, ToHour: 23, Content: "Harrow talks to another post about your traffic patterns."},
		{ID: "empty_frequency", Band: BandUtility, Frequency: 8.412, MinNight: 1, MaxNight: 20, FromHour: 0, ToHour: 23, Content: "Carrier only. Something almost there."},
		{ID: "maritime_distress", Band: BandUtility, Frequency: 8.0, MinNight: 1, MaxNight: 20, FromHour: 0, ToHour: 23, Content: "Distress relay loop with one real call buried inside."},
		{ID: "foreign_broadcast", Band: BandHigh, Frequency: 12.13, MinNight: 7, MaxNight: 20, FromHour: 1, ToHour: 5, Content: "Calm foreign-language message addressed to Operator Seven."},
	}
}

func Scan(signals []Signal, band Band, freq float64, night, hour int) (Signal, bool) {
	for _, s := range signals {
		if s.Band != band || night < s.MinNight || night > s.MaxNight {
			continue
		}
		if hour < s.FromHour || hour > s.ToHour {
			continue
		}
		if math.Abs(s.Frequency-freq) <= 0.05 {
			return s, true
		}
	}
	return Signal{}, false
}

func RenderWaterfall(width int, band Band, tune float64, signals []Signal, night, hour int) string {
	if width < 10 {
		width = 10
	}
	row := make([]rune, width)
	for i := range row {
		row[i] = '.'
	}
	for _, s := range signals {
		if s.Band != band || night < s.MinNight || night > s.MaxNight {
			continue
		}
		if hour < s.FromHour || hour > s.ToHour {
			continue
		}
		pos := int((s.Frequency / 30.0) * float64(width-1))
		if pos < 0 || pos >= width {
			continue
		}
		row[pos] = '#'
		if pos > 0 {
			row[pos-1] = '='
		}
		if pos+1 < width {
			row[pos+1] = '='
		}
	}
	cursor := int((tune / 30.0) * float64(width-1))
	if cursor >= 0 && cursor < width {
		row[cursor] = '|'
	}
	return fmt.Sprintf("%s %s %.3f", string(row), strings.ToUpper(string(band)), tune)
}
