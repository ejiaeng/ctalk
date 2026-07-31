package main

import (
	tea "charm.land/bubbletea/v2"
)

type model struct {
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	v := tea.NewView("Barebones Bubble Tea boilerplate. Press q to quit.")
	v.AltScreen = true
	return v
}

