package main

import "fmt"

func main() {
	//channerl boffered
	//here we have a number of channels defined
	c := make(chan int, 1)
	c <- 1
	fmt.Println(<-c)
}
