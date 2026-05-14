package ending

import "testing"

func TestResolveEnding_priorityConvoyOverBroadcast(t *testing.T) {
	cfg := DefaultConfig()
	s := RunState{
		ConvoyBetrayal:   true,
		OseiFullRelease:  true,
		KidInvestigation: cfg.KMax + 1,
	}
	if g := ResolveEnding(cfg, s); g != TheConvoy {
		t.Fatalf("got %v want TheConvoy", g)
	}
}

func TestResolveEnding_broadcastBeatsKid(t *testing.T) {
	cfg := DefaultConfig()
	s := RunState{
		OseiFullRelease:  true,
		KidInvestigation: cfg.KMax,
	}
	if g := ResolveEnding(cfg, s); g != FullBroadcast {
		t.Fatalf("got %v want FullBroadcast", g)
	}
}

func TestResolveEnding_kidWhenNoBroadcast(t *testing.T) {
	cfg := DefaultConfig()
	s := RunState{
		KidInvestigation: cfg.KMax,
	}
	if g := ResolveEnding(cfg, s); g != TheKidWasRight {
		t.Fatalf("got %v want TheKidWasRight", g)
	}
}

func TestResolveEnding_darkFrequencyAfterKidBlock(t *testing.T) {
	cfg := DefaultConfig()
	s := RunState{
		HarrowDarkPlan:   true,
		KidInvestigation: cfg.KMax,
	}
	if g := ResolveEnding(cfg, s); g != DarkFrequency {
		t.Fatalf("got %v want DarkFrequency", g)
	}
}

func TestResolveEnding_deadAirEarlyTerminalNight(t *testing.T) {
	cfg := DefaultConfig()
	s := RunState{
		Fuel:              0,
		TerminalDarkNight: cfg.DeadAirExclusiveMaxTerminalNight - 1,
	}
	if g := ResolveEnding(cfg, s); g != DeadAir {
		t.Fatalf("got %v want DeadAir", g)
	}
}

func TestResolveEnding_relayLateWithMarenScores(t *testing.T) {
	cfg := DefaultConfig()
	s := RunState{
		Fuel:              0,
		TerminalDarkNight: cfg.RelayMinTerminalNight,
		MarenHubSupport:   cfg.MThreshold,
		MarenTrust:        cfg.TThreshold,
	}
	if g := ResolveEnding(cfg, s); g != TheRelay {
		t.Fatalf("got %v want TheRelay", g)
	}
}

func TestResolveEnding_relayFailsLowTrust(t *testing.T) {
	cfg := DefaultConfig()
	s := RunState{
		Fuel:              0,
		TerminalDarkNight: cfg.RelayMinTerminalNight,
		MarenHubSupport:   cfg.MThreshold,
		MarenTrust:        cfg.TThreshold - 1,
	}
	if g := ResolveEnding(cfg, s); g != Fallback {
		t.Fatalf("got %v want Fallback", g)
	}
}

func TestResolveEnding_relayFailsLowHub(t *testing.T) {
	cfg := DefaultConfig()
	s := RunState{
		Fuel:              0,
		TerminalDarkNight: cfg.RelayMinTerminalNight,
		MarenHubSupport:   cfg.MThreshold - 1,
		MarenTrust:        cfg.TThreshold,
	}
	if g := ResolveEnding(cfg, s); g != Fallback {
		t.Fatalf("got %v want Fallback", g)
	}
}
