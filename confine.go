package strm

// Confine a stream process and automatically close the receiver
// after the job is done
func Confine[T any](process StreamProcessor[T]) StreamReceiver[T] {
	sender, receiver, free := NewStream[T]()

	go func() {
		defer free()
		process(sender)
	}()

	return receiver
}

// ConfineEnqueue confine a queue stream process and automatically close the receiver
// after the job is done
func ConfineEnqueue[T any](queueSize int, process StreamProcessor[T]) StreamReceiver[T] {
	sender, receiver, free := NewStreamBuffer[T](queueSize)

	go func() {
		defer free()
		process(sender)
	}()

	return receiver
}
