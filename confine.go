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
