package rest

import (
	"context"
	"errors"
	"rec-vendor-api/internal/telemetry"
	"testing"
	"time"

	"bitbucket.org/plaxieappier/rec-go-kit/httpkit"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	gomock "go.uber.org/mock/gomock"
)

type clientTestSuite struct {
	suite.Suite

	mockHTTPClient *httpkit.MockClient
	client         Client
}

func (ts *clientTestSuite) SetupTest() {
	ts.mockHTTPClient = httpkit.NewMockClient(gomock.NewController(ts.T()))
	ts.client = NewClient(ts.mockHTTPClient)
}

func (ts *clientTestSuite) TestGet() {
	timeout := 5 * time.Second

	tests := []struct {
		name     string
		request  *httpkit.Request
		wantErr  bool
		mockFunc func()
	}{
		{
			name:    "GIVEN a valid request THEN return the response",
			request: httpkit.NewRequest("http://example.com").SetQueryParams(map[string]string{"key": "value"}),
			mockFunc: func() {
				expectedRequest := httpkit.NewRequest("http://example.com").
					SetQueryParams(map[string]string{"key": "value"}).
					SetMetrics(
						telemetry.Metrics.RestApiDurationSeconds.WithLabelValues("test-component"),
						telemetry.Metrics.RestApiErrorTotal.WithLabelValues("test-component"),
					)
				ts.mockHTTPClient.EXPECT().Get(gomock.Any(), expectedRequest, timeout, []int{200}, gomock.Any()).Return(&httpkit.Response{}, nil)
			},
		},
		{
			name:    "GIVEN mock error THEN return the error",
			request: httpkit.NewRequest("http://example.com"),
			wantErr: true,
			mockFunc: func() {
				expectedRequest := httpkit.NewRequest("http://example.com").
					SetMetrics(
						telemetry.Metrics.RestApiDurationSeconds.WithLabelValues("test-component"),
						telemetry.Metrics.RestApiErrorTotal.WithLabelValues("test-component"),
					)
				ts.mockHTTPClient.EXPECT().Get(gomock.Any(), expectedRequest, timeout, []int{200}, gomock.Any()).Return(nil, errors.New("mock error"))
			},
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			ctx := context.Background()
			_, err := ts.client.Get(ctx, tt.request, timeout, "test-component")
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestClientTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &clientTestSuite{})
}
