package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type chooserModel struct {
	choices  []string
	cursor   int
	selected string
}

func (m chooserModel) Init() tea.Cmd {
	return nil
}

func (m chooserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.choices[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m chooserModel) View() string {
	s := "\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	s += "\n(press q to quit)"
	return s
}

func newChooser(choices []string) *chooserModel {
	return &chooserModel{
		choices: choices,
	}
}

func runChooser(choices []string) (string, error) {
	p := tea.NewProgram(newChooser(choices))
	m, err := p.Run()
	if err != nil {
		return "", err
	}
	if m, ok := m.(chooserModel); ok {
		return m.selected, nil
	}
	return "", nil
}
