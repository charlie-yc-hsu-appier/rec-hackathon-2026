package rest

import (
	"context"
	"rec-vendor-api/internal/telemetry"
	"time"

	"bitbucket.org/plaxieappier/rec-go-kit/httpkit"
	"github.com/avast/retry-go/v4"
)

var (
	validStatus = []int{200}
)

//go:generate mockgen -source=./client.go -destination=./client_mock.go -package=rest
type Client interface {
	Get(ctx context.Context, restReq *httpkit.Request, timeout time.Duration, vendorKey string, opts ...retry.Option) (*httpkit.Response, error)
}

type client struct {
	httpClient httpkit.Client
}

func NewClient(httpClient httpkit.Client) Client {
	return &client{httpClient: httpClient}
}

func (c *client) Get(ctx context.Context, req *httpkit.Request, timeout time.Duration, vendorKey string, opts ...retry.Option) (*httpkit.Response, error) {
	req = req.SetMetrics(
		telemetry.Metrics.RestApiDurationSeconds.WithLabelValues(vendorKey),
		telemetry.Metrics.RestApiErrorTotal.WithLabelValues(vendorKey),
	)
	resp, err := c.httpClient.Get(ctx, req, timeout, validStatus, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
