package main

import (
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oktalz/terminal-of-fortune/data"
)

const (
	padding = 2
)

type (
	tickMsg   time.Time
	winnerMsg struct{}
)

type model struct {
	w, h        int
	vp          viewport.Model
	maxNameLen  int
	gameOver    bool
	raceStarted bool
}

func tickCmd() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) recalculateMaxNameLen() {
	m.maxNameLen = 0
	data.Process(func(users []data.User) {
		for _, u := range users {
			if len(u.Name) > m.maxNameLen {
				m.maxNameLen = len(u.Name)
			}
		}
	})
	if m.maxNameLen > 15 {
		m.maxNameLen = 15
	}
}

func (m *model) updateSize(w, h int) {
	m.w = w
	m.h = h
	data.Process(func(users []data.User) {
		for i := range users {
			users[i].Progress.Width = w - padding*2 - 4 - m.maxNameLen - 1 // name + space
		}
	})

	m.vp = viewport.New(w-2, h-2) // -2 for border
	m.vp.Style = lipgloss.NewStyle()

	m.updateViewportContent()
}
