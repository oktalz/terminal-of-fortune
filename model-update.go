package main

import (
	"math/rand"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.state = StateFinished
			return m, tea.Quit
		case "r":
			if m.state != StateRunning {
				// curently its not safe to do it while running
				return m, nil
			}
			// randomise users in slice
			for i := len(m.users) - 1; i > 0; i-- {
				j := rand.Intn(i)
				m.users[i], m.users[j] = m.users[j], m.users[i]
				if m.winner == (i) {
					m.winner = j
				} else if m.winner == (j) {
					m.winner = i
				}
			}
			return m, nil
		// case "delete":
		// 	if m.state != StateRunning {
		// 		// curently its not safe to do it while running
		// 		return m, nil
		// 	}
		// 	// randomise users in slice
		// 	Process(func(users []User) {
		// 		// remove all users that have flag Removed true
		// 		slices.DeleteFunc(users, func(u User) bool {
		// 			return u.Removed
		// 		})
		// 	})
		// 	return m, nil
		case "s", "c", " ", "enter":
			if !m.raceStarted {
				m.raceStarted = true
				m.state = StateRunning
				return m, tickCmd()
			}
			if m.gameOver {

				if len(m.users) == 0 {
					return m, tea.Quit
				}

				m.recalculateMaxNameLen()

				// Reset game
				m.gameOver = false
				for i := range m.users {
					m.users[i].Percentage = 0
					m.users[i].progress.SetPercent(0)
				}
				return m, tickCmd()
			}
			m.state = StateRunning
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.updateSize(msg.Width, msg.Height)
		return m, nil
	case tickMsg:
		if m.gameOver {
			return m, nil
		}
		// Update progress bars
		var cmds []tea.Cmd
		if !m.gameOver {
			cmds = append(cmds, tickCmd())
		}
		updateUsers(m)
		for i := range m.users {
			if m.users[i].progress.Percent() == float64(m.users[i].Percentage)/100 {
				continue
			}
			cmds = append(cmds, m.users[i].progress.SetPercent(float64(m.users[i].Percentage)/100))
		}
		m.updateViewportContent()
		return m, tea.Batch(cmds...)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		for i := range m.users {
			newModel, newCmd := m.users[i].progress.Update(msg)
			m.users[i].progress = newModel.(progress.Model)
			cmd = newCmd
			cmds = append(cmds, cmd)
		}
		m.updateViewportContent()
		return m, tea.Batch(cmds...)

	case winnerMsg:
		m.gameOver = true
		return m, nil

	default:
		return m, nil
	}
}
