package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	workerPool()
}

func workerPool() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	wg := &sync.WaitGroup{}
	toComplete, completedNumbers := make(chan int, 5), make(chan int, 5)

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ctx, toComplete, completedNumbers)
		}()
	}

	go func() {
		for i := 0; i < 50; i++ {
			toComplete <- i
		}
		close(toComplete)
	}()

	go func() {
		defer close(completedNumbers)
		wg.Wait()
	}()

	counter := 0
	for num := range completedNumbers {
		fmt.Println(num)
		counter += 1
		fmt.Println(counter)
	}
}

func worker(ctx context.Context, toComplete <-chan int, completedNumbers chan<- int) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		for num := range toComplete {
			completedNumbers <- num * num
		}
		return
	}
}
