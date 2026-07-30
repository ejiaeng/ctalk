package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const ADDRESS = "localhost:8080"
const defaultPrompt = "┃ "

type sentMsg struct {
	err error
}

type connMsg struct {
    err error
}

type resetStatusMsg struct{}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
	}
}

type model struct {
	viewport    viewport.Model
	messages    []string
	textinput   textinput.Model
	senderStyle lipgloss.Style
	status      string
	err         error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Send a message..."
	ti.Focus()
	ti.Prompt = defaultPrompt
	ti.CharLimit = 280
	ti.SetWidth(30)

	vp := viewport.New(viewport.WithWidth(30), viewport.WithHeight(5))
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)
	vp.KeyMap.Left.SetEnabled(false)
	vp.KeyMap.Right.SetEnabled(false)

	return model{
		textinput:   ti,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		status:      "",
		err:         nil,
	}
}

func sendMessage(text string) tea.Cmd {
	return func() tea.Msg {
		conn, err := net.DialTimeout("tcp", ADDRESS, 5*time.Second)
		if err != nil {
			return sentMsg{err: err}
		}
		defer conn.Close()

		_, err = conn.Write([]byte(text))
		return sentMsg{err: err}
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width)
		m.textinput.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height - 1)

		if len(m.messages) > 0 {
			// Wrap content before setting it.
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()

	case sentMsg:
		if msg.err != nil {
			redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
			m.status = redStyle.Render(fmt.Sprintf("[%v]", msg.err))
			return m, tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
				return resetStatusMsg{}
			})
		}
		greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
		m.status = greenStyle.Render("Message sent")
		return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return resetStatusMsg{}
		})

	case resetStatusMsg:
		m.status = ""

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			val := m.textinput.Value()
			if val == "" {
				return m, nil
			}
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+val)
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")))
			m.textinput.Reset()
			m.viewport.GotoBottom()
			return m, sendMessage(val)
		default:
			// Send all other keypresses to the textinput.
			var cmd tea.Cmd
			m.textinput, cmd = m.textinput.Update(msg)
			return m, cmd
		}

	case cursor.BlinkMsg:
		// Textinput should also process cursor blinks.
		var cmd tea.Cmd
		m.textinput, cmd = m.textinput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() tea.View {
	viewportView := m.viewport.View()
	var content string
	if m.status != "" {
		content = viewportView + "\n" + m.status + "\n" + m.textinput.View()
	} else {
		content = viewportView + "\n" + m.textinput.View()
	}

	v := tea.NewView(content)
	c := m.textinput.Cursor()
	if c != nil {
		c.Y += lipgloss.Height(viewportView)
		if m.status != "" {
			c.Y += 1
		}
	}
	v.Cursor = c
	v.AltScreen = true
	return v
}