package strm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeat(t *testing.T) {
	until := NewClosableStream()
	defer close(until)

	values := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expectedValues := append(values, values...)
	for value := range DownstreamFrom(Repeat(until, values...)) {
		if len(expectedValues) == 0 {
			return
		}

		assert.Equal(t, expectedValues[0], value)
		expectedValues = expectedValues[1:]
	}
}

func TestRepeatWithInvalidValues(t *testing.T) {
	until := NewClosableStream()
	defer close(until)

	var values []string
	for range DownstreamFrom(Repeat(until, values...)) {
		assert.FailNow(t, "for range block cannot be reached")
	}
}
