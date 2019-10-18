package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	dataCh := make(chan int, 100)
	stopCh := make(chan struct{})

	// senders
	for i := 0; i < 5; i++ {
		go func() {
			for {
				select {
				case <-stopCh:
					fmt.Println("stop")
					return
				case dataCh <- rand.Intn(5):
				}
			}
		}()
	}
	for value := range dataCh {
		if value == 1 {
			fmt.Println("send stop signal to senders.")
			close(stopCh)
			break
		}
		fmt.Println(value)
	}
	time.Sleep(5 * time.Second)
}
