package header

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdPopcornHeader(t *testing.T) {
	h := &AdpopcornHeader{UserAgent: "tzyu.net"}
	headers := h.GenerateHeaders(Params{})
	assert.Equal(t, map[string]string{"Content-Type": "application/json", "User-Agent": "tzyu.net"}, headers)
}
