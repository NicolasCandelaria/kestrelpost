package game

import (
	"testing"

	"kestrelpost/internal/ending"
	"kestrelpost/internal/night"
)

func TestSession_earlyFuelLoss_deadAir(t *testing.T) {
	s := NewSession()
	s.BeginFromIntro()
	if !s.SetDogName(1) {
		t.Fatal("expected dog naming on night 1")
	}
	for s.Phase == PhaseNight && s.State.Night <= 4 {
		s.BeginReceiveWindow()
		s.ApplyChoice(1)
		s.ContinueAfterIncident()
		s.EndNight()
	}
	s.Power.Consume(1000)
	s.State.Fuel = 0
	s.Mode = night.ModeLogbook
	s.EndNight()
	if s.Phase != PhaseGameOver {
		t.Fatalf("phase = %v want game over", s.Phase)
	}
	if s.State.TerminalDarkNight >= 16 {
		t.Fatalf("terminal night %d want <16 for gone dark path", s.State.TerminalDarkNight)
	}
	if g := s.Ending(); g != ending.DeadAir {
		t.Fatalf("ending %v want GoneDark/DeadAir", g)
	}
}

func TestSession_surviveLateRun_relay(t *testing.T) {
	s := NewSession()
	s.BeginFromIntro()
	_ = s.SetDogName(1)
	s.State.Night = 16
	s.startNight()
	s.State.MarenHubSupport = 20
	s.State.MarenTrust = 10
	s.Power.Consume(1000)
	s.State.Fuel = 0
	s.Mode = night.ModeLogbook
	s.EndNight()
	if s.Phase != PhaseGameOver {
		t.Fatalf("phase = %v want game over", s.Phase)
	}
	if g := s.Ending(); g != ending.TheRelay {
		t.Fatalf("ending %v want TheRelay", g)
	}
}
