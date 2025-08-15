package main

import (
	"bufio"
	"math/rand"
	"os"
	"time"

	"github.com/oktalz/terminal-of-fortune/data"
)

func readUsers(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data.AddUser(data.User{Name: " " + scanner.Text()})
	}
	data.RandomiseUsers()

	return scanner.Err()
}

func updateUsers() {
	for {
		winnerLocal := data.Winner()
		st := data.State()
		if st == data.StateWaiting {
			time.Sleep(2 * time.Second)
			animationDone := true
			dtStart := time.Now()
			for animationDone && time.Since(dtStart) < 1*time.Second {
				data.Process(func(users []data.User) {
					for i := range users {
						if users[i].Progress.IsAnimating() {
							animationDone = false
							break
						}
					}
				})
				time.Sleep(50 * time.Millisecond)
			}
			data.Process(func(users []data.User) {
				for i := range users {
					if i != winnerLocal {
						users[i].Percentage = 0
					}
				}
			})
			data.StateSet(data.StatePaused)
		}
		if st != data.StateRunning {
			time.Sleep(50 * time.Millisecond)
			continue
		}
		if winnerLocal != -1 {
			time.Sleep(25 * time.Millisecond)
			data.Process(func(users []data.User) {
				for i := range users {
					if i != winnerLocal {
						users[i].Percentage = 0
					}
				}
				// remove winner
				// users = append(users[:winnerLocal], users[winnerLocal+1:]...)
				users[winnerLocal].Removed = true
				users[winnerLocal].Percentage = 0
			})
			data.WinnerSet(-1)
			winnerLocal = -1
		}
		data.Process(func(users []data.User) {
			for i := range users {
				if users[i].Removed {
					continue
				}
				if users[i].Percentage < 100 {
					users[i].Percentage += rand.Intn(3)
					if users[i].Percentage > 100 {
						users[i].Percentage = 100
					}
				}
			}
			var potentialWinners []int
			for i := range users {
				if users[i].Percentage >= 99 {
					potentialWinners = append(potentialWinners, i)
				}
			}
			if len(potentialWinners) > 0 {
				winnerLocal = potentialWinners[rand.Intn(len(potentialWinners))]
				users[winnerLocal].Percentage = 100
				data.WinnerSet(winnerLocal)
				data.StateSet(data.StateWaiting)
			}
		})
		time.Sleep(25 * time.Millisecond)
	}
}
