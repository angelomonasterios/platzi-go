package main

import (
	"fmt"
	"sync"
	"time"
)

func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

type Memory struct {
	f     Function
	cache map[int]FuntionResult
	lock  sync.Mutex
}
type FuntionResult struct {
	value interface{}
	err   error
}
type Function func(key int) (interface{}, error)

func NewCache(f Function) *Memory {
	return &Memory{
		f:     f,
		cache: make(map[int]FuntionResult),
	}
}

func (m *Memory) Get(key int) (interface{}, error) {
	m.lock.Lock()

	result, exists := m.cache[key]

	m.lock.Unlock()

	if !exists {
		m.lock.Lock()
		result.value, result.err = m.f(key)
		m.cache[key] = result
		m.lock.Unlock()

	}
	return result.value, result.err
}

func GetFibonacci(n int) (interface{}, error) {
	return Fibonacci(n), nil
}

func main() {
	cache := NewCache(GetFibonacci)
	fibo := []int{42, 40, 40, 41, 42, 38, 40, 38}
	var wg sync.WaitGroup
	for _, n := range fibo {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			start := time.Now()
			result, err := cache.Get(index)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%d, %s, %d \n", index, time.Since(start), result)
		}(n)
		wg.Wait()
	}

}
