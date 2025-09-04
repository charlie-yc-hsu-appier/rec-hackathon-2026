package header

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdPopcornHeader(t *testing.T) {
	h := &AdPopcornHeader{UserAgent: "tzyu.net"}
	headers := h.GenerateHeaders(Params{})
	assert.Equal(t, map[string]string{"User-Agent": "tzyu.net"}, headers)
}
