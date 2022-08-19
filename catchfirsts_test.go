package strm

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCatchFirsts(t *testing.T) {
	quit := NewClosableStream()
	defer close(quit)

	cargo := Confine(func(sender StreamSender[uint64]) {
		rand.Seed(time.Now().Unix())
		for SendTo(sender, rand.Uint64(), quit) {
		}
	})

	const expected int = 20
	var got int
	for range DownstreamFrom(CatchFirsts(expected, NewTransporter(cargo, quit))) {
		got++
	}
	assert.Equal(t, expected, got)
}
