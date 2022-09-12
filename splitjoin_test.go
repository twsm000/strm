package strm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitJoin(t *testing.T) {
	const queueSize = 100
	const totalSize = 100_000
	quit := NewClosableStream()
	defer close(quit)

	produce := func(sender StreamSender[int], quit StreamDeadline) {
		for i := 0; i < totalSize; i++ {
			if !SendTo(sender, i, quit) {
				return
			}
		}
	}
	isEven := func(v int) bool {
		return v%2 == 0
	}
	transfer := func(cargo Transporter[int]) Transporter[int] {
		return FilterEnqueue(queueSize, cargo, isEven)
	}

	var count int
	for value := range DownstreamFrom(SplitJoin(GenerateEnqueue(queueSize, quit, produce), transfer)) {
		assert.True(t, isEven(value))
		count++
	}

	assert.Equal(t, totalSize/2, count)
}
