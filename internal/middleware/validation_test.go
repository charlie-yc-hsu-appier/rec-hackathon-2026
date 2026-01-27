package middleware

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"
)

type mockValidatableRequest struct {
	field string
}

func (m *mockValidatableRequest) ValidateAll() error {
	if m.field == "" {
		return status.Errorf(codes.InvalidArgument, "field is required")
	}
	return nil
}

// mockValidatableRequestLegacy implements only Validate() method
type mockValidatableRequestLegacy struct {
	field string
}

func (m *mockValidatableRequestLegacy) Validate() error {
	if m.field == "" {
		return status.Errorf(codes.InvalidArgument, "field is required")
	}
	return nil
}

// mockNonValidatableRequest doesn't implement validation
type mockNonValidatableRequest struct {
	field string
}

func TestValidationUnaryInterceptor_WithValidateAll(t *testing.T) {
	interceptor := ValidationUnaryInterceptor
	handler := func(ctx context.Context, req any) (any, error) {
		return "success", nil
	}

	// Test valid request
	validReq := &mockValidatableRequest{field: "value"}
	resp, err := interceptor(context.Background(), validReq, nil, handler)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp)

	// Test invalid request
	invalidReq := &mockValidatableRequest{field: ""}
	resp, err = interceptor(context.Background(), invalidReq, nil, handler)
	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestValidationUnaryInterceptor_WithValidate(t *testing.T) {
	interceptor := ValidationUnaryInterceptor
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	// Test valid request
	validReq := &mockValidatableRequestLegacy{field: "value"}
	resp, err := interceptor(context.Background(), validReq, nil, handler)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp)

	// Test invalid request
	invalidReq := &mockValidatableRequestLegacy{field: ""}
	resp, err = interceptor(context.Background(), invalidReq, nil, handler)
	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestValidationUnaryInterceptor_NonValidatableRequest(t *testing.T) {
	interceptor := ValidationUnaryInterceptor
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	// Test non-validatable request (should pass through)
	req := &mockNonValidatableRequest{field: "value"}
	resp, err := interceptor(context.Background(), req, nil, handler)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp)
}
