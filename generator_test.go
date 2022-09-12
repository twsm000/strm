package strm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	quit := NewClosableStream()
	defer close(quit)

	expected := []int{1, 2, 3, 4, 5}
	producer := func(sender StreamSender[int], quit StreamDeadline) {
		for _, v := range expected {
			if !SendTo(sender, v, quit) {
				return
			}
		}
	}

	got := []int{}
	for v := range DownstreamFrom(Generate(quit, producer)) {
		got = append(got, v)
	}

	assert.ObjectsAreEqualValues(expected, got)
}

func TestGenerateEnqueue(t *testing.T) {
	quit := NewClosableStream()
	defer close(quit)

	expected := []int{1, 2, 3, 4, 5}
	producer := func(sender StreamSender[int], quit StreamDeadline) {
		for loopCount := 1; loopCount <= 2; loopCount++ {
			for _, v := range expected {
				if !SendTo(sender, v, quit) {
					return
				}
			}
		}
	}

	got := []int{}
	for v := range DownstreamFrom(GenerateEnqueue(5, quit, producer)) {
		got = append(got, v)
	}

	assert.ObjectsAreEqualValues(append(expected, expected...), got)
}
