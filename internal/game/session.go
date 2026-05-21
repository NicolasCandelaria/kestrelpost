package game

import (
	"fmt"

	"kestrelpost/internal/content"
	"kestrelpost/internal/ending"
	"kestrelpost/internal/night"
	"kestrelpost/internal/resources"
	"kestrelpost/internal/rig"
	"kestrelpost/internal/world"
)

// Phase is the coarse UI mode for a play session.
type Phase byte

const (
	PhaseIntro Phase = iota
	PhaseTutorial
	PhaseNight
	PhaseGameOver
)

// Session holds per-connection run state and advances nights + fuel.
type Session struct {
	Phase        Phase
	State        ending.RunState
	Mode         night.Mode
	LastAction   string
	NeedsDogName bool
	Dog          world.Dog
	Power        resources.Power
	Antenna      resources.Antenna
	Pins         resources.Pins
	Tuner        rig.Tuner
	Signals      []rig.Signal
	Campaign     content.Campaign
	CurNight     content.Night
	Logbook      []string
	ScannedText  string
	// TxLog is an append-only transmission log (trimmed server-side).
	TxLog []string
}

const txLogMaxLines = 64

// NewSession starts night 1 with a full tank; automation is still “offline” narratively.
func NewSession() *Session {
	campaign, err := content.LoadCampaign()
	if err != nil {
		panic(err)
	}
	return &Session{
		Phase: PhaseIntro,
		State: ending.RunState{
			Night:           1,
			Fuel:            100,
			DogAlive:        true,
			DogName:         "",
			FrequenciesSeen: map[string]bool{},
		},
		Mode:         night.ModePreShift,
		NeedsDogName: true,
		Dog:          world.NewDog("Old Gray"),
		Power:        resources.NewPower(100, 20),
		Antenna:      resources.NewAntenna(),
		Pins:         resources.NewPins(4),
		Tuner:        rig.NewTuner(),
		Signals:      rig.DiscoverableSignals(),
		Campaign:     campaign,
		LastAction:   "Press [enter] to begin Night 1, or [t] for tutorial.",
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
		s.startNight()
		s.LastAction = "Campaign started. You are in pre-shift."
	}
}

func (s *Session) BeginTutorial() {
	if s.Phase == PhaseIntro {
		s.Phase = PhaseTutorial
		s.appendTx("Tutorial Night 0: previous operator walkthrough started.")
		s.LastAction = "Tutorial started."
	}
}

func (s *Session) FinishTutorial() {
	if s.Phase == PhaseTutorial {
		s.Phase = PhaseNight
		s.startNight()
		s.LastAction = "Tutorial complete. Night 1 started."
	}
}

func (s *Session) startNight() {
	nightData, ok := s.Campaign.Nights[s.State.Night]
	if !ok {
		s.Phase = PhaseGameOver
		return
	}
	s.CurNight = nightData
	s.Mode = night.ModePreShift
	w := world.WeatherForNight(s.State.Night)
	s.State.Weather = string(w)
	f := world.FaultForNight(s.State.Night, w)
	if f != world.FaultNone {
		s.State.NightFaults = append(s.State.NightFaults, string(f))
	}
	s.appendTx(fmt.Sprintf("Night %d weather=%s fault=%s", s.State.Night, w, f))
	s.Dog.AdvanceNight()
	s.State.DogAlive = s.Dog.Alive
	s.State.DogHunger = s.Dog.Hunger
	s.LastAction = fmt.Sprintf("Night %d loaded. Complete pre-shift checks.", s.State.Night)
}

func (s *Session) CurrentNight() content.Night {
	return s.CurNight
}

func (s *Session) CanAdvanceFromPreShift() bool {
	return s.Mode == night.ModePreShift
}

func (s *Session) BeginReceiveWindow() {
	if s.Mode != night.ModePreShift {
		return
	}
	if s.NeedsDogName && s.State.Night == 1 {
		s.LastAction = "Name your dog first: press [1], [2], or [3]."
		return
	}
	_ = s.Antenna.BeginSwitch(resources.AntennaReceive)
	s.Antenna.CompleteSwitch()
	s.Mode = night.ModeReceive
	s.LastAction = "Switched to receive mode."
}

func (s *Session) SetDogName(choice int) bool {
	if s.State.Night != 1 || !s.NeedsDogName {
		return false
	}
	name := ""
	switch choice {
	case 1:
		name = "Scout"
	case 2:
		name = "Ash"
	case 3:
		name = "Bramble"
	default:
		return false
	}
	s.Dog.Name = name
	s.State.DogName = name
	s.NeedsDogName = false
	s.LastAction = fmt.Sprintf("Dog named %s.", name)
	s.appendTx(fmt.Sprintf("Night %d [DOG]: named %s", s.State.Night, name))
	return true
}

func (s *Session) EnterScan() {
	if s.Mode != night.ModeReceive {
		return
	}
	if s.Antenna.BeginSwitch(resources.AntennaScan) != nil {
		return
	}
	s.Antenna.CompleteSwitch()
	s.Mode = night.ModeScan
	s.LastAction = "Switched to scan mode."
}

func (s *Session) ExitScan() {
	if s.Mode != night.ModeScan {
		return
	}
	if s.Antenna.BeginSwitch(resources.AntennaReceive) != nil {
		return
	}
	s.Antenna.CompleteSwitch()
	s.Mode = night.ModeReceive
	s.LastAction = "Returned to receive mode."
}

func (s *Session) FeedDog() {
	s.Dog.Feed()
	s.State.DogHunger = s.Dog.Hunger
	s.appendTx(fmt.Sprintf("Night %d [DOG]: fed %s", s.State.Night, s.Dog.Name))
	s.LastAction = fmt.Sprintf("Fed %s.", s.Dog.Name)
}

func (s *Session) PinSource() {
	_ = s.Pins.Pin(s.CurNight.Source)
	s.State.ThreadsPinned = s.Pins.Items()
	s.LastAction = fmt.Sprintf("Pinned thread: %s.", s.CurNight.Source)
}

func (s *Session) UnpinSource() {
	s.Pins.Unpin(s.CurNight.Source)
	s.State.ThreadsPinned = s.Pins.Items()
	s.LastAction = fmt.Sprintf("Unpinned thread: %s.", s.CurNight.Source)
}

func (s *Session) ScanCurrentFrequency(hour int) {
	if s.Mode != night.ModeScan {
		return
	}
	s.Power.Consume(1)
	s.State.Fuel = s.Power.Fuel()
	sig, ok := rig.Scan(s.Signals, s.Tuner.Band, s.Tuner.Frequency, s.State.Night, hour)
	if !ok {
		s.ScannedText = "No lock; mostly static."
		s.LastAction = "Scan found no clear signal."
		return
	}
	if sig.ID == "harrow_secondary" && s.State.HarrowDarkPlan {
		s.ScannedText = "No lock; Harrow secondary is quiet while cooperating."
		s.LastAction = "Harrow secondary unavailable on this route."
		return
	}
	s.State.FrequenciesSeen[sig.ID] = true
	s.ScannedText = sig.Content
	if sig.ID == "numbers_station" && s.State.Night >= 14 {
		s.ScannedText = "Numbers station partial decode: route grid and timing markers are now readable."
	}
	s.appendTx(fmt.Sprintf("Night %d [SCAN]: discovered %s on %.3f", s.State.Night, sig.ID, sig.Frequency))
	s.LastAction = fmt.Sprintf("Scan locked: %s.", sig.ID)
}

func (s *Session) Waterfall(hour int, width int) string {
	return rig.RenderWaterfall(width, s.Tuner.Band, s.Tuner.Frequency, s.Signals, s.State.Night, hour)
}

func (s *Session) Tune(delta float64) {
	s.Tuner.Tune(delta)
	s.LastAction = fmt.Sprintf("Tuned to %.3f MHz.", s.Tuner.Frequency)
}

func (s *Session) SetBand(index int) {
	switch index {
	case 1:
		s.Tuner.SetBand(rig.BandLow)
	case 2:
		s.Tuner.SetBand(rig.BandMid)
	case 3:
		s.Tuner.SetBand(rig.BandHigh)
	case 4:
		s.Tuner.SetBand(rig.BandUtility)
	}
	s.LastAction = fmt.Sprintf("Band set to %s.", s.Tuner.Band)
}

func (s *Session) LoadLogbookEntry() string {
	entry, err := content.SeededLogbookEntry(s.State.Night)
	if err != nil {
		entry = fmt.Sprintf("Night %d pre-shift note: keep contacts short and precise.", s.State.Night)
	}
	s.Logbook = append(s.Logbook, entry)
	s.LastAction = "Loaded seeded logbook entry."
	return entry
}

// ApplyChoice handles operator response for the current night (keys 1–3).
func (s *Session) ApplyChoice(choice int) {
	if s.Phase != PhaseNight || s.Mode != night.ModeReceive {
		return
	}
	if choice < 1 || choice > 3 {
		return
	}
	n := s.State.Night
	ch := s.CurNight.Choices[choice-1]

	if ch.Effects.Fuel >= 0 {
		s.Power.Consume(ch.Effects.Fuel)
	} else {
		s.Power.Refuel(-ch.Effects.Fuel)
	}
	s.State.Fuel = s.Power.Fuel()
	s.State.MarenHubSupport += ch.Effects.Hub
	s.State.MarenTrust += ch.Effects.Trust
	s.State.KidInvestigation += ch.Effects.Kid
	if ch.Effects.SetHarrow {
		s.State.HarrowDarkPlan = true
		s.State.HarrowDarkNights++
	}
	if ch.Effects.OseiRelease {
		s.State.OseiFullRelease = true
	}
	if ch.Effects.ConvoyBetray {
		s.State.ConvoyBetrayal = true
		s.appendTx(fmt.Sprintf("Night %d [%s]: [%d] %s — extraction call signed.", n, s.CurNight.Source, choice, ch.Text))
		s.State.Fuel = 0
		s.Power.Consume(1000)
		s.State.TerminalDarkNight = n
		s.Phase = PhaseGameOver
		s.LastAction = "You chose convoy extraction."
		return
	}

	s.appendTx(fmt.Sprintf("Night %d [%s]: [%d] %s", n, s.CurNight.Source, choice, ch.Text))
	s.Mode = night.NextAfterChoice(choice)
	s.LastAction = fmt.Sprintf("Selected option %d: %s", choice, ch.Text)
}

func (s *Session) ContinueAfterIncident() {
	if s.Mode != night.ModeIncident {
		return
	}
	s.Mode = night.ModeLogbook
	s.LastAction = "Incident acknowledged. Entering logbook."
}

func (s *Session) WriteLogbook(note string) {
	if s.Mode != night.ModeLogbook {
		return
	}
	note = fmt.Sprintf("Night %d: %s", s.State.Night, note)
	s.Logbook = append(s.Logbook, note)
	s.appendTx(note)
	s.LastAction = "Wrote end-of-shift log entry."
}

func (s *Session) EndNight() {
	if s.Mode != night.ModeLogbook {
		return
	}
	if s.State.Fuel <= 0 {
		s.State.TerminalDarkNight = s.State.Night
		s.Phase = PhaseGameOver
		s.LastAction = "Fuel exhausted. Run ended."
		return
	}
	s.State.Night++
	if s.State.Night > 20 {
		s.State.Fuel = 0
		s.State.TerminalDarkNight = 20
		s.Phase = PhaseGameOver
		s.LastAction = "Night 20 complete. Run ended."
		return
	}
	s.startNight()
	s.LastAction = fmt.Sprintf("Moved to Night %d pre-shift.", s.State.Night)
}

// Ending returns the resolved epilogue id for the current State (call in PhaseGameOver).
func (s *Session) Ending() ending.Ending {
	return ending.ResolveEnding(ending.DefaultConfig(), s.State)
}
