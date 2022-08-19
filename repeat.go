package strm

// Repeat: send the values repeatedly through the transporter until reaching the deadline
func Repeat[T any](quit StreamDeadline, values ...T) Transporter[T] {
	if len(values) == 0 {
		return NewTransporter[T](nil, nil)
	}
	
	cargo := Confine(func(sender StreamSender[T]) {
		for {
			for _, value := range values {
				if !SendTo(sender, value, quit) {
					return
				}
			}
		}
	})

	return NewTransporter(cargo, quit)
}
