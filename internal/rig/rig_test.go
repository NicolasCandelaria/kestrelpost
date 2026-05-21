package rig

import (
	"strings"
	"testing"
)

func TestScanFindsKnownSignal(t *testing.T) {
	signals := DiscoverableSignals()
	s, ok := Scan(signals, BandUtility, 8.412, 8, 1)
	if !ok {
		t.Fatal("expected to discover empty frequency")
	}
	if s.ID != "empty_frequency" {
		t.Fatalf("signal %s", s.ID)
	}
}

func TestRenderWaterfallWidth(t *testing.T) {
	line := RenderWaterfall(60, BandMid, 7.24, DiscoverableSignals(), 4, 1)
	if !strings.Contains(line, "|") {
		t.Fatal("expected tuning cursor")
	}
}
