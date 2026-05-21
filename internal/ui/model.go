package ui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"kestrelpost/internal/ending"
	"kestrelpost/internal/game"
	"kestrelpost/internal/night"
	"kestrelpost/internal/save"
)

const txLogViewLines = 6
const minTerminalWidth = 80
const minTerminalHeight = 24

type waterfallTickMsg struct{}

// Model is the root Bubble Tea model for one SSH session.
type Model struct {
	session *game.Session
	width   int // from WindowSizeMsg; chrome uses effectiveWidth(width)
	height  int
}

func NewModel() *Model {
	return &Model{session: game.NewSession(), width: 80, height: 24}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(tea.RequestWindowSize, tickWaterfall())
}

func tickWaterfall() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(_ time.Time) tea.Msg {
		return waterfallTickMsg{}
	})
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Bubble Tea v2: presses are [tea.KeyPressMsg], releases are [tea.KeyReleaseMsg].
	// Some SSH/terminal paths may surface other [tea.KeyMsg] implementations.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width > 0 {
			m.width = msg.Width
		}
		if msg.Height > 0 {
			m.height = msg.Height
		}
		return m, nil
	case waterfallTickMsg:
		if m.session.Phase == game.PhaseNight && m.session.Mode == night.ModeScan {
			m.session.ScanCurrentFrequency(2)
			return m, tickWaterfall()
		}
		return m, tickWaterfall()
	case tea.KeyPressMsg:
		return m.dispatchKey(msg.String(), msg.Key())
	case tea.KeyReleaseMsg:
		// Some clients/PTY stacks surface releases without a press; allow intro to advance.
		if m.session.Phase == game.PhaseIntro {
			return m.dispatchKey(msg.String(), msg.Key())
		}
		return m, nil
	case tea.KeyMsg:
		return m.dispatchKey(msg.String(), msg.Key())
	}
	return m, nil
}

func (m *Model) dispatchKey(ks string, k tea.Key) (tea.Model, tea.Cmd) {
	if isQuitKey(ks, k) {
		return m, tea.Quit
	}

	switch m.session.Phase {
	case game.PhaseIntro:
		if ks == "t" {
			m.session.BeginTutorial()
		} else if ks == "o" || introConfirmKeys(ks, k) {
			m.session.BeginFromIntro()
		}
		return m, nil
	case game.PhaseTutorial:
		if introConfirmKeys(ks, k) {
			m.session.FinishTutorial()
		}
		return m, nil
	case game.PhaseNight:
		switch m.session.Mode {
		case night.ModePreShift:
			handled := false
			if m.session.NeedsDogName && m.session.State.Night == 1 {
				if choice, ok := nightChoiceFromKey(ks, k); ok {
					handled = m.session.SetDogName(choice)
				}
				if !handled {
					m.session.LastAction = "Name your dog first: [1] Scout, [2] Ash, [3] Bramble."
				}
				return m, nil
			}
			if ks == "f" {
				handled = true
				m.session.FeedDog()
			} else if ks == "k" {
				handled = true
				_ = save.Save("", save.Snapshot{
					State:       m.session.State,
					Logbook:     append([]string(nil), m.session.Logbook...),
					UnlockCount: len(m.session.State.FrequenciesSeen),
				})
				m.session.LastAction = "Saved snapshot."
			} else if ks == "o" {
				handled = true
				if snap, err := save.Load(""); err == nil {
					m.session.State = snap.State
					m.session.Logbook = append([]string(nil), snap.Logbook...)
					m.session.LastAction = "Loaded snapshot."
				} else {
					m.session.LastAction = "Load failed: no snapshot found."
				}
			} else if ks == "p" {
				handled = true
				m.session.PinSource()
			} else if ks == "u" {
				handled = true
				m.session.UnpinSource()
			} else if ks == "l" {
				handled = true
				m.session.LoadLogbookEntry()
			} else if introConfirmKeys(ks, k) {
				handled = true
				m.session.BeginReceiveWindow()
			}
			if !handled {
				m.session.LastAction = fmt.Sprintf("No action for key %q in pre-shift.", displayKey(ks))
			}
		case night.ModeReceive:
			handled := false
			if ks == "s" {
				handled = true
				m.session.EnterScan()
				return m, nil
			}
			if choice, ok := nightChoiceFromKey(ks, k); ok {
				handled = true
				m.session.ApplyChoice(choice)
			}
			if !handled {
				m.session.LastAction = fmt.Sprintf("No action for key %q in receive mode.", displayKey(ks))
			}
		case night.ModeScan:
			handled := true
			switch ks {
			case "left":
				m.session.Tune(-0.1)
			case "right":
				m.session.Tune(0.1)
			case "S-left":
				m.session.Tune(-0.01)
			case "S-right":
				m.session.Tune(0.01)
			case "1", "2", "3", "4":
				switch ks {
				case "1":
					m.session.SetBand(1)
				case "2":
					m.session.SetBand(2)
				case "3":
					m.session.SetBand(3)
				case "4":
					m.session.SetBand(4)
				}
			case "r":
				m.session.ExitScan()
			case "enter":
				m.session.ScanCurrentFrequency(2)
			default:
				handled = false
			}
			if !handled {
				m.session.LastAction = fmt.Sprintf("No action for key %q in scan mode.", displayKey(ks))
			}
		case night.ModeIncident:
			if introConfirmKeys(ks, k) {
				m.session.ContinueAfterIncident()
			} else {
				m.session.LastAction = fmt.Sprintf("Press [enter] to continue (got %q).", displayKey(ks))
			}
		case night.ModeLogbook:
			handled := false
			if ks == "w" {
				handled = true
				m.session.WriteLogbook("operator note filed.")
			}
			if ks == "n" || introConfirmKeys(ks, k) {
				handled = true
				m.session.EndNight()
			}
			if !handled {
				m.session.LastAction = fmt.Sprintf("No action for key %q in logbook mode.", displayKey(ks))
			}
		}
		return m, nil
	case game.PhaseGameOver:
		return m, nil
	}
	return m, nil
}

func isQuitKey(ks string, k tea.Key) bool {
	if ks == "ctrl+c" || ks == "q" {
		return true
	}
	if k.Mod.Contains(tea.ModCtrl) && (k.Code == 'c' || k.Code == 'C' || strings.EqualFold(k.Text, "c")) {
		return true
	}
	return false
}

func introConfirmKeys(ks string, k tea.Key) bool {
	switch ks {
	case "enter", "return", "space":
		return true
	}
	switch k.Code {
	case tea.KeyEnter, tea.KeySpace, tea.KeyKpEnter: // KeyReturn aliases KeyEnter in ultraviolet
		return true
	default:
		return false
	}
}

func isLikelyBareModifier(k tea.Key) bool {
	switch k.Code {
	case tea.KeyLeftShift, tea.KeyRightShift,
		tea.KeyLeftAlt, tea.KeyRightAlt,
		tea.KeyLeftCtrl, tea.KeyRightCtrl,
		tea.KeyLeftSuper, tea.KeyRightSuper,
		tea.KeyLeftHyper, tea.KeyRightHyper,
		tea.KeyLeftMeta, tea.KeyRightMeta,
		tea.KeyCapsLock, tea.KeyNumLock, tea.KeyScrollLock:
		return true
	default:
		return false
	}
}

func introAdvances(ks string, k tea.Key) bool {
	if introConfirmKeys(ks, k) {
		return true
	}
	if isLikelyBareModifier(k) {
		return false
	}
	return true
}

func nightChoiceFromKey(ks string, k tea.Key) (int, bool) {
	switch ks {
	case "1":
		return 1, true
	case "2":
		return 2, true
	case "3":
		return 3, true
	}
	switch k.Code {
	case tea.KeyKp1:
		return 1, true
	case tea.KeyKp2:
		return 2, true
	case tea.KeyKp3:
		return 3, true
	}
	if len(k.Text) == 1 {
		switch k.Text[0] {
		case '1':
			return 1, true
		case '2':
			return 2, true
		case '3':
			return 3, true
		}
	}
	return 0, false
}

func endingHeadline(e ending.Ending) string {
	switch e {
	case ending.TheRelay:
		return "THE RELAY: You kept the network up long enough to help the hall."
	case ending.DarkFrequency:
		return "DARK FREQUENCY: You committed to Harrow's closed network."
	case ending.TheKidWasRight:
		return "THE KID WAS RIGHT: The clue chain reached its final stage."
	case ending.FullBroadcast:
		return "FULL BROADCAST: You released Osei's recording publicly."
	case ending.TheConvoy:
		return "THE CONVOY: You traded Maren's safety for extraction."
	case ending.DeadAir:
		return "GONE DARK: Fuel ran out before the network stabilized."
	case ending.Fallback:
		return "FALLBACK: No major ending condition was fully met."
	default:
		return "Run complete."
	}
}

func epilogueParagraph(e ending.Ending, s game.Session) string {
	dogLine := "The old dog sleeps by the stove, still here for morning."
	if !s.Dog.Alive {
		dogLine = "The bunk is empty where the dog used to sleep."
	}
	freqLine := ""
	if len(s.State.FrequenciesSeen) > 0 {
		freqLine = fmt.Sprintf(" You logged %d discovered side-frequency threads.", len(s.State.FrequenciesSeen))
	}
	switch e {
	case ending.TheRelay:
		return "You kept supporting Maren's hub through the late nights and held trust. The hall remained organized long enough to move people safely. " + dogLine + freqLine
	case ending.DarkFrequency:
		return "You prioritized Harrow's closed-network strategy for multiple nights. Wide-band traffic dropped, and your run ended focused on controlled channels over public broadcast. " + dogLine + freqLine
	case ending.TheKidWasRight:
		return "You followed the Kid's clues deeply enough to reach the hidden payoff path. The final pattern matched the warnings you tracked across the campaign. " + dogLine + freqLine
	case ending.FullBroadcast:
		return "You broadcast Osei's recording openly. Everyone got the same information at once, including groups you could not control. " + dogLine + freqLine
	case ending.TheConvoy:
		return "You accepted convoy extraction by giving up protected location data. The immediate outcome favored your escape, but it broke trust with Maren's side. " + dogLine + freqLine
	case ending.DeadAir:
		return "Power failed before the campaign reached a stable late-game state. Your run ended from resource collapse rather than a major faction decision. " + dogLine + freqLine
	case ending.Fallback:
		return "Your choices did not lock into a top-priority ending path. The run closed on mixed outcomes and incomplete threads. " + dogLine + freqLine
	default:
		return "The run ended outside named paths. " + dogLine + freqLine
	}
}

func tailStrings(lines []string, n int) []string {
	if n <= 0 || len(lines) == 0 {
		return nil
	}
	if len(lines) <= n {
		return append([]string(nil), lines...)
	}
	return append([]string(nil), lines[len(lines)-n:]...)
}

func (m *Model) View() tea.View {
	var b strings.Builder
	switch m.session.Phase {
	case game.PhaseIntro:
		if m.width < minTerminalWidth || m.height < minTerminalHeight {
			b.WriteString(fmt.Sprintf("Terminal too small (%dx%d). DEAD AIR requires at least %dx%d.\n", m.width, m.height, minTerminalWidth, minTerminalHeight))
			b.WriteString("Resize your terminal and reconnect.\n")
			break
		}
		b.WriteString("RELAY POST KESTREL\n\n")
		b.WriteString("You are Operator Seven at a remote radio post.\n")
		b.WriteString("Each night follows the same loop: pre-shift checks, receive calls, optional scan, incident, then logbook.\n\n")
		b.WriteString("Controls on shift:\n")
		b.WriteString("  intro: [t] tutorial, [enter] start campaign\n")
		b.WriteString("  pre-shift: [enter] open radio, [f] feed dog, [p]/[u] pin or unpin thread, [l] read logbook\n")
		b.WriteString("  receive: [1]/[2]/[3] choose response, [s] open scan\n")
		b.WriteString("  scan: arrows tune, [1]-[4] change band, [enter] lock signal, [r] return to receive\n")
		b.WriteString("  incident/logbook: [enter] continue, [w] write note, [n] next night\n")
	case game.PhaseTutorial:
		b.WriteString("TUTORIAL NIGHT 0\n\n")
		b.WriteString("Recorded note: Check generator, feed the dog, and confirm band before transmitting.\n\n")
		b.WriteString("Press [enter] to begin Night 1.\n")
	case game.PhaseNight:
		s := &m.session.State
		card := m.session.CurrentNight()
		b.WriteString(game.ActTitle(card.Act))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("SHIFT · NIGHT %d · MODE %s\n\n", s.Night, m.session.Mode.String()))
		b.WriteString("LAST ACTION\n")
		b.WriteString("────────────────────────────────────────\n")
		b.WriteString(m.session.LastAction + "\n\n")
		if tail := tailStrings(m.session.TxLog, txLogViewLines); len(tail) > 0 {
			b.WriteString("TRANSMISSION LOG (tail)\n")
			for _, line := range tail {
				b.WriteString("  │ ")
				b.WriteString(line)
				b.WriteString("\n")
			}
			b.WriteString("\n")
		}
		b.WriteString(fmt.Sprintf("DOG %s (%s)  POWER %s  FUEL CHECK %d\n\n", m.session.Dog.Name, dogState(m.session), m.session.Power.Gauge(), s.Fuel))
		switch m.session.Mode {
		case night.ModePreShift:
			b.WriteString("PRE-SHIFT CHECKS\n")
			b.WriteString("────────────────────────────────────────\n")
			b.WriteString(card.PreShift + "\n\n")
			if m.session.NeedsDogName && m.session.State.Night == 1 {
				b.WriteString("DOG NAME REQUIRED\n")
				b.WriteString("  [1] Scout   [2] Ash   [3] Bramble\n\n")
			}
			if len(m.session.Logbook) > 0 {
				b.WriteString("LOGBOOK ENTRY (latest)\n")
				b.WriteString("  " + m.session.Logbook[len(m.session.Logbook)-1] + "\n\n")
			}
			b.WriteString("Pinned threads: ")
			pins := m.session.Pins.Items()
			if len(pins) == 0 {
				b.WriteString("(none)\n")
			} else {
				b.WriteString(strings.Join(pins, ", ") + "\n")
			}
		case night.ModeReceive:
			b.WriteString(fmt.Sprintf("INCOMING — %s  (%s)\n", card.Source, card.Hash))
			b.WriteString("────────────────────────────────────────\n")
			b.WriteString(card.Quote + "\n\n")
			for i := range card.Choices {
				ch := card.Choices[i]
				b.WriteString(fmt.Sprintf("  [%d]  %s\n", i+1, ch.Text))
			}
		case night.ModeScan:
			b.WriteString("SCAN MODE\n")
			b.WriteString("────────────────────────────────────────\n")
			b.WriteString(m.session.Waterfall(2, 60) + "\n\n")
			if m.session.ScannedText != "" {
				b.WriteString("Signal lock: " + m.session.ScannedText + "\n")
			} else {
				b.WriteString("Signal lock: none\n")
			}
		case night.ModeIncident:
			b.WriteString("INCIDENT\n")
			b.WriteString("────────────────────────────────────────\n")
			b.WriteString(card.Incident + "\n\n")
			b.WriteString("Press [enter] to acknowledge and open the end-of-shift log.\n")
		case night.ModeLogbook:
			b.WriteString("END-OF-SHIFT LOG\n")
			b.WriteString("────────────────────────────────────────\n")
			b.WriteString("Press [w] to write a note, then [n] or [enter] for next night.\n")
		}
	case game.PhaseGameOver:
		e := m.session.Ending()
		b.WriteString("END OF RUN\n\n")
		b.WriteString(endingHeadline(e))
		b.WriteString("\n\n")
		b.WriteString(epilogueParagraph(e, *m.session))
		b.WriteString("\n")
	}
	framed := shopFrame(m.width, m.session.Phase, m.session.State.Night, m.session.State.Fuel, b.String())
	v := tea.NewView(framed)
	v.AltScreen = true
	return v
}

func (m *Model) String() string { return fmt.Sprintf("%T", m) }

func dogState(s *game.Session) string {
	if !s.Dog.Alive {
		return "missing"
	}
	if s.Dog.Hunger == 0 {
		return "fed"
	}
	if s.Dog.Hunger >= 3 {
		return "critical"
	}
	return "hungry"
}

func displayKey(ks string) string {
	if strings.TrimSpace(ks) == "" {
		return "unknown"
	}
	return ks
}
