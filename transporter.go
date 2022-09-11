package strm

// Transporter is any function that returns a stream receiver, deadline
// and a consumable status
type Transporter[T any] func() (StreamReceiver[T], StreamDeadline)

// Transfer is any function that receives a Transporter[T] and return another Transporter[T]
type Transfer[T any] func(t Transporter[T]) Transporter[T]

// NewTransporter: returns a Transporter that send the cargo from a place to another.
// An invalid cargo will return an transporter with a drained channel
func NewTransporter[T any](cargo StreamReceiver[T], endPoint StreamDeadline) Transporter[T] {
	if cargo == nil || endPoint == nil {
		return func() (StreamReceiver[T], StreamDeadline) {
			return Drainer[T](), GetValidDeadline(endPoint)
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
