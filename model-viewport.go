package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *model) updateViewportContent() {
	if m.gameOver {
		m.vp.SetContent("")
		return
	}

	logoPart := logo
	lines := strings.Split(logoPart, "\n")
	if m.state == StatePaused && m.winner != -1 {
		winnerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#045e9aff")).
			Padding(1, 3)

		var ui string
		ui = lipgloss.Place(m.w, m.h,
			lipgloss.Right, lipgloss.Top,
			winnerStyle.Render(m.users[m.winner].Name),
		)

		rendered := lipgloss.NewStyle().Render(ui)
		linesWinner := strings.Split(rendered, "\n")
		trimedName := strings.TrimLeft(linesWinner[1], " ")

		for i := range lines {
			trimed := strings.TrimLeft(linesWinner[i], " ")
			space := strings.Repeat(" ", m.w-len(trimedName)-12)
			if i == 1 {
				space += " "
			}
			lines[i] = lines[i] + space + trimed
		}
		logoPart = strings.Join(lines, "\n")
	}

	var s strings.Builder
	for i, u := range m.users {
		name := u.Name
		if len(name) > m.maxNameLen {
			name = name[:m.maxNameLen]
		}
		if i < len(m.users) {
			removedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
			// normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("45"))

			if u.Removed {
				rendered := removedStyle.Render(fmt.Sprintf(" %-*s", m.maxNameLen, name))
				rendered = strings.ReplaceAll(rendered, "\n", "")
				s.WriteString(rendered)
				s.WriteString("\n\n")
				continue
			}
			s.WriteString(fmt.Sprintf(" %-*s %s\n\n", m.maxNameLen, name, m.users[i].progress.View()))
		}
	}
	s.WriteString("\n")
	// s.WriteString("\nCurrent time: " + time.Now().Format("15:04:05"))

	rendered := textStyle.Width(m.w).Height(m.h).Render(s.String())
	rendered = logoPart + "\n" + rendered
	m.vp.SetContent(rendered)
}
