package main

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m *model) View() string {
	currentTime := time.Now().Format("15:04:05")
	currentTime += " " + string(m.state)
	currentTime = " "
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

	if m.state == StatePaused {
		menu := "(c)ontinue (q)uit" + " " + currentTime
		border := lipgloss.NormalBorder()
		repeatCount := m.w - len(menu) - 5
		if repeatCount < 0 {
			repeatCount = 0
		}
		border.Bottom = strings.Repeat("─", repeatCount) + " " + menu + " " + "───"
		return lipgloss.NewStyle().Border(border).Render(m.vp.View())
	}

	menu := "(q)uit" + " " + currentTime
	border := lipgloss.NormalBorder()
	repeatCount := m.w - len(menu) - 5
	if repeatCount < 0 {
		repeatCount = 0
	}
	border.Bottom = strings.Repeat("─", repeatCount) + " " + menu + " " + "───"
	return lipgloss.NewStyle().Border(border).Render(m.vp.View())
}
