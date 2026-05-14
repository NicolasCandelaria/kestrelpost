package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"kestrelpost/internal/ending"
	"kestrelpost/internal/game"
)

const txLogViewLines = 6

// Model is the root Bubble Tea model for one SSH session.
type Model struct {
	session *game.Session
	width   int // from WindowSizeMsg; chrome uses effectiveWidth(width)
}

func NewModel() *Model {
	return &Model{session: game.NewSession(), width: 80}
}

func (m *Model) Init() tea.Cmd {
	return tea.RequestWindowSize
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Bubble Tea v2: presses are [tea.KeyPressMsg], releases are [tea.KeyReleaseMsg].
	// Some SSH/terminal paths may surface other [tea.KeyMsg] implementations.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width > 0 {
			m.width = msg.Width
		}
		return m, nil
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
		if introAdvances(ks, k) {
			m.session.BeginFromIntro()
		}
		return m, nil
	case game.PhaseNight:
		if choice, ok := nightChoiceFromKey(ks, k); ok {
			m.session.ApplyChoice(choice)
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
		return "The relay held — long enough for the hall to believe morning."
	case ending.DarkFrequency:
		return "You went tight. The wide band went quiet without forgiving anyone."
	case ending.TheKidWasRight:
		return "The kid’s thread closed on something that wasn’t a voice."
	case ending.FullBroadcast:
		return "The loop left your tower and kept running without your hand on it."
	case ending.TheConvoy:
		return "You traded the map for a seat. The road remembers both."
	case ending.DeadAir:
		return "The generator quit before the story did."
	case ending.Fallback:
		return "The log ended in a verdict the night wouldn’t name."
	default:
		return "Closure."
	}
}

func epilogueParagraph(e ending.Ending) string {
	switch e {
	case ending.TheRelay:
		return "The hall stayed lit. Maren’s voice thinned to static, then steadied as the convoy hash matched yours. You stayed on the key until the window closed—not heroics, just timing. The north held one more night because someone kept answering when answering still mattered."
	case ending.DarkFrequency:
		return "Harrow’s carrier outlasted your denial. What you thought was noise lined up into steps—someone else walking the same band plan. You finished the shift with clean hands and a dirty frequency; the next op inherits the hum you wouldn’t name."
	case ending.TheKidWasRight:
		return "The kid’s tally wasn’t gossip; it was triangulation. When the proof landed, it wasn’t loud—just inevitable, like a fuse counting down in someone else’s pocket. You kept the chair warm; they took the story out the door."
	case ending.FullBroadcast:
		return "Osei’s packet wasn’t a leak—it was a handoff. The full release washed the board clean and drowned nuance in signal. You did what the procedure demanded; the air belongs to everyone now, including the parts nobody’s ready to hear."
	case ending.TheConvoy:
		return "The convoy key lied, and the log caught the seam too late. Trust didn’t break in one message—it sheared along an old fault you’d been papering over with procedure. You signed off correct; the road still went wrong."
	case ending.DeadAir:
		return "The last frames are breath and backoff. Maren’s side went dark while you were still reaching for the next prefix. The ending isn’t moral—it’s mechanical: the machine ran out of room to be kind."
	case ending.Fallback:
		return "The board closed on a verdict that fit the numbers more than the night. Nothing dramatic in the log—just drift, compromise, and the ordinary way a relay stops being yours. You shut it down; the story keeps going without a headline."
	default:
		return "The run ended outside the labeled paths—telemetry intact, narrative thin. File it under operator variance: the machine got an answer even if the myth didn’t."
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
		b.WriteString("RELAY POST KESTREL\n\n")
		b.WriteString("Northern Manitoba. A repeater tower that used to run itself—until the chain broke mid-event and left you holding the key.\n\n")
		b.WriteString("You are OPERATOR 7. On paper you’re redundancy. In practice you’re the voice people find when the usual nets go strange: Maren at the hall, strangers on scan, ")
		b.WriteString("sometimes a second operator who signs like a colleague and argues like a strategist.\n\n")
		b.WriteString("Each night the band fills with requests for truth, time, and silence you can’t give everyone. Listen, answer, or refuse—then live in what that costs.\n\n")
		b.WriteString("Use the footer keys when you’re ready. Bare modifier taps won’t start a shift.\n")
	case game.PhaseNight:
		s := &m.session.State
		card := m.session.CurrentNightCard()
		b.WriteString(game.ActTitle(card.Act))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("SHIFT · NIGHT %d\n\n", s.Night))
		if tail := tailStrings(m.session.TxLog, txLogViewLines); len(tail) > 0 {
			b.WriteString("TRANSMISSION LOG (tail)\n")
			for _, line := range tail {
				b.WriteString("  │ ")
				b.WriteString(line)
				b.WriteString("\n")
			}
			b.WriteString("\n")
		}
		b.WriteString(fmt.Sprintf("INCOMING — %s  (%s)\n", card.Source, card.Hash))
		b.WriteString("────────────────────────────────────────\n")
		b.WriteString(card.Quote)
		b.WriteString("\n\n")
		for i := range card.Choices {
			ch := card.Choices[i]
			b.WriteString(fmt.Sprintf("  [%d]  %s\n", i+1, ch.Reply))
		}
	case game.PhaseGameOver:
		e := m.session.Ending()
		b.WriteString("END OF RUN\n\n")
		b.WriteString(endingHeadline(e))
		b.WriteString("\n\n")
		b.WriteString(epilogueParagraph(e))
		b.WriteString("\n")
	}
	framed := shopFrame(m.width, m.session.Phase, m.session.State.Night, m.session.State.Fuel, b.String())
	v := tea.NewView(framed)
	v.AltScreen = true
	return v
}

func (m *Model) String() string { return fmt.Sprintf("%T", m) }
