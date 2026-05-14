package game

import (
	"fmt"

	"kestrelpost/internal/ending"
)

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
	// TxLog is an append-only transmission log (trimmed server-side).
	TxLog []string
}

const txLogMaxLines = 64

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

func (s *Session) appendTx(line string) {
	s.TxLog = append(s.TxLog, line)
	if len(s.TxLog) > txLogMaxLines {
		s.TxLog = s.TxLog[len(s.TxLog)-txLogMaxLines:]
	}
}

// BeginFromIntro moves past the welcome screen into the first shift.
func (s *Session) BeginFromIntro() {
	if s.Phase == PhaseIntro {
		s.Phase = PhaseNight
	}
}

// CurrentNightCard is the script beat for the active night.
func (s *Session) CurrentNightCard() NightCard {
	if s.Phase != PhaseNight {
		return NightScript(1)
	}
	return NightScript(s.State.Night)
}

// ApplyChoice handles operator response for the current night (keys 1–3).
func (s *Session) ApplyChoice(choice int) {
	if s.Phase != PhaseNight {
		return
	}
	if choice < 1 || choice > 3 {
		return
	}
	n := s.State.Night
	card := NightScript(n)
	ch := card.Choices[choice-1]

	s.State.Fuel -= ch.Fuel
	s.State.MarenHubSupport += ch.Hub
	s.State.MarenTrust += ch.Trust
	s.State.KidInvestigation += ch.Kid
	if ch.SetHarrow {
		s.State.HarrowDarkPlan = true
	}
	if ch.OseiRelease {
		s.State.OseiFullRelease = true
	}
	if ch.ConvoyBetray {
		s.State.ConvoyBetrayal = true
		s.appendTx(fmt.Sprintf("Night %d [%s]: [%d] %s — extraction call signed.", n, card.Source, choice, ch.Label))
		s.State.Fuel = 0
		s.State.TerminalDarkNight = n
		s.Phase = PhaseGameOver
		return
	}

	s.appendTx(fmt.Sprintf("Night %d [%s]: [%d] %s", n, card.Source, choice, ch.Label))

	if s.State.Fuel <= 0 {
		s.State.Fuel = 0
		s.State.TerminalDarkNight = n
		s.Phase = PhaseGameOver
		return
	}

	s.State.Night++
	if s.State.Night > 9 {
		s.State.Fuel = 0
		s.State.TerminalDarkNight = 9
		s.Phase = PhaseGameOver
		return
	}
}

// Ending returns the resolved epilogue id for the current State (call in PhaseGameOver).
func (s *Session) Ending() ending.Ending {
	return ending.ResolveEnding(ending.DefaultConfig(), s.State)
}
