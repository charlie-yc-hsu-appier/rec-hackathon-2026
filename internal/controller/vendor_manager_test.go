package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"rec-vendor-api/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type VendorsTestSuite struct {
	suite.Suite
}

func (ts *VendorsTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (ts *VendorsTestSuite) TestGetVendors() {
	tt := []struct {
		name         string
		vendorConfig config.VendorConfig
		wantBody     string
	}{
		{
			name: "GIVEN valid vendor config THEN expect response with all vendors",
			vendorConfig: config.VendorConfig{
				Vendors: []config.Vendor{
					{
						Name: "vendor1",
						Request: config.URLPattern{
							URL: "https://api.vendor1.com/recommend?key=value&param=test",
						},
					},
					{
						Name: "vendor2",
						Request: config.URLPattern{
							URL: "https://api.vendor2.com/v1/recommend",
						},
					},
				},
			},
			wantBody: `[
				{
					"vendor_key": "vendor1",
					"request_host": "api.vendor1.com"
				},
				{
					"vendor_key": "vendor2",
					"request_host": "api.vendor2.com"
				}
			]`,
		},
		{
			name: "GIVEN empty vendor config THEN expect response with empty array",
			vendorConfig: config.VendorConfig{
				Vendors: []config.Vendor{},
			},
			wantBody: `[]`,
		},
		{
			name: "GIVEN vendor with invalid URL THEN expect vendor with empty request host",
			vendorConfig: config.VendorConfig{
				Vendors: []config.Vendor{
					{
						Name: "invalid_url_vendor",
						Request: config.URLPattern{
							URL: "://invalid-url",
						},
					},
				},
			},
			wantBody: `[
				{
					"vendor_key": "invalid_url_vendor",
					"request_host": ""
				}
			]`,
		},
	}

	for _, tc := range tt {
		ts.Run(tc.name, func() {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/vendors", nil)

			vm := NewVendorManager(tc.vendorConfig)
			vm.GetVendors(c)

			require.Equal(ts.T(), http.StatusOK, w.Code)
			require.JSONEq(ts.T(), tc.wantBody, w.Body.String())
		})
	}
}

func TestVendorsTestSuite(t *testing.T) {
	suite.Run(t, new(VendorsTestSuite))
}
