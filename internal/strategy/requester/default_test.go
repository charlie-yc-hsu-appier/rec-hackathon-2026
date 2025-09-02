package requester

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	sizeCodes := map[string]string{
		"300x300":  "code-300",
		"1200x627": "code-1200-627",
		"1200x600": "code-1200-600",
	}
	tt := []struct {
		name      string
		params    Params
		sizeCodes map[string]string
		want      string
		wantErr   bool
	}{
		{
			name: "GIVEN valid parameters THEN return the expected URL",
			params: Params{
				RequestURL: "https://example.com/image?size={width}x{height}&user={user_id_lower}",
				UserID:     "TestUser",
				ImgWidth:   200,
				ImgHeight:  100,
			},
			sizeCodes: nil,
			want:      "https://example.com/image?size=200x100&user=testuser",
		},
		{
			name: "GIVEN missing placeholders THEN return the expected URL",
			params: Params{
				RequestURL: "https://example.com/image/user/abc",
				UserID:     "User",
				ImgWidth:   50,
				ImgHeight:  50,
			},
			sizeCodes: nil,
			want:      "https://example.com/image/user/abc",
		},
		{
			name: "GIVEN valid parameters and size_code macro THEN expect URL with valid size code",
			params: Params{
				RequestURL: "https://example.com/image?code={size_code}&user={user_id_lower}&click_id={click_id_base64}",
				UserID:     "TestUser",
				ImgWidth:   300,
				ImgHeight:  300,
				ClickID:    "test-id",
			},
			sizeCodes: sizeCodes,
			want:      "https://example.com/image?code=code-300&user=testuser&click_id=dGVzdC1pZA",
		},
		{
			name: "GIVEN invalid size and size_code macro THEN expect error",
			params: Params{
				RequestURL: "https://example.com/image?code={size_code}&user={user_id_lower}&click_id={click_id_base64}",
				UserID:     "TestUser",
				ImgWidth:   300,
				ImgHeight:  250,
				ClickID:    "test-id",
			},
			sizeCodes: sizeCodes,
			wantErr:   true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &Default{SizeCodes: tc.sizeCodes}
			got, err := strategy.GenerateRequestURL(tc.params)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}
		})
	}
}
