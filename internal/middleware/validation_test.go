package middleware

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ValidationTestSuite struct {
	suite.Suite
}

type mockValidatableRequest struct {
	field string
}

func (m *mockValidatableRequest) ValidateAll() error {
	if m.field == "" {
		return status.Errorf(codes.InvalidArgument, "field is required")
	}
	return nil
}

// mockValidatableRequestWithoutValidateAll implements only Validate() method
type mockValidatableRequestWithoutValidateAll struct {
	field string
}

func (m *mockValidatableRequestWithoutValidateAll) Validate() error {
	if m.field == "" {
		return status.Errorf(codes.InvalidArgument, "field is required")
	}
	return nil
}

// mockNonValidatableRequest doesn't implement validation
type mockNonValidatableRequest struct {
	field string
}

func (ts *ValidationTestSuite) TestValidationUnaryInterceptor_WithValidateAll() {
	tt := []struct {
		name     string
		request  *mockValidatableRequest
		wantCode codes.Code
		wantResp any
		wantErr  bool
	}{
		{
			name:     "GIVEN a valid request with ValidateAll THEN expect successful response",
			request:  &mockValidatableRequest{field: "value"},
			wantCode: codes.OK,
			wantResp: "success",
			wantErr:  false,
		},
		{
			name:     "GIVEN an invalid request with ValidateAll THEN expect invalid argument error",
			request:  &mockValidatableRequest{field: ""},
			wantCode: codes.InvalidArgument,
			wantResp: nil,
			wantErr:  true,
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			interceptor := ValidationUnaryInterceptor
			handler := func(ctx context.Context, req any) (any, error) {
				return "success", nil
			}

			resp, err := interceptor(context.Background(), tc.request, nil, handler)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, resp)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tc.wantCode, st.Code())
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tc.wantResp, resp)
			}
		})
	}
}

func (ts *ValidationTestSuite) TestValidationUnaryInterceptor_WithValidate() {
	tt := []struct {
		name     string
		request  *mockValidatableRequestWithoutValidateAll
		wantCode codes.Code
		wantResp any
		wantErr  bool
	}{
		{
			name:     "GIVEN a valid request with Validate THEN expect successful response",
			request:  &mockValidatableRequestWithoutValidateAll{field: "value"},
			wantCode: codes.OK,
			wantResp: "success",
			wantErr:  false,
		},
		{
			name:     "GIVEN an invalid request with Validate THEN expect invalid argument error",
			request:  &mockValidatableRequestWithoutValidateAll{field: ""},
			wantCode: codes.InvalidArgument,
			wantResp: nil,
			wantErr:  true,
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			interceptor := ValidationUnaryInterceptor
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return "success", nil
			}

			resp, err := interceptor(context.Background(), tc.request, nil, handler)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, resp)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tc.wantCode, st.Code())
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tc.wantResp, resp)
			}
		})
	}
}

func (ts *ValidationTestSuite) TestValidationUnaryInterceptor_NonValidatableRequest() {
	tt := []struct {
		name     string
		request  *mockNonValidatableRequest
		wantResp any
		wantErr  bool
	}{
		{
			name:     "GIVEN a non-validatable request THEN expect request to pass through",
			request:  &mockNonValidatableRequest{field: "value"},
			wantResp: "success",
			wantErr:  false,
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			interceptor := ValidationUnaryInterceptor
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return "success", nil
			}

			resp, err := interceptor(context.Background(), tc.request, nil, handler)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tc.wantResp, resp)
			}
		})
	}
}

func TestValidationTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ValidationTestSuite{})
}
