package strm

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBroadcastInvalid(t *testing.T) {
	const clients int = 1
	fleet, err := Broadcast[int](clients, nil)
	assert.Nil(t, fleet)
	assert.NotNil(t, err)
	assert.Errorf(t, err, broadcastErrMsgFmt, clients)
}

func TestBroadcastInvalidCargo(t *testing.T) {
	const clients int = 2
	fleet, err := Broadcast[int](clients, nil)
	assert.NotNil(t, fleet)
	assert.Nil(t, err)
	assert.Len(t, fleet, clients)

	wg := sync.WaitGroup{}
	wg.Add(len(fleet))
	for _, cargo := range fleet {
		go func(wg *sync.WaitGroup, cargo Transporter[int]) {
			defer wg.Done()
			var count int
			for range DownstreamFrom(cargo) {
				count++
			}
			assert.Equal(t, 0, count)
		}(&wg, cargo)
	}
	wg.Wait()
}

func TestBroadcastValidCargo(t *testing.T) {
	quit := NewClosableStream()
	defer close(quit)

	generator := func(sender StreamSender[int], quit StreamDeadline) {
		for i := 1; i <= 10; i++ {
			if !SendTo(sender, i, quit) {
				return
			}
		}
	}

	const clients int = 3
	fleet, err := Broadcast(clients, Generate(quit, generator))
	assert.NotNil(t, fleet)
	assert.Nil(t, err)
	assert.Len(t, fleet, clients)

	wg := sync.WaitGroup{}
	wg.Add(len(fleet))
	for _, cargo := range fleet {
		go func(wg *sync.WaitGroup, cargo Transporter[int]) {
			defer wg.Done()
			var expected int
			for got := range DownstreamFrom(cargo) {
				expected++
				assert.Equal(t, expected, got)
			}
		}(&wg, cargo)
	}
	wg.Wait()
}
