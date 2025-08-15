package main

type state string

const (
	StateWaiting      state = "Waiting"
	StateWaitProgress state = "WaitingProgress"
	StateRunning      state = "Running"
	StatePaused       state = "Paused"
	StateFinished     state = "Finished"
)
