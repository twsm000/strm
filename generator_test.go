package strm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	quit := NewClosableStream()
	defer close(quit)

	expected := []int{1, 2, 3, 4, 5}
	generatedValues := Generate(quit, func(sender StreamSender[int]) {
		for _, v := range expected {
			sender <- v
		}
	})

	got := []int{}
	for v := range DownstreamFrom(generatedValues) {
		got = append(got, v)
	}

	assert.ObjectsAreEqualValues(expected, got)
}

func TestGenerateEnqueue(t *testing.T) {
	quit := NewClosableStream()
	defer close(quit)

	expected := []int{1, 2, 3, 4, 5}
	generatedValues := GenerateEnqueue(5, quit, func(sender StreamSender[int]) {
		for loopCount := 1; loopCount <= 2; loopCount++ {
			for _, v := range expected {
				sender <- v
			}
		}
	})

	got := []int{}
	for v := range DownstreamFrom(generatedValues) {
		got = append(got, v)
	}

	assert.ObjectsAreEqualValues(append(expected, expected...), got)
}
