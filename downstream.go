package strm

// Downstream: returns an valid receiver selected between the value or the deadline
// An invalid receiver will generate an already drained stream to be consumed.
func Downstream[T any](receiver StreamReceiver[T], exit StreamDeadline) StreamReceiver[T] {
	exit = GetValidDeadline(exit)
	if receiver == nil {
		receiver = Drainer[T]()
	}

	return Confine(func(sender StreamSender[T]) {
		for {
			select {
			case <-exit:
				return
			case value, ok := <-receiver:
				if !ok {
					return
				}
				select {
				case sender <- value:
				case <-exit:
					return
				}
			}
		}
	})
}
