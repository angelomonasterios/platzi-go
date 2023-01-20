package main

import "time"

type job struct {
	Name   string
	Delay  time.Duration
	Number int
}

type Worker struct {
	Id         int
	JobQueue   chan job
	WorkerPool chan chan job
	QuitChan   chan bool
}
type Dispatcher struct {
	WorkerPool chan chan job
	MaxWorkers int
	JobQueue   chan job
}

func NewWorker(id int, workerPool chan chan job) *Worker {
	return &Worker{
		Id:         id,
		JobQueue:   make(chan job),
		WorkerPool: workerPool,
		QuitChan:   make(chan bool),
	}

}

func main() {

}
