package main

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding = 2
)

type (
	tickMsg   time.Time
	winnerMsg struct{}
)

type User struct {
	Name       string
	Percentage int // 0 to 100
	Removed    bool
	progress   progress.Model
}

type model struct {
	w            int
	h            int
	vp           viewport.Model
	maxNameLen   int
	gameOver     bool
	raceStarted  bool
	users        []User
	winner       int
	state        state
	textInput    textinput.Model
	selectedUser int
	choices      []string
	cursor       int
}

func tickCmd() tea.Cmd {
	return tea.Tick(15*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *model) Init() tea.Cmd {
	return tickCmd()
}

func (m *model) recalculateMaxNameLen() {
	m.maxNameLen = 0
	for _, u := range m.users {
		if len(u.Name) > m.maxNameLen {
			m.maxNameLen = len(u.Name)
		}
	}
	if m.maxNameLen > 35 {
		m.maxNameLen = 35
	}
}

func (m *model) updateSize(w, h int) {
	m.w = w
	m.h = h
	for i := range m.users {
		m.users[i].progress.Width = w - padding*2 - 4 - m.maxNameLen - 1 // name + space
	}

	m.vp = viewport.New(w-2, h-2) // -2 for border
	m.vp.Style = lipgloss.NewStyle()

	m.updateViewportContent()
}
