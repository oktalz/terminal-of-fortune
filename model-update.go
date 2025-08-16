package main

import (
	"math/rand"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == StateChooseFile {
			switch msg.String() {
			case "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter":
				readUsers(m.choices[m.cursor])
				m.state = StatePaused
			}
			return m, nil
		}
		if m.state == StateAddUser {
			switch msg.String() {
			case "enter":
				if m.textInput.Value() == "" {
					m.state = StatePaused
					return m, nil
				}
				m.users = append(m.users, User{Name: m.textInput.Value()})
				m.recalculateMaxNameLen()
				m.users[len(m.users)-1].progress = progress.New(progress.WithDefaultGradient())
				for i := range m.users {
					m.users[i].progress.Width = m.w - padding*2 - 4 - m.maxNameLen - 1
				}
				m.textInput.Reset()
				m.state = StatePaused
				return m, nil
			case "esc":
				m.state = StatePaused
				return m, nil
			}
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
		if m.state == StateDeleteUser {
			switch msg.String() {
			case "up":
				if m.selectedUser > 0 {
					m.selectedUser--
				}
			case "down":
				if m.selectedUser < len(m.users)-1 {
					m.selectedUser++
				}
			case "enter":
				m.users = append(m.users[:m.selectedUser], m.users[m.selectedUser+1:]...)
				m.recalculateMaxNameLen()
				for i := range m.users {
					m.users[i].progress.Width = m.w - padding*2 - 4 - m.maxNameLen - 1
				}
				m.winner = -1
				m.state = StatePaused
			case "esc":
				m.state = StatePaused
			}
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c", "q":
			m.state = StateExit
			m.exitTime = time.Now()
			m.h = m.h - 1
			// return m, tea.Quit
			cmd := func() tea.Msg {
				return tea.WindowSizeMsg{Width: m.w, Height: m.h}
			}
			return m, cmd
		case "r":
			if m.state != StatePaused {
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
		case "i", "insert":
			if m.state != StatePaused {
				// curently its not safe to do it while running
				return m, nil
			}
			m.state = StateAddUser
			return m, nil
		case "d", "delete":
			if m.state != StatePaused {
				// curently its not safe to do it while running
				return m, nil
			}
			m.state = StateDeleteUser
			m.selectedUser = 0
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
			newState := StateRunning
			if !m.raceStarted {
				m.raceStarted = true
				m.state = newState
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
			m.state = newState
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.updateSize(msg.Width, msg.Height)
		return m, nil
	case tickMsg:
		if m.state == StateExit {
			if time.Now().After(m.exitTime) {
				return m, tea.Quit
			} else {
				return m, tickCmd()
			}
		}
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
