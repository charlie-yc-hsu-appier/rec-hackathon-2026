package requester

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	tt := []struct {
		name   string
		params Params
		want   string
	}{
		{
			name: "GIVEN valid parameters THEN return the expected URL",
			params: Params{
				RequestURL: "https://example.com/image?size={width}x{height}&user={user_id_lower}",
				UserID:     "TestUser",
				ImgWidth:   200,
				ImgHeight:  100,
			},
			want: "https://example.com/image?size=200x100&user=testuser",
		},
		{
			name: "GIVEN missing placeholders THEN return the expected URL",
			params: Params{
				RequestURL: "https://example.com/image/user/abc",
				UserID:     "User",
				ImgWidth:   50,
				ImgHeight:  50,
			},
			want: "https://example.com/image/user/abc",
		},
		{
			name: "GIVEN 300x300 size and size_code macro THEN expect URL with 300x300 size code",
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
			name: "GIVEN 1200x627 size and size_code macro THEN expect URL with 1200x627 size code",
			params: Params{
				RequestURL: "https://example.com/image?code={size_code}&user={user_id_lower}&click_id={click_id_base64}",
				UserID:     "TestUser",
				ImgWidth:   1200,
				ImgHeight:  627,
				ClickID:    "test-id",
			},
			want: "https://example.com/image?code=code-1200-627&user=testuser&click_id=dGVzdC1pZA",
		},
		{
			name: "GIVEN 1200x600 size and size_code macro THEN expect URL with 1200x600 size code",
			params: Params{
				RequestURL: "https://example.com/image?code={size_code}&user={user_id_lower}&click_id={click_id_base64}",
				UserID:     "TestUser",
				ImgWidth:   1200,
				ImgHeight:  600,
				ClickID:    "test-id",
			},
			want: "https://example.com/image?code=code-1200-600&user=testuser&click_id=dGVzdC1pZA",
		},
		{
			name: "GIVEN invalid size and size_code macro THEN expect URL with fallback 300x300 size code",
			params: Params{
				RequestURL: "https://example.com/image?code={size_code}&user={user_id_lower}&click_id={click_id_base64}",
				UserID:     "TestUser",
				ImgWidth:   300,
				ImgHeight:  250,
				ClickID:    "test-id",
			},
			want: "https://example.com/image?code=code-300&user=testuser&click_id=dGVzdC1pZA",
		},
		{
			name: "GIVEN 300x300 size and width height macro THEN expect URL with 300x300 size",
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
			strategy := &Default{SizeCodes: map[string]string{
				"300x300":  "code-300",
				"1200x627": "code-1200-627",
				"1200x600": "code-1200-600",
			}}
			got := strategy.GenerateRequestURL(tc.params)
			require.Equal(t, tc.want, got)
		})
	}
}
