package unmarshaler

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVendorPost_UnmarshalResponse(t *testing.T) {
	response := vendorPostResponse{
		RCode:    "abc",
		RMessage: "msg",
		Data: data{
			Result: []vendorPostResult{
				{
					ProductID:    123,
					ProductURL:   "http://url",
					ProductImage: "http://image",
				},
			},
		},
	}
	body, _ := json.Marshal(response)

	unmarshaler := &VendorPost{}
	products, err := unmarshaler.UnmarshalResponse(context.Background(), body)

	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, "123", products[0].ProductID)
	assert.Equal(t, "http://url", products[0].ProductURL)
	assert.Equal(t, "http://image", products[0].ProductImage)
}
