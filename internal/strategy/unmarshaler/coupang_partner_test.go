package unmarshaler

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCoupangPartner(t *testing.T) {
	tt := []struct {
		name        string
		input       []byte
		want        []PartnerResp
		wantedError error
	}{
		{
			name:  "GIVEN valid JSON THEN return the expected struct",
			input: []byte(`[{"productId":1,"productUrl":"url1","productImage":"img1"},{"productId":2,"productUrl":"url2","productImage":"img2"}]`),
			want:  []PartnerResp{{ProductID: "1", ProductImage: "img1", ProductURL: "url1"}, {ProductID: "2", ProductImage: "img2", ProductURL: "url2"}},
		},
		{
			name:        "GIVEN invalid JSON THEN return an error",
			input:       []byte("invalid json"),
			wantedError: errors.New("invalid format"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &CoupangPartner{}
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
