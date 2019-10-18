package main

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)

func main() {
	var g singleflight.Group

	c := make(chan string)
	var calls int32

	fn := func() (interface{}, error) {
		atomic.AddInt32(&calls, 1)
		return <-c, nil
	}

	const n = 10
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			v, err,isShared := g.Do("key", fn)
			if err != nil {
				fmt.Printf("Do error: %v\n", err)
				return
			}
			if v.(string) != "bar" {
				fmt.Printf("got %q; want %q\n", v, "bar")
				return
			}
			fmt.Println("v:",v, "isShared",isShared)
			wg.Done()
		}()
	}
	time.Sleep(100 * time.Millisecond) // let goroutines above block
	c <- "bar"
	wg.Wait()
	if got := atomic.LoadInt32(&calls); got != 1 {
		fmt.Printf("number of calls = %d; want 1", got)
		os.Exit(1)
	}
	fmt.Println("test end")
}
