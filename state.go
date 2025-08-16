package main

type state string

const (
	StateWaiting      state = "Waiting"
	StateWaitProgress state = "WaitingProgress"
	StateRunning      state = "Running"
	StatePaused       state = "Paused"
	StateExit         state = "Exit"
	StateAddUser      state = "AddUser"
	StateDeleteUser   state = "DeleteUser"
	StateChooseFile   state = "ChooseFile"
)
