package requester

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultRequestURLStrategy(t *testing.T) {
	tt := []struct {
		name       string
		requestURL string
		userID     string
		imgWidth   int
		imgHeight  int
		want       string
	}{
		{
			name:       "GIVEN valid parameters THEN return the expected URL",
			requestURL: "https://example.com/image?size={width}x{height}&user={user_id}",
			userID:     "TestUser",
			imgWidth:   200,
			imgHeight:  100,
			want:       "https://example.com/image?size=200x100&user=testuser",
		},
		{
			name:       "GIVEN missing placeholders THEN return the expected URL",
			requestURL: "https://example.com/image/user/abc",
			userID:     "User",
			imgWidth:   50,
			imgHeight:  50,
			want:       "https://example.com/image/user/abc",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &DefaultRequestURLStrategy{}
			params := Params{
				RequestURL: tc.requestURL,
				UserID:     tc.userID,
				ImgWidth:   tc.imgWidth,
				ImgHeight:  tc.imgHeight,
			}
			got := strategy.GenerateRequestURL(params)
			require.Equal(t, tc.want, got)
		})
	}
}
