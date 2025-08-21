package unmarshaler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultUnmarshalStrategy(t *testing.T) {
	tt := []struct {
		name        string
		input       []byte
		want        *[]CoupangPartnerResp
		wantedError error
	}{
		{
			name:  "GIVEN valid JSON THEN return the expected struct",
			input: []byte(`[{"productId":1,"productUrl":"url1","productImage":"img1"},{"productId":2,"productUrl":"url2","productImage":"img2"}]`),
			want:  &[]CoupangPartnerResp{{ProductID: 1, ProductImage: "img1", ProductUrl: "url1"}, {ProductID: 2, ProductImage: "img2", ProductUrl: "url2"}},
		},
		{
			name:        "GIVEN invalid JSON THEN return an error",
			input:       []byte("invalid json"),
			wantedError: fmt.Errorf("invalid format. body: %v", "invalid json"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &DefaultUnmarshalStrategy{}
			got, err := strategy.UnmarshalResponse(tc.input)
			if tc.wantedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.wantedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}
		})
	}
}
