package main

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/gorilla/websocket"
)

type model struct {
	conn *websocket.Conn
	err  error
	ti   textinput.Model
	vp   viewport.Model
	text []string
}

func initialModel() model {
	ti := initTextInput()
	ti.Focus()

	vp := initViewPort()

	// Create and return the model
	return model{
		ti:   ti,
		vp:   vp,
		text: []string{},
	}
}

func (m model) Init() tea.Cmd {
	return connectCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Tea messages from connect
	case connectionMsg:
		if msg.err != nil {
			m.err = msg.err
			m.conn = nil
			errMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render("Connection failed: " + msg.err.Error())
			m.text = append(m.text, errMsg)
		} else {
			m.conn = msg.conn
			successMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("Connection successful!")
			m.text = append(m.text, successMsg)
		}
		// Refresh viewport content with updated text slice
		m.vp.SetContent(strings.Join(m.text, "\n"))

	case sentMsg:
		if msg.err != nil {
			errMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render("Message Failed: " + msg.err.Error())
			m.text = append(m.text, errMsg)
		} else {
			retMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Render("You") + msg.text
			m.text = append(m.text, retMsg)
		}
		m.vp.SetContent(strings.Join(m.text, "\n"))

	// Key presses
	case tea.KeyPressMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.conn == nil {
				errMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render("No connection with server, establishing connection")
				m.text = append(m.text, errMsg)
				m.vp.SetContent(strings.Join(m.text, "\n"))
				return m, connectCmd()
			}
			text := m.ti.Value()
			if strings.TrimSpace(text) == "" {
				return m, nil
			}
			m.ti.Reset()
			return m, sendMessageCmd(m.conn, text)
		}

	case tea.WindowSizeMsg:
		m.vp.SetWidth(msg.Width)
		m.vp.SetHeight(msg.Height - 3)
	}

	var tiCmd, vpCmd tea.Cmd

	m.ti, tiCmd = m.ti.Update(msg)

	m.vp, vpCmd = m.vp.Update(msg)

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() tea.View {
	// 1. Stack the viewport on top of the text input
	ui := lipgloss.JoinVertical(
		lipgloss.Left, // Align everything to the left
		m.vp.View(),   // The chat history (top)
		m.ti.View(),   // The text input (bottom)
	)

	// 2. Wrap it in a tea.View and return
	v := tea.NewView(ui)
	v.AltScreen = true // Keep it in the alternate terminal buffer (full screen)
	return v
}
