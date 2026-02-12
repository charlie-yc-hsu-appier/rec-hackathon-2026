package unmarshaler

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostFormat_UnmarshalResponse(t *testing.T) {
	strategy := &PostFormat{}
	data, err := ioutil.ReadFile("post_format_test.json")
	assert.NoError(t, err)

	response, err := strategy.UnmarshalResponse(context.Background(), data)
	assert.NoError(t, err)

	assert.Len(t, response, 1)
	assert.Equal(t, "12345", response[0].ProductID)
	assert.Equal(t, "http://example.com/image.jpg", response[0].ProductImage)
	assert.Equal(t, "http://example.com/product", response[0].ProductURL)
}

func TestPostFormat_UnmarshalResponse_Error(t *testing.T) {
	strategy := &PostFormat{}
	data := []byte(`{"invalid_json"}`)

	_, err := strategy.UnmarshalResponse(context.Background(), data)
	assert.Error(t, err)
}
