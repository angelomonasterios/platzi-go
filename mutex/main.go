package main

import (
	"fmt"
	"sync"
)

var (
	balance int = 100
)

func Deposit(amount int, wg *sync.WaitGroup, lock *sync.RWMutex) {
	defer wg.Done()
	lock.Lock()
	b := balance
	balance = b + amount
	lock.Unlock()

}
func Balance(locker *sync.RWMutex) int {
	locker.RLock()
	b := balance
	locker.RUnlock()
	return b
}

// solo 1 deposit -> write
// n Balance() -> read

func main() {
	var wg sync.WaitGroup
	var lock sync.RWMutex
	for i := 0; i <= 5; i++ {
		wg.Add(1)
		go Deposit(i*100, &wg, &lock)
	}
	wg.Wait()
	fmt.Println(Balance(&lock))
}
