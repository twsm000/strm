package strm

import (
	"fmt"
	"sync"
)

const broadcastErrMsgFmt = "broadcast allows sent only to 2 or more receivers(%d)"

// Broadcast replicates the transporter data to the clients
// An error will be returned when they are less than two clients
func Broadcast[T any](clients int, cargo Transporter[T]) ([]Transporter[T], error) {
	if clients <= 1 {
		return nil, fmt.Errorf(broadcastErrMsgFmt, clients)
	}
	senders := make([]StreamSender[T], clients)
	receivers := make([]Transporter[T], clients)
	finishers := make([]StreamFinisher, clients)

	pkg, quit := OpenCargo(cargo)
	for i := 0; i < len(receivers); i++ {
		sender, receiver, free := NewStream[T]()
		senders[i] = sender
		receivers[i] = NewTransporter(receiver, quit)
		finishers[i] = free
	}

	go func() {
		for _, free := range finishers {
			defer free()
		}

		broadcast := func(wg *sync.WaitGroup, sender StreamSender[T], value T, quit StreamDeadline) {
			defer wg.Done()
			SendTo(sender, value, quit)
		}

		wg := sync.WaitGroup{}
		for value := range Downstream(pkg, quit) {
			wg.Add(len(senders))
			for i := 0; i < len(senders); i++ {
				go broadcast(&wg, senders[i], value, quit)
			}
			wg.Wait()
		}
	}()

	return receivers, nil
}
