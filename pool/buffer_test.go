package pool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Pool(t *testing.T) {
	b := GetBuffer()
	PutBuffer(b)
	assert.NotNil(t, BufferPool)
	// check we can put back nil
	PutBuffer(nil)

	MaxBufferSize = 10
	b = GetBuffer()
	b.WriteString("12345678901234567890")
	PutBuffer(b)
}
