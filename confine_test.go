package strm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfine(t *testing.T) {
	expected := []int{1, 2, 3, 4, 5}
	stream := Confine(func(sender StreamSender[int]) {
		for _, v := range expected {
			sender <- v
		}
	})

	var got []int
	for value := range stream {
		got = append(got, value)
	}

	assert.ObjectsAreEqualValues(expected, got)
}

func TestConfineEnqueue(t *testing.T) {
	expected := []int{1, 2, 3, 4, 5}
	stream := ConfineEnqueue(5, func(sender StreamSender[int]) {
		for loopCount := 1; loopCount <= 2; loopCount++ {
			for _, v := range expected {
				output := v * loopCount
				if output == 6 {
					break
				}

				sender <- output
			}
		}
	})

	var got []int
	for value := range stream {
		got = append(got, value)
	}

	assert.ObjectsAreEqualValues(append(expected, 2, 4), got)
}
