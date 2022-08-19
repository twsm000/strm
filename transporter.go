package strm

// Transporter is any function that returns a stream receiver, deadline
// and a consumable status
type Transporter[T any] func() (StreamReceiver[T], StreamDeadline)

// NewTransporter: returns a Transporter that send the cargo from a place to another.
// An invalid cargo will return an transporter with a drained channel
func NewTransporter[T any](cargo StreamReceiver[T], endPoint StreamDeadline) Transporter[T] {
	endPoint = GetValidDeadline(endPoint)
	if cargo == nil {
		return func() (StreamReceiver[T], StreamDeadline) {
			return Drainer[T](), endPoint
		}
	}

	newCargo := Confine(func(sender StreamSender[T]) {
		for value := range Downstream(cargo, endPoint) {
			if !SendTo(sender, value, endPoint) {
				return
			}
		}
	})

	return func() (StreamReceiver[T], StreamDeadline) {
		return newCargo, endPoint
	}
}

// OpenCargo: return a valid stream receiver and deadline
func OpenCargo[T any](cargo Transporter[T]) (StreamReceiver[T], StreamDeadline) {
	if cargo == nil {
		return NewTransporter[T](nil, nil)()
	}

	return cargo()
}
