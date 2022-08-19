package strm

import "math"

// CatchFirsts: return the first 'n' elements from the transported cargo
func CatchFirsts[T any](n int, cargo Transporter[T]) Transporter[T] {
	if n == 0 {
		return NewTransporter[T](nil, nil)
	}

	pkg, allDelivered := OpenCargo(cargo)

	n = int(math.Abs(float64(n)))
	newCargo := Confine(func(sender StreamSender[T]) {
		for value := range Downstream(pkg, allDelivered) {
			if !SendTo(sender, value, allDelivered) {
				return
			}

			n--
			if n == 0 {
				return
			}
		}
	})

	return NewTransporter(newCargo, allDelivered)
}
