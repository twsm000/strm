package strm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestClosableStream(t *testing.T) {
	done := NewClosableStream()
	close(done)
	assert.Panics(t, func() {
		close(done)
	})
}

func TestDrainer(t *testing.T) {
	drain := Drainer[int]()
	for range drain {
		assert.FailNow(t, "for/range block cannot be reached")
	}
}

func TestGetValidDeadlineWhenPassingValidDeadline(t *testing.T) {
	var deadline StreamDeadline = make(chan EmptyStruct)
	validDeadline := GetValidDeadline(deadline)
	assert.Equal(t, deadline, validDeadline)
}

func TestGetValidDeadlineWhenPassingInvalidDeadline(t *testing.T) {
	var invalidDeadline StreamDeadline
	validDeadline := GetValidDeadline(invalidDeadline)
	assert.NotEqual(t, invalidDeadline, validDeadline)
}

func TestSendToWhenPassingValidSenderAndDeadline(t *testing.T) {
	quit := NewClosableStream()
	sender, receiver, free := NewStream[int]()

	const (
		expected   int = 10
		unexpected int = 20
	)

	go func() {
		assert.True(t, SendTo(sender, expected, quit))
		close(quit)
		assert.False(t, SendTo(sender, unexpected, quit))
		free()
	}()

	for value := range receiver {
		assert.Equal(t, expected, value)
	}
}

func TestSendToWhenPassingInvalidSender(t *testing.T) {
	quit := NewClosableStream()
	assert.False(t, SendTo(nil, 01, quit))
	close(quit)
}

func TestSendToWhenPassingInvalidSenderAndDeadline(t *testing.T) {
	assert.False(t, SendTo(nil, 01, nil))
}
