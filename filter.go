package strm

func Filter[T any](cargo Transporter[T], isOk func(t T) bool) Transporter[T] {
	if isOk == nil {
		return NewTransporter[T](nil, nil)
	}

	pkg, allDelivered := OpenCargo(cargo)

	newCargo := Confine(func(sender StreamSender[T]) {
		for value := range Downstream(pkg, allDelivered) {
			if !isOk(value) {
				continue
			}

			if !SendTo(sender, value, allDelivered) {
				return
			}
		}
	})

	return NewTransporter(newCargo, allDelivered)
}

func FilterEnqueue[T any](queueSize int, cargo Transporter[T], isOk func(t T) bool) Transporter[T] {
	if isOk == nil {
		return NewTransporter[T](nil, nil)
	}

	pkg, allDelivered := OpenCargo(cargo)

	newCargo := ConfineEnqueue(queueSize, func(sender StreamSender[T]) {
		for value := range Downstream(pkg, allDelivered) {
			if !isOk(value) {
				continue
			}

			if !SendTo(sender, value, allDelivered) {
				return
			}
		}
	})

	return NewTransporter(newCargo, allDelivered)
}
