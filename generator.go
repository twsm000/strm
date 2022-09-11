package strm

func Generate[T any](quit StreamDeadline, produce func(StreamSender[T])) Transporter[T] {
	if produce == nil {
		return NewTransporter[T](nil, nil)
	}

	cargo := Confine(func(sender StreamSender[T]) {
		produce(sender)
	})

	return NewTransporter(cargo, quit)
}

func GenerateEnqueue[T any](queueSize int, quit StreamDeadline, produce func(StreamSender[T])) Transporter[T] {
	if produce == nil {
		return NewTransporter[T](nil, nil)
	}

	cargo := ConfineEnqueue(queueSize, func(sender StreamSender[T]) {
		produce(sender)
	})

	return NewTransporter(cargo, quit)
}
