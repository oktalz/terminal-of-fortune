package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/oktalz/terminal-of-fortune/data"
)

func (m *model) View() string {
	currentTime := time.Now().Format("15:04:05")
	st := data.State()
	currentTime = " " + string(st)
	currentTime = " " // do not like it
	if !m.raceStarted {
		menu := "(s)tart (q)uit" + " " + currentTime
		border := lipgloss.NormalBorder()
		repeatCount := m.w - len(menu) - 5
		if repeatCount < 0 {
			repeatCount = 0
		}
		border.Bottom = strings.Repeat("─", repeatCount) + " " + menu + " " + "───"
		return lipgloss.NewStyle().Border(border).Render(m.vp.View())
	}

	if data.State() == data.StatePaused {
		menu := "(c)ontinue (q)uit" + " " + currentTime
		border := lipgloss.NormalBorder()
		repeatCount := m.w - len(menu) - 5
		if repeatCount < 0 {
			repeatCount = 0
		}
		border.Bottom = strings.Repeat("─", repeatCount) + " " + menu + " " + "───"
		return lipgloss.NewStyle().Border(border).Render(m.vp.View())
	}

	// if state.Load() == StateRunning {
	menu := "(q)uit" + " " + currentTime
	border := lipgloss.NormalBorder()
	repeatCount := m.w - len(menu) - 5
	if repeatCount < 0 {
		repeatCount = 0
	}
	border.Bottom = strings.Repeat("─", repeatCount) + " " + menu + " " + "───"
	return lipgloss.NewStyle().Border(border).Render(m.vp.View())
	//}

	// return borderStyle.Render(m.vp.View())
}

func (m *model) updateViewportContent() {
	if m.gameOver {
		m.vp.SetContent("")
		return
	}
	if data.State() == data.StatePaused && data.Winner() != -1 && false {
		winnerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(1, 2)

		var ui string
		data.Process(func(users []data.User) {
			ui = lipgloss.Place(m.w, m.h,
				lipgloss.Center, lipgloss.Center,
				winnerStyle.Render(users[data.Winner()].Name),
			)
		})

		rendered := lipgloss.NewStyle().Render(ui)
		m.vp.SetContent(rendered)
		return
	}

	var s strings.Builder
	data.Process(func(users []data.User) {
		for i, u := range users {
			name := u.Name
			if len(name) > m.maxNameLen {
				name = name[:m.maxNameLen]
			}
			if i < len(users) {
				removedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
				// normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("45"))

				if u.Removed {
					rendered := removedStyle.Render(fmt.Sprintf("%-*s", m.maxNameLen, name))
					rendered = strings.ReplaceAll(rendered, "\n", "")
					s.WriteString(rendered)
					s.WriteString("\n\n")
					continue
				}
				// s.WriteString(normalStyle.Render(fmt.Sprintf("%-*s\n", m.maxNameLen, name)))
				s.WriteString(fmt.Sprintf("%-*s %s\n\n", m.maxNameLen, name, users[i].Progress.View()))
				// s.WriteString(fmt.Sprintf("%s %s\n\n",
				// 	normalStyle.Render(fmt.Sprintf("%-*s", m.maxNameLen, name)),
				// 	users[i].Progress.View()))
			}
		}
	})
	s.WriteString("\n")
	// s.WriteString("\nCurrent time: " + time.Now().Format("15:04:05"))

	rendered := textStyle.Width(m.w).Height(m.h).Render(s.String())
	rendered = logo + "\n" + rendered
	m.vp.SetContent(rendered)
}
