package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("sub g end")
					return
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}

	ctx, cancel := context.WithCancel(context.Background())
	dstChan := gen(ctx)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ticker.C:
			cancel()
			break loop
		case a := <-dstChan:
			fmt.Println("a:", a)
		}
	}
	fmt.Println("ctx end:")
}
