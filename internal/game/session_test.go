package game

import (
	"testing"

	"kestrelpost/internal/ending"
)

func TestSession_earlyFuelLoss_deadAir(t *testing.T) {
	s := NewSession()
	s.BeginFromIntro()
	for s.Phase == PhaseNight {
		s.ApplyChoice(1) // −22 fuel each; blackout before night 7 → DEAD AIR
	}
	if s.Phase != PhaseGameOver {
		t.Fatalf("phase = %v want game over", s.Phase)
	}
	if s.State.TerminalDarkNight >= 7 {
		t.Fatalf("terminal night %d want <7 for dead air path", s.State.TerminalDarkNight)
	}
	if g := s.Ending(); g != ending.DeadAir {
		t.Fatalf("ending %v want DeadAir", g)
	}
}

func TestSession_surviveNineNights_relay(t *testing.T) {
	s := NewSession()
	s.BeginFromIntro()
	for s.Phase == PhaseNight {
		s.ApplyChoice(2) // −10/night; nine responses then runway shutdown
	}
	if s.Phase != PhaseGameOver {
		t.Fatalf("phase = %v want game over", s.Phase)
	}
	if g := s.Ending(); g != ending.TheRelay {
		t.Fatalf("ending %v want TheRelay", g)
	}
}
