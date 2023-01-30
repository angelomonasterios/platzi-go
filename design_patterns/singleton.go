package main

import (
	"fmt"
	"sync"
	"time"
)

type Database struct {
}

var db *Database
var lock sync.Mutex

func (d Database) CreateSingleConnection() {
	fmt.Println("creating singleton for database")
	time.Sleep(2 * time.Second)
	fmt.Println("Creation done")
}
func getDatabase() *Database {
	lock.Lock()
	defer lock.Unlock()

	if db == nil {
		fmt.Println("creating DB connection...")
		db = &Database{}
		db.CreateSingleConnection()
	} else {
		fmt.Println("DB already Created")
	}
	return db
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			getDatabase()
		}()
	}
	wg.Wait()
}
