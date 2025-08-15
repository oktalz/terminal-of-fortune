package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oktalz/terminal-of-fortune/data"
)

var textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("45"))

func main() {
	readFiles()

	data.WinnerSet(-1)
	data.StateSet(data.StatePaused)

	maxNameLen := 0
	data.Process(func(users []data.User) {
		for i, u := range users {
			if len(u.Name) > maxNameLen {
				maxNameLen = len(u.Name)
			}
			users[i].Progress = progress.New(progress.WithDefaultGradient())
			users[i].Progress.Init()
		}
	})

	if maxNameLen > 35 {
		maxNameLen = 35
	}

	m := &model{
		w:          1,
		h:          1,
		vp:         viewport.New(1, 1),
		maxNameLen: maxNameLen,
	}
	p := tea.NewProgram(m)
	go updateUsers()

	if _, err := p.Run(); err != nil {
		fmt.Print(err)
	}
}
