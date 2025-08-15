package data

type state string

const (
	StateWaiting  state = "Waiting"
	StateRunning  state = "Running"
	StatePaused   state = "Paused"
	StateFinished state = "Finished"
)
