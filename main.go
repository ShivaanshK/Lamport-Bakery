package main

import (
	"flag"
	"fmt"
	"lamport-bakery/bakery"
	"sync"
)

func main() {
	numProcesses := flag.Int("numProcesses", 2, "Number of processes for the bakery")
	numIterations := flag.Int("numIterations", 1, "Number of iterations to run")
	flag.Parse()
	numProcessesVal := *numProcesses
	numIterationsVal := *numIterations

	sharedMemory := &bakery.SharedMemory{MagicNum: 0}

	bakeryState := bakery.InitBakery(numProcessesVal)

	start := make(chan struct{}) // Channel to signal start
	var wg sync.WaitGroup
	for i := 0; i < numProcessesVal; i++ {
		wg.Add(1)
		go func(pid int) {
			defer wg.Done()
			<-start // Wait for the start signal
			for i := 0; i < numIterationsVal; i++ {
				bakeryState.BakeryAlgorithm(sharedMemory, pid)
			}
		}(i)
	}
	close(start) // Signal all goroutines to start
	wg.Wait()
	fmt.Printf("Final Magic Number: %v\n", sharedMemory.MagicNum)
}
