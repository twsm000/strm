package strm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	quit := NewClosableStream()
	defer close(quit)

	onlyEven := func(v int) bool {
		return v%2 == 0
	}
	gen := Generate(quit, func(sender StreamSender[int]) {
		for i := 1; i < 100; i++ {
			sender <- i
		}
	})

	for v := range DownstreamFrom(Filter(gen, onlyEven)) {
		assert.True(t, onlyEven(v))
	}
}

func TestFilterEnqueue(t *testing.T) {
	quit := NewClosableStream()
	defer close(quit)

	onlyEven := func(v int) bool {
		return v%2 == 0
	}
	gen := GenerateEnqueue(100, quit, func(sender StreamSender[int]) {
		for i := 1; i < 100; i++ {
			sender <- i
		}
	})

	for v := range DownstreamFrom(FilterEnqueue(100, gen, onlyEven)) {
		assert.True(t, onlyEven(v))
	}
}
