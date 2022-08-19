package strm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTransporterWithInvalidCargoAndDeadline(t *testing.T) {
	transporter := NewTransporter[int](nil, nil)

	pkg, delivered := transporter()
	assert.NotNil(t, pkg)
	assert.NotNil(t, delivered)

	for range Downstream(pkg, delivered) {
		assert.FailNow(t, `
			content cannot be received because a transporter has nothing to delivery
		`)
	}
}

func TestNewTransporterWithInvalidDeadline(t *testing.T) {
	const expected int = 10
	sender := make(chan int, 1)
	sender <- expected

	transporter := NewTransporter(sender, nil)

	pkg, delivered := transporter()
	assert.NotNil(t, pkg)
	assert.NotNil(t, delivered)

	// Invalid cenario below. Instead use Downstream to consume the package content
	// transferred by the transporter or DownstreamFrom to consume the cargo directly
	// from the transporter
	var count int
	for content := range pkg {
		// for range block will never be executed
		assert.Equal(t, expected, content)
		count++
	}
	assert.Equal(t, 0, count)

	// Correct cenario
	for range Downstream(pkg, delivered) {
		assert.FailNow(t, `
			content cannot be received because a transporter without an endpoint
			should not have a destination to reach
		`)
	}
}
