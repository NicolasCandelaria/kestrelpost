package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type Model struct{}

func NewModel() *Model { return &Model{} }

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *Model) View() tea.View {
	s := "KESTREL POST\n\nRelay automation offline.\n\nPress q to disconnect.\n"
	v := tea.NewView(s)
	v.AltScreen = true
	return v
}

func (m *Model) String() string { return fmt.Sprintf("%T", m) }
