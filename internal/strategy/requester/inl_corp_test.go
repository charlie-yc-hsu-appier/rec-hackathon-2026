package requester

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInlCorp(t *testing.T) {
	SizeCodeMap := map[string]string{
		"300x300":  "code-300",
		"1200x627": "code-1200-627",
	}
	tt := []struct {
		name        string
		SizeCodeMap map[string]string
		params      Params
		want        string
		wantErr     bool
	}{
		{
			name:        "GIVEN valid size and size_code macro THEN expect URL with correct size code",
			SizeCodeMap: SizeCodeMap,
			params: Params{
				RequestURL: "https://example.com/image?code={size_code}&user={user_id_lower}&click_id={click_id_base64}",
				UserID:     "TestUser",
				ImgWidth:   300,
				ImgHeight:  300,
				ClickID:    "test-id",
			},
			want: "https://example.com/image?code=code-300&user=testuser&click_id=dGVzdC1pZA",
		},
		{
			name:        "GIVEN invalid size and size_code macro THEN expect error",
			SizeCodeMap: SizeCodeMap,
			params: Params{
				RequestURL: "https://example.com/image?code={size_code}&user={user_id_lower}&click_id={click_id_base64}",
				UserID:     "TestUser",
				ImgWidth:   300,
				ImgHeight:  250,
				ClickID:    "test-id",
			},
			wantErr: true,
		},
		{
			name:        "GIVEN width and height macro THEN expect URL with correct size",
			SizeCodeMap: nil,
			params: Params{
				RequestURL: "https://example.com/image?size={width}x{height}&user={user_id_lower}&click_id={click_id_base64}",
				UserID:     "TestUser",
				ImgWidth:   300,
				ImgHeight:  300,
				ClickID:    "test-id",
			},
			want: "https://example.com/image?size=300x300&user=testuser&click_id=dGVzdC1pZA",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &InlCorp{SizeCodeMap: tc.SizeCodeMap}
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
