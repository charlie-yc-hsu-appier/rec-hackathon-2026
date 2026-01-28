package header

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdPopcornHeader(t *testing.T) {
	h := &AdpopcornHeader{UserAgent: "tzyu.net", ContentType: "application/json"}
	headers := h.GenerateHeaders(Params{})
	assert.Equal(t, map[string]string{"User-Agent": "tzyu.net", "Content-Type": "application/json"}, headers)
}
