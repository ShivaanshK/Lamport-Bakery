package bakery

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type SharedMemory struct {
	MagicNum int
}

type Bakery struct {
	Number   []int
	Choosing []bool
}

func InitBakery(numProcesses int) *Bakery {
	return &Bakery{
		Number:   make([]int, numProcesses),
		Choosing: make([]bool, numProcesses),
	}
}

func (bakery *Bakery) FindMax() (max int) {
	for _, num := range bakery.Number {
		if num > max {
			max = num
		}
	}
	return
}

func (bakeryState *Bakery) BakeryAlgorithm(sharedMemory *SharedMemory, pid int) {
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
			time.Sleep(5 * time.Millisecond)
		}
		// Wait until the other process has a higher number or the same number but higher PID
		for bakeryState.Number[i] != 0 && (bakeryState.Number[i] < bakeryState.Number[pid] || (bakeryState.Number[i] == bakeryState.Number[pid] && i < pid)) {
			time.Sleep(5 * time.Millisecond)
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
		rand.IntN(100)
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
