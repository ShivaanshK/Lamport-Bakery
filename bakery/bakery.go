package bakery

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
