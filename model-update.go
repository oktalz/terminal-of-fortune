package main

import (
	"math/rand"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oktalz/terminal-of-fortune/data"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			data.StateSet(data.StateFinished)
			return m, tea.Quit
		case "r":
			if data.State() != data.StateRunning {
				// curently its not safe to do it while running
				return m, nil
			}
			// randomise users in slice
			data.Process(func(users []data.User) {
				for i := len(users) - 1; i > 0; i-- {
					j := rand.Intn(i)
					users[i], users[j] = users[j], users[i]
					if data.Winner() == (i) {
						data.WinnerSet(j)
					} else if data.Winner() == (j) {
						data.WinnerSet(i)
					}
				}
			})
			return m, nil
		// case "delete":
		// 	if data.State() != data.StateRunning {
		// 		// curently its not safe to do it while running
		// 		return m, nil
		// 	}
		// 	// randomise users in slice
		// 	data.Process(func(users []data.User) {
		// 		// remove all users that have flag Removed true
		// 		slices.DeleteFunc(users, func(u data.User) bool {
		// 			return u.Removed
		// 		})
		// 	})
		// 	return m, nil
		case "s", "c", " ", "enter":
			if !m.raceStarted {
				m.raceStarted = true
				data.StateSet(data.StateRunning)
				return m, tickCmd()
			}
			if m.gameOver {

				if data.UsersLen() == 0 {
					return m, tea.Quit
				}

				m.recalculateMaxNameLen()

				// Reset game
				m.gameOver = false
				data.Process(func(users []data.User) {
					for i := range users {
						users[i].Percentage = 0
						users[i].Progress.SetPercent(0)
					}
				})
				return m, tickCmd()
			}
			data.StateSet(data.StateRunning)
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
		data.Process(func(users []data.User) {
			for i := range users {
				if users[i].Progress.Percent() == float64(users[i].Percentage)/100 {
					continue
				}
				cmds = append(cmds, users[i].Progress.SetPercent(float64(users[i].Percentage)/100))
			}
		})
		m.updateViewportContent()
		return m, tea.Batch(cmds...)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		data.Process(func(users []data.User) {
			for i := range users {
				newModel, newCmd := users[i].Progress.Update(msg)
				users[i].Progress = newModel.(progress.Model)
				cmd = newCmd
				cmds = append(cmds, cmd)
			}
		})
		m.updateViewportContent()
		return m, tea.Batch(cmds...)

	case winnerMsg:
		m.gameOver = true
		return m, nil

	default:
		return m, nil
	}
}
