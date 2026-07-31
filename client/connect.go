package main

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/gorilla/websocket"
)

type connectionMsg struct {
	conn *websocket.Conn
	err  error
}

const ADDRESS = "ws://localhost:8080/"

type closedMsg struct {
	err error
}

type sentMsg struct {
	err error
}

// connectCmd opens a websocket connection to the ADDRESS.
func connectCmd() tea.Cmd {
	return func() tea.Msg {
		dialer := websocket.Dialer{HandshakeTimeout: 5 * time.Second}
		conn, _, err := dialer.Dial(ADDRESS, nil)
		return connectionMsg{
			conn: conn,
			err:  err,
		}
	}
}

// closeCmd gracefully closes the websocket connection.
func closeCmd(conn *websocket.Conn) tea.Cmd {
	return func() tea.Msg {
		if conn == nil {
			return closedMsg{err: fmt.Errorf("no active connection")}
		}
		err := conn.Close()
		return closedMsg{err: err}
	}
}

// sendMessageCmd sends a text message via the websocket.
func sendMessageCmd(conn *websocket.Conn, text string) tea.Cmd {
	return func() tea.Msg {
		if conn == nil {
			return sentMsg{err: fmt.Errorf("no active connection")}
		}
		err := conn.WriteMessage(websocket.TextMessage, []byte(text))
		return sentMsg{err: err}
	}
}
