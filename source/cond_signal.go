package main

import (
	"fmt"
	"os"
	"sync"
)

func main() {
	var m sync.Mutex
	c := sync.NewCond(&m)
	n := 10
	running := make(chan bool, n)
	awake := make(chan int, n)
	for i := 0; i < n; i++ {
		go func(i int) {
			m.Lock()
			running <- true
			c.Wait()
			awake <- i
			m.Unlock()
		}(i)
		if i > 0 {
			a := <-awake
			if a != i-1 {
				fmt.Printf("wrong goroutine woke up: want %d, got %d\n", i-1, a)
				os.Exit(1)
			}
			fmt.Println("awake i:", i)
		}
		<-running
		m.Lock()
		c.Signal()
		m.Unlock()
	}
	fmt.Println("cond signal end")
}
