package data

import (
	"sync"
	"sync/atomic"

	"github.com/charmbracelet/bubbles/progress"
)

type User struct {
	Name       string
	Percentage int // 0 to 100
	Progress   progress.Model
	Removed    bool
}

type data struct {
	users  []User
	Winner atomic.Int32
	state  state
}

var (
	dataInstance = data{}
	mu           sync.Mutex
	muState      sync.Mutex
)

func State() state {
	muState.Lock()
	defer muState.Unlock()
	st := dataInstance.state
	return st
}

func StateSet(st state) {
	muState.Lock()
	defer muState.Unlock()
	dataInstance.state = st
}

func Winner() int {
	return int(dataInstance.Winner.Load())
}

func WinnerSet(winner int) {
	dataInstance.Winner.Store(int32(winner))
}
