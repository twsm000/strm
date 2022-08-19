package strm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownstreamWhenPassInvalidReceiverAndDeadline(t *testing.T) {
	for range Downstream[int](nil, nil) {
		assert.FailNow(t, `for range cannot be reached`)
	}
}

func TestDownstreamWhenPassInvalidReceiver(t *testing.T) {
	quit := NewClosableStream()
	defer close(quit)
	for range Downstream[int](nil, quit) {
		assert.FailNow(t, `for range cannot be reached`)
	}
}

func TestDownstreamWhenPassInvalidDeadline(t *testing.T) {
	sender := make(chan int, 1)
	sender <- 1
	for range Downstream(sender, nil) {
		assert.FailNow(t, `for range cannot be reached`)
	}
}

func TestDownstream(t *testing.T) {
	quit := NewClosableStream()

	sender := make(chan int, 2)
	sender <- 1
	sender <- 2

	var expected int
	for got := range Downstream(sender, quit) {
		expected++
		assert.Equal(t, expected, got)
		if expected == 2 {
			close(quit)
		}
	}
}