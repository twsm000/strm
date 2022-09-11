package strm

import (
	"runtime"
	"sync"
)

// SplitJoin split the work into the number of CPUs and return join the workers result
// into a single Transporter
func SplitJoin[T any](cargo Transporter[T], transfer Transfer[T]) Transporter[T] {
	return SplitJoinInto(0, cargo, transfer)
}

// SplitJoinInto split the work into the amount of workers and return join the workers result
// into a single Transporter
// Workers lower than two will be replaced by the number of CPU cores
func SplitJoinInto[T any](workers int, cargo Transporter[T], transfer Transfer[T]) Transporter[T] {
	if transfer == nil {
		return NewTransporter[T](nil, nil)
	}

	if workers < 2 {
		workers = runtime.NumCPU()
	}

	_, allDelivered := OpenCargo(cargo)
	newCargo := Confine(func(sender StreamSender[T]) {
		wg := sync.WaitGroup{}
		wg.Add(workers)

		for index := 0; index < workers; index++ {
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				for value := range DownstreamFrom(transfer(cargo)) {
					if !SendTo(sender, value, allDelivered) {
						return
					}
				}
			}(&wg)
		}

		wg.Wait()
	})

	return NewTransporter(newCargo, allDelivered)
}
