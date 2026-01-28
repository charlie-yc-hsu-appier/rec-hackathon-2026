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

// mockRequestType represents the type of mock request to create
type mockRequestType int

const (
	mockTypeValidateAll mockRequestType = iota
	mockTypeValidate
	mockTypeNonValidatable
)

func newMockRequest(mockType mockRequestType, fieldValue string) any {
	switch mockType {
	case mockTypeValidateAll:
		return &mockValidatableRequest{field: fieldValue}
	case mockTypeValidate:
		return &mockValidatableRequestWithoutValidateAll{field: fieldValue}
	case mockTypeNonValidatable:
		return &mockNonValidatableRequest{field: fieldValue}
	default:
		return nil
	}
}

func (ts *ValidationTestSuite) TestValidationUnaryInterceptor() {
	tt := []struct {
		name     string
		mockType mockRequestType
		fieldVal string
		wantCode codes.Code
		wantResp any
		wantErr  bool
	}{
		{
			name:     "GIVEN a valid request with ValidateAll THEN expect successful response",
			mockType: mockTypeValidateAll,
			fieldVal: "value",
			wantCode: codes.OK,
			wantResp: "success",
			wantErr:  false,
		},
		{
			name:     "GIVEN an invalid request with ValidateAll THEN expect invalid argument error",
			mockType: mockTypeValidateAll,
			fieldVal: "",
			wantCode: codes.InvalidArgument,
			wantResp: nil,
			wantErr:  true,
		},
		{
			name:     "GIVEN a valid request with Validate THEN expect successful response",
			mockType: mockTypeValidate,
			fieldVal: "value",
			wantCode: codes.OK,
			wantResp: "success",
			wantErr:  false,
		},
		{
			name:     "GIVEN an invalid request with Validate THEN expect invalid argument error",
			mockType: mockTypeValidate,
			fieldVal: "",
			wantCode: codes.InvalidArgument,
			wantResp: nil,
			wantErr:  true,
		},
		{
			name:     "GIVEN a non-validatable request THEN expect request to pass through",
			mockType: mockTypeNonValidatable,
			fieldVal: "value",
			wantCode: codes.OK,
			wantResp: "success",
			wantErr:  false,
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			interceptor := ValidationUnaryInterceptor
			handler := func(ctx context.Context, req any) (any, error) {
				return "success", nil
			}

			request := newMockRequest(tc.mockType, tc.fieldVal)
			resp, err := interceptor(context.Background(), request, nil, handler)

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

func TestValidationTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ValidationTestSuite{})
}
