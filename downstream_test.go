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

func TestDownstreamFromWhenPassInvalidTransporter(t *testing.T) {
	for range DownstreamFrom[int](nil) {
		assert.FailNow(t, "for range block cannot be reached")
	}
}

func TestDownstreamFromWhenPassValidTransporter(t *testing.T) {
	const expected int = 10
	quit := NewClosableStream()

	cargo := func() (StreamReceiver[int], StreamDeadline) {
		sender := make(chan int)
		go func() {
			sender <- expected
		}()
		return sender, quit
	}

	for got := range DownstreamFrom(cargo) {
		assert.Equal(t, expected, got)
		close(quit)
	}
}
