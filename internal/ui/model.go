package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"kestrelpost/internal/game"
)

// Model is the root Bubble Tea model for one SSH session.
type Model struct {
	session *game.Session
}

func NewModel() *Model {
	return &Model{session: game.NewSession()}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
		if m.session.Phase == game.PhaseIntro {
			m.session.BeginFromIntro()
			return m, nil
		}
		if m.session.Phase == game.PhaseNight {
			switch msg.String() {
			case "1", "2", "3":
				m.session.ApplyChoice(int(msg.String()[0] - '0'))
				return m, nil
			}
		}
	}
	return m, nil
}

func (m *Model) View() tea.View {
	var b strings.Builder
	switch m.session.Phase {
	case game.PhaseIntro:
		b.WriteString("RELAY POST KESTREL — NORTHERN MANITOBA\n\n")
		b.WriteString("You are OPERATOR 7. Automation failed during the event.\n")
		b.WriteString("The HF rig is yours until the fuel is gone.\n\n")
		b.WriteString("Press any letter key to begin your shift.\n\n")
		b.WriteString("Then: 1–3 choose how to spend tonight’s bandwidth.\n")
		b.WriteString("Press q to disconnect.\n")
	case game.PhaseNight:
		s := &m.session.State
		b.WriteString(fmt.Sprintf("OPERATOR 7 · NIGHT %d · FUEL %d\n\n", s.Night, s.Fuel))
		b.WriteString("INCOMING — MAREN (relay hash 9f2c)\n")
		b.WriteString("\"We’re holding eight at the hall. Fever in two. What do I do?\"\n\n")
		b.WriteString("  [1] Long medical routing + reassurance (−22 fuel, +hub, +trust)\n")
		b.WriteString("  [2] Short factual packet (−10 fuel, +hub, +trust)\n")
		b.WriteString("  [3] Standby ping only (−3 fuel, −trust)\n\n")
		b.WriteString("Keys 1–3 commit. q quits.\n")
	case game.PhaseGameOver:
		b.WriteString("GENERATOR / BATTERY — END OF RUN\n\n")
		b.WriteString(fmt.Sprintf("Resolved ending: %s\n\n", m.session.Ending().String()))
		b.WriteString("Press q to disconnect.\n")
	}
	v := tea.NewView(b.String())
	v.AltScreen = true
	return v
}

func (m *Model) String() string { return fmt.Sprintf("%T", m) }
