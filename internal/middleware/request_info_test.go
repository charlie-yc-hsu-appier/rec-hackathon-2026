// TODO: to be removed after gin/nginx retirement
package middleware

import (
	"net/http"
	"net/http/httptest"
	"rec-vendor-api/internal/telemetry"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRequestInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tt := []struct {
		name                string
		headers             map[string]string
		expectedRequestInfo telemetry.RequestInfo
	}{
		{
			name: "GIVEN both headers present THEN expect both values extracted",
			headers: map[string]string{
				"x-rec-siteid": "test-site-123",
				"x-rec-oid":    "test-oid-456",
			},
			expectedRequestInfo: telemetry.RequestInfo{
				SiteID: "test-site-123",
				OID:    "test-oid-456",
			},
		},
		{
			name:                "GIVEN no headers present THEN expect empty values",
			headers:             map[string]string{},
			expectedRequestInfo: telemetry.RequestInfo{},
		},
		{
			name: "GIVEN empty header values THEN expect empty values",
			headers: map[string]string{
				"x-rec-siteid": "",
				"x-rec-oid":    "",
			},
			expectedRequestInfo: telemetry.RequestInfo{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			r.Use(RequestInfo())
			r.GET("/test", func(c *gin.Context) {
				requestInfo := telemetry.RequestInfoFromContext(c.Request.Context())
				require.Equal(t, tc.expectedRequestInfo, requestInfo)
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusOK, w.Code)
		})
	}
}
