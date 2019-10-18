package main

import (
	"fmt"
	"os"
	"sync"
)

func main() {
	var m sync.Mutex
	c := sync.NewCond(&m)
	n := 20

	running := make(chan int, n)
	awake := make(chan int, n)
	exit := false
	for i := 0; i < n; i++ {
		go func(g int) {
			m.Lock()
			for !exit {
				running <- g
				c.Wait()
				awake <- g
			}
			m.Unlock()
		}(i)
	}
	for i := 0; i < 3; i++ {
		fmt.Println("-------------------------------------")
		for i := 0; i < n; i++ {
			<-running // Will deadlock unless n are running.
		}
		if i == 2 {
			m.Lock()
			exit = true
			m.Unlock()
		}
		select {
		case <-awake:
			fmt.Println("goroutine not asleep")
			os.Exit(1)
		default:
		}
		m.Lock()
		c.Broadcast()
		m.Unlock()
		seen := make([]bool, n)
		for i := 0; i < n; i++ {
			g := <-awake
			fmt.Println("g is : ", g)
			if seen[g] {
				fmt.Println("goroutine woke up twice")
				os.Exit(1)
			}
			seen[g] = true
		}
	}
	select {
	case <-running:
		fmt.Println("goroutine did not exit")
		os.Exit(1)
	default:
	}
	c.Broadcast()
}
