package main

import (
	"flag"
	"fmt"
	"lamport-bakery/bakery"
	"math/rand/v2"
	"sync"
)

type SharedMemory struct {
	MagicNum int
}

func main() {
	numProcesses := flag.Int("numProcesses", 2, "Number of processes for the bakery")
	numIterations := flag.Int("numIterations", 1, "Number of iterations to run")
	flag.Parse()
	numProcessesVal := *numProcesses
	numIterationsVal := *numIterations

	sharedMemory := &SharedMemory{MagicNum: 0}

	bakeryState := bakery.InitBakery(numProcessesVal)

	start := make(chan struct{}) // Channel to signal start
	var wg sync.WaitGroup
	for i := 0; i < numProcessesVal; i++ {
		wg.Add(1)
		go func(pid int) {
			defer wg.Done()
			<-start // Wait for the start signal
			for i := 0; i < numIterationsVal; i++ {
				bakeryAlgorithm(bakeryState, sharedMemory, pid)
			}
		}(i)
	}
	close(start) // Signal all goroutines to start
	wg.Wait()
	fmt.Printf("Final Magic Number: %v\n", sharedMemory.MagicNum)
}

func bakeryAlgorithm(bakeryState *bakery.Bakery, sharedMemory *SharedMemory, pid int) {
	// Run entry code
	num := entry(pid)

	// Choosing Phase
	bakeryState.Choosing[pid] = true
	bakeryState.Number[pid] = bakeryState.FindMax() + 1
	bakeryState.Choosing[pid] = false

	// Waiting Phase
	for i := 0; i < len(bakeryState.Number); i++ {
		if i == pid {
			continue
		}
		// Wait until the other process is done choosing
		for bakeryState.Choosing[i] {
			// Leave empty to force a busy wait
		}
		// Wait until the other process has a higher number or the same number but higher PID
		for bakeryState.Number[i] != 0 && (bakeryState.Number[i] < bakeryState.Number[pid] || (bakeryState.Number[i] == bakeryState.Number[pid] && i < pid)) {
			// Leave empty to force a busy wait
		}
	}

	// Enter critical section
	magicNum := criticalSection(sharedMemory, pid, num)
	// Run exit code
	exit(magicNum, pid)
	// Hand over lock
	bakeryState.Number[pid] = 0
	// Run rest of code
	remainder(pid)
}

// Preparation for entering the critical section
func entry(pid int) (num int) {
	fmt.Printf("Process %v entered entry portion\n", pid)
	numOperations := rand.IntN(100)
	for i := 0; i < numOperations; i++ {
		num += rand.IntN(100)
	}
	fmt.Printf("Process %v leaving entry portion\n", pid)
	return
}

// Executed immediately after critical section (before lock is given up)
func exit(magicNum, pid int) {
	fmt.Printf("Process %v entered exit portion\n", pid)
	fmt.Printf("Process %v has magic number of %v\n", pid, magicNum)
	fmt.Printf("Process %v leaving exit portion\n", pid)
}

// Simulate remainder of program
func remainder(pid int) {
	fmt.Printf("Process %v entered remainder portion\n", pid)
	numOperations := rand.IntN(10000)
	for i := 0; i < numOperations; i++ {
	}
	fmt.Printf("Process %v leaving remainder portion\n", pid)
}

func criticalSection(sharedMemory *SharedMemory, pid, num int) int {
	fmt.Printf("Process %v entered critical section\n", pid)
	switch num % 4 {
	case 0:
		sharedMemory.MagicNum = sharedMemory.MagicNum + num
	case 1:
		sharedMemory.MagicNum = sharedMemory.MagicNum - num
	case 2:
		sharedMemory.MagicNum = sharedMemory.MagicNum * num
	case 3:
		sharedMemory.MagicNum = sharedMemory.MagicNum / num
	}
	fmt.Printf("Process %v left critical section\n", pid)
	return sharedMemory.MagicNum
}
