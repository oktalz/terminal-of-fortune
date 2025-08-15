package data

func Process(fn func(users []User)) {
	mu.Lock()
	defer mu.Unlock()
	fn(dataInstance.users)
}
