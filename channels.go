package strm

type StreamSender[T any] chan<- T
type StreamReceiver[T any] <-chan T
type StreamFinisher func()
type StreamProcessor[T any] func(sender StreamSender[T])
type EmptyStruct struct{}
type StreamDeadline StreamReceiver[EmptyStruct]

func NewStream[T any]() (StreamSender[T], StreamReceiver[T], StreamFinisher) {
	return NewStreamBuffer[T](0)
}

func NewStreamBuffer[T any](queueSize int) (StreamSender[T], StreamReceiver[T], StreamFinisher) {
	stream := make(chan T, queueSize)
	return stream, stream, func() { close(stream) }
}

// NewClosableStream: return an channel that can be closed
func NewClosableStream() chan EmptyStruct {
	return make(chan EmptyStruct)
}

// Drainer: returns a channel already closed
func Drainer[T any]() StreamReceiver[T] {
	done := make(chan T)
	close(done)
	return done
}

// GetValidDeadline: check whether the deadline is assigned and return.
// When the deadline is not assigned a StreamDeadline drained will be returned.
func GetValidDeadline(deadline StreamDeadline) StreamDeadline {
	if deadline == nil {
		deadline = StreamDeadline(Drainer[EmptyStruct]())
	}

	return deadline
}

// SendTo: selects between send a value to the sender or quit.
// Return true only when a value is sent.
// Send to an invalid sender or invalid deadline will return false.
func SendTo[T any](sender StreamSender[T], value T, quit StreamDeadline) bool {
	quit = GetValidDeadline(quit)
	var invalidSender StreamDeadline
	if sender == nil {
		invalidSender = GetValidDeadline(invalidSender)
	}

	select {
	case <-quit:
		return false
	case <-invalidSender:
		return false
	case sender <- value:
		return true
	}
}
