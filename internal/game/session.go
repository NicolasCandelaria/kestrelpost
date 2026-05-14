package game

import "kestrelpost/internal/ending"

// Phase is the coarse UI mode for a play session.
type Phase byte

const (
	PhaseIntro Phase = iota
	PhaseNight
	PhaseGameOver
)

// Session holds per-connection run state and advances nights + fuel.
type Session struct {
	Phase Phase
	State ending.RunState
}

// NewSession starts night 1 with a full tank; automation is still “offline” narratively.
func NewSession() *Session {
	return &Session{
		Phase: PhaseIntro,
		State: ending.RunState{
			Night: 1,
			Fuel:  100,
		},
	}
}

// BeginFromIntro moves past the welcome screen into the first shift.
func (s *Session) BeginFromIntro() {
	if s.Phase == PhaseIntro {
		s.Phase = PhaseNight
	}
}

// ApplyChoice handles operator response for the current night (keys 1–3).
// Costs and deltas are placeholders until narrative content exists.
func (s *Session) ApplyChoice(choice int) {
	if s.Phase != PhaseNight {
		return
	}
	n := s.State.Night
	switch choice {
	case 1: // long open channel to Maren
		s.State.Fuel -= 22
		s.State.MarenHubSupport += 2
		s.State.MarenTrust += 1
	case 2: // short acknowledgment (tuned so nine nights can finish without early blackout)
		s.State.Fuel -= 10
		s.State.MarenHubSupport += 1
		s.State.MarenTrust += 1
	case 3: // standby — minimal draw, relationship cools
		s.State.Fuel -= 3
		s.State.MarenTrust -= 1
	default:
		return
	}

	if s.State.Fuel <= 0 {
		s.State.Fuel = 0
		s.State.TerminalDarkNight = n
		s.Phase = PhaseGameOver
		return
	}

	s.State.Night++
	if s.State.Night > 9 {
		// Completed the nine-night runway with fuel remaining — force shutdown for resolver.
		s.State.Fuel = 0
		s.State.TerminalDarkNight = 9
		s.Phase = PhaseGameOver
	}
}

// Ending returns the resolved epilogue id for the current State (call in PhaseGameOver).
func (s *Session) Ending() ending.Ending {
	return ending.ResolveEnding(ending.DefaultConfig(), s.State)
}
