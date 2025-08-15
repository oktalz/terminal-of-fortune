package data

import (
	"math/rand"
)

func UsersLen() int {
	mu.Lock()
	defer mu.Unlock()
	return len(dataInstance.users)
}

func AddUser(user User) {
	mu.Lock()
	defer mu.Unlock()
	dataInstance.users = append(dataInstance.users, user)
}

func RandomiseUsers() {
	mu.Lock()
	defer mu.Unlock()
	// randomise users in slice
	for i := len(dataInstance.users) - 1; i > 0; i-- {
		j := rand.Intn(i)
		dataInstance.users[i], dataInstance.users[j] = dataInstance.users[j], dataInstance.users[i]
	}
}
