package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *model) updateViewportContent() {
	if m.state == StateChooseFile {
		s := "Select file\n"
		for i, choice := range m.choices {
			if i == m.cursor {
				s += "> "
			} else {
				s += "  "
			}
			s += choice + "\n"
		}
		m.vp.SetContent(logo + "\n\n" + s)
		return
	}
	if m.state == StateAddUser {
		input := fmt.Sprintf(
			"Enter new user name\n\n%s\n\n%s",
			m.textInput.View(),
			"(esc to quit)",
		) + "\n"
		m.vp.SetContent(logo + "\n\n" + input)
		return
	}
	if m.state == StateDeleteUser {
		s := "Select user to delete:\n\n"
		for i, u := range m.users {
			if i == m.selectedUser {
				s += fmt.Sprintf("> %s\n", u.Name)
			} else {
				s += fmt.Sprintf("  %s\n", u.Name)
			}
		}
		s += "\n(esc to quit)"
		m.vp.SetContent(logo + "\n\n" + s)
		return
	}

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
		var winner string
		if m.winner > -1 && m.winner < len(m.users) {
			winner = m.users[m.winner].Name
		}
		ui = lipgloss.Place(m.w, m.h,
			lipgloss.Right, lipgloss.Top,
			winnerStyle.Render(winner),
		)

		rendered := lipgloss.NewStyle().Render(ui)
		linesWinner := strings.Split(rendered, "\n")
		trimedName := strings.TrimLeft(linesWinner[1], " ")

		for i := range lines {
			trimed := strings.TrimLeft(linesWinner[i], " ")
			rep := m.w - len(trimedName) - 12
			if rep < 0 {
				rep = 0
			}
			space := strings.Repeat(" ", rep)
			if i == 1 {
				space += " "
			}
			lines[i] = lines[i] + space + trimed
		}
		logoPart = strings.Join(lines, "\n")
	}

	var s strings.Builder

	isOneSpace := len(m.users) > (m.h-len(strings.Split(logoPart, "\n"))-4)/2
	isTwoColumns := len(m.users) > (m.h - len(strings.Split(logoPart, "\n")) - 4)
	if isTwoColumns {
		newWidth := (m.w / 2) - m.maxNameLen - 6 // some padding
		if newWidth < 1 {
			newWidth = 1
		}
		for i := range m.users {
			m.users[i].progress.Width = newWidth
		}
		isOneSpace = false
	}

	for i, u := range m.users {
		name := u.Name
		if len(name) > m.maxNameLen {
			name = name[:m.maxNameLen]
		}
		newLine := true
		if isTwoColumns {
			if i%2 == 0 {
				newLine = false
			}
		}
		if i < len(m.users) {
			s.WriteString(m.renderUser(u, newLine, isOneSpace))
		}
		if !newLine {
			s.WriteString("   ") // space between columns
			if u.Removed {
				s.WriteString(strings.Repeat(" ", u.progress.Width+1))
			}
		}
	}
	s.WriteString("\n")
	// s.WriteString("\nCurrent time: " + time.Now().Format("15:04:05"))

	rendered := textStyle.Width(m.w).Height(m.h).Render(s.String())
	rendered = logoPart + "\n" + rendered
	m.vp.SetContent(rendered)
}

func (m *model) renderUser(u User, newLine, isOneSpace bool) string {
	nextLine := ""
	if newLine {
		if isOneSpace {
			nextLine = "\n"
		} else {
			nextLine = "\n\n"
		}
	}
	name := u.Name
	if len(name) > m.maxNameLen {
		name = name[:m.maxNameLen]
	}
	removedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("45"))

	if u.Removed {
		rendered := removedStyle.Render(fmt.Sprintf(" %-*s", m.maxNameLen, name))
		rendered = strings.ReplaceAll(rendered, "\n", "")
		return rendered + nextLine
	}
	rendered := normalStyle.Render(fmt.Sprintf(" %-*s", m.maxNameLen, name))
	rendered = strings.ReplaceAll(rendered, "\n", "")
	rendered = fmt.Sprintf("%s %s%s", rendered, u.progress.View(), nextLine)
	return rendered
}
