package main

import (
	"bufio"
	"math/rand"
	"os"
	"regexp"
	"time"

	"github.com/charmbracelet/bubbles/progress"
)

func readUsers(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m.users = append(m.users, User{Name: scanner.Text()})
		m.users[len(m.users)-1].progress = progress.New(progress.WithDefaultGradient())
		m.users[len(m.users)-1].progress.Init()
	}
	m.recalculateMaxNameLen()
	for i := range m.users {
		m.users[i].progress.Width = m.w - padding*2 - 4 - m.maxNameLen - 1
	}
	// randomise users
	for i := len(m.users) - 1; i > 0; i-- {
		j := rand.Intn(i)
		m.users[i], m.users[j] = m.users[j], m.users[i]
	}

	return scanner.Err()
}

func updateUsers(m *model) {
	switch m.state {
	case StateWaiting:
		animationDone := true
		dtStart := time.Now()
		for animationDone && time.Since(dtStart) < 1*time.Second {
			for i := range m.users {
				if m.users[i].progress.IsAnimating() {
					animationDone = false
					break
				}
			}
			time.Sleep(50 * time.Millisecond)
		}
		for i := range m.users {
			if i != m.winner {
				m.users[i].Percentage = 0
			}
		}
		m.state = StatePaused
	case StateRunning:
		if m.winner != -1 {
			time.Sleep(25 * time.Millisecond)
			for i := range m.users {
				if i != m.winner {
					m.users[i].Percentage = 0
				}
			}
			// remove winner
			// users = append(users[:m.winner], users[m.winner+1:]...)
			m.users[m.winner].Removed = true
			m.users[m.winner].Percentage = 0
			m.winner = -1
			m.lastUserUpdate = time.Now()
		}
		activeUsers := len(m.users)
		for i := range m.users {
			if m.users[i].Removed {
				activeUsers--
			}
		}
		if activeUsers == 0 {
			m.state = StatePaused
		}

		if time.Since(m.lastUserUpdate) < 125*time.Millisecond {
			return
		}
		m.lastUserUpdate = time.Now()

		for i := range m.users {
			if m.users[i].Removed {
				continue
			}
			// if m.users[i].progress.IsAnimating() {
			// 	continue
			// }
			if m.users[i].Percentage < 100 {
				max := 10
				m.users[i].Percentage += rand.Intn(max)
				if m.users[i].Percentage > 100 {
					m.users[i].Percentage = 100
				}
			}
		}
		// sort users based on current percentage
		// slices.SortFunc(m.users, func(a, b User) int {
		// 	viewB := b.progress.View()
		// 	viewA := a.progress.View()
		// 	var percentA, percentB int
		// 	match := re.FindStringSubmatch(viewA)
		// 	if len(match) > 1 {
		// 		percentA, _ = strconv.Atoi(match[1])
		// 	}
		// 	match = re.FindStringSubmatch(viewB)
		// 	if len(match) > 1 {
		// 		percentB, _ = strconv.Atoi(match[1])
		// 	}
		// 	return percentB - percentA
		// })
		var potentialwinners []int
		for i := range m.users {
			if m.users[i].Percentage >= 99 {
				potentialwinners = append(potentialwinners, i)
			}
		}
		if len(potentialwinners) > 0 {
			m.winner = potentialwinners[rand.Intn(len(potentialwinners))]
			m.users[m.winner].Percentage = 100
			m.state = StateWaitProgress
			// data.RandomiseUsers()
		}
		if len(potentialwinners) > 0 {
			// set other users to 99
			for i := range m.users {
				if i != m.winner && m.users[i].Percentage > 99 {
					m.users[i].Percentage = 99
				}
			}
		}
	case StateWaitProgress:
		uiUpdateDone := true
		for i := range m.users {
			if m.users[i].progress.IsAnimating() {
				uiUpdateDone = false
				break
			}
		}
		if uiUpdateDone {
			m.state = StateWaiting
		}
	}
}

var re = regexp.MustCompile(`(\d+)%`)

// func updateUsers() {
// 	for {
// 		m.winner := data.winner()
// 		st := data.State()
// 		if st == data.StateWaiting {
// 			time.Sleep(2 * time.Second)
// 			animationDone := true
// 			dtStart := time.Now()
// 			for animationDone && time.Since(dtStart) < 1*time.Second {
// 				for i := range users {
// 					if users[i].Progress.IsAnimating() {
// 						animationDone = false
// 						break
// 					}
// 				}
// 				time.Sleep(50 * time.Millisecond)
// 			}
// 			data.Process(func(users []data.User) {
// 				for i := range users {
// 					if i != m.winner {
// 						users[i].Percentage = 0
// 					}
// 				}
// 			})
// 			data.StateSet(data.StatePaused)
// 		}
// 		if st != data.StateRunning {
// 			time.Sleep(50 * time.Millisecond)
// 			continue
// 		}
// 		if m.winner != -1 {
// 			time.Sleep(25 * time.Millisecond)
// 			data.Process(func(users []data.User) {
// 				for i := range users {
// 					if i != m.winner {
// 						users[i].Percentage = 0
// 					}
// 				}
// 				// remove winner
// 				// users = append(users[:m.winner], users[m.winner+1:]...)
// 				users[m.winner].Removed = true
// 				users[m.winner].Percentage = 0
// 			})
// 			data.winnerSet(-1)
// 			m.winner = -1
// 		}
// 		data.Process(func(users []data.User) {
// 			for i := range users {
// 				if users[i].Removed {
// 					continue
// 				}
// 				if users[i].Percentage < 100 {
// 					users[i].Percentage += rand.Intn(3)
// 					if users[i].Percentage > 100 {
// 						users[i].Percentage = 100
// 					}
// 				}
// 			}
// 			var potentialwinners []int
// 			for i := range users {
// 				if users[i].Percentage >= 99 {
// 					potentialwinners = append(potentialwinners, i)
// 				}
// 			}
// 			if len(potentialwinners) > 0 {
// 				m.winner = potentialwinners[rand.Intn(len(potentialwinners))]
// 				users[m.winner].Percentage = 100
// 				data.winnerSet(m.winner)
// 				data.StateSet(data.StateWaiting)
// 				// data.RandomiseUsers()
// 			}
// 		})
// 		time.Sleep(25 * time.Millisecond)
// 	}
// }
