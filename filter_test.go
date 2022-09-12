package strm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	quit := NewClosableStream()
	defer close(quit)

	producer := func(sender StreamSender[int], quit StreamDeadline) {
		for i := 1; i < 100; i++ {
			sender <- i
		}
	}
	onlyEven := func(v int) bool {
		return v%2 == 0
	}

	for v := range DownstreamFrom(Filter(Generate(quit, producer), onlyEven)) {
		assert.True(t, onlyEven(v))
	}
}

func TestFilterEnqueue(t *testing.T) {
	quit := NewClosableStream()
	defer close(quit)

	producer := func(sender StreamSender[int], quit StreamDeadline) {
		for i := 1; i < 100; i++ {
			if !SendTo(sender, i, quit) {
				return
			}
		}
	}
	onlyEven := func(v int) bool {
		return v%2 == 0
	}

	for v := range DownstreamFrom(FilterEnqueue(100, GenerateEnqueue(100, quit, producer), onlyEven)) {
		assert.True(t, onlyEven(v))
	}
}
