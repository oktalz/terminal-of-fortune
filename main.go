package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("45"))

var m *model

func main() {
	ti := textinput.New()
	ti.Placeholder = "New User"
	ti.Focus()
	ti.CharLimit = 15
	ti.Width = 20

	m = &model{
		w:         1,
		h:         1,
		vp:        viewport.New(1, 1),
		winner:    -1,
		state:     StatePaused,
		textInput: ti,
	}
	readFiles()

	maxNameLen := 0
	for i, u := range m.users {
		if len(u.Name) > maxNameLen {
			maxNameLen = len(u.Name)
		}
		m.users[i].progress = progress.New(progress.WithDefaultGradient())
		m.users[i].progress.Init()
	}

	if maxNameLen > 35 {
		maxNameLen = 35
	}
	m.maxNameLen = maxNameLen
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Print(err)
	}
	fmt.Println("└" + strings.Repeat("─", m.w-2) + "┘")
}
