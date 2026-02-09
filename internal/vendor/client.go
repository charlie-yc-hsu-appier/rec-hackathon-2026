package vendor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"rec-vendor-api/internal/config"
	"rec-vendor-api/internal/strategy/body"
	"rec-vendor-api/internal/strategy/header"
	"rec-vendor-api/internal/strategy/unmarshaler"
	"rec-vendor-api/internal/strategy/url"
	"rec-vendor-api/internal/telemetry"

	"github.com/plaxieappier/rec-go-kit/httpkit"
)

const (
	errNetworkTimeout        = "network timeout"
	errRemoteConnectionReset = "remote connection reset"
	errInvalidHTTPStatus     = "invalid http status: "
	errUnknownNetworkError   = "unknown network error"
)

type vendorClient struct {
	cfg                   config.Vendor
	client                httpkit.Client
	timeout               time.Duration
	headerStrategy        header.Strategy
	requestURLStrategy    url.Strategy
	bodyStrategy          body.Strategy
	respUnmarshalStrategy unmarshaler.Strategy
	trackingURLStrategy   url.Strategy
}

//go:generate mockgen -source=./client.go -destination=./client_mock.go -package=vendor
type Client interface {
	GetUserRecommendationItems(ctx context.Context, req Request) ([]ProductInfo, error)
}

func NewClient(cfg config.Vendor, client httpkit.Client, timeout time.Duration,
	headerStrategy header.Strategy, requestURLStrategy url.Strategy,
	bodyStrategy body.Strategy, respUnmarshalStrategy unmarshaler.Strategy,
	trackingURLStrategy url.Strategy) Client {
	return &vendorClient{
		cfg:                   cfg,
		client:                client,
		timeout:               timeout,
		headerStrategy:        headerStrategy,
		requestURLStrategy:    requestURLStrategy,
		bodyStrategy:          bodyStrategy,
		respUnmarshalStrategy: respUnmarshalStrategy,
		trackingURLStrategy:   trackingURLStrategy,
	}
}

func (v *vendorClient) GetUserRecommendationItems(ctx context.Context, req Request) ([]ProductInfo, error) {
	requestInfo := telemetry.RequestInfoFromContext(ctx)

	requestURL, err := v.requestURLStrategy.GenerateURL(v.cfg.Request, req.toURLParams())
	if err != nil {
		return nil, err
	}
	restReq := httpkit.NewRequest(requestURL)

	headerParams := header.Params{RequestURL: requestURL, UserID: req.UserID, HTTPMethod: v.cfg.HTTPMethod}
	headers := v.headerStrategy.GenerateHeaders(headerParams)
	restReq = restReq.PatchHeaders(headers)

	restReq = restReq.SetMetrics(
		telemetry.Metrics.RestApiDurationSeconds.WithLabelValues(v.cfg.Name, requestInfo.SiteID, requestInfo.OID),
		telemetry.Metrics.RestApiErrorTotal.WithLabelValues(v.cfg.Name, requestInfo.SiteID, requestInfo.OID),
	)

	var restResp *httpkit.Response
	switch v.cfg.HTTPMethod {
	case http.MethodGet:
		restResp, err = v.client.Get(ctx, restReq, v.timeout, []int{200})
	case http.MethodPost:
		bodyObj := v.bodyStrategy.GenerateBody(req.toBodyParams())
		restReq = restReq.SetBody(bodyObj)
		restResp, err = v.client.Post(ctx, restReq, v.timeout, []int{200})
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s (supported: GET, POST)", v.cfg.HTTPMethod)
	}
	if err != nil {
		categorized := categorizeError(restResp, err)
		telemetry.Metrics.RestApiAnomalyTotal.WithLabelValues(v.cfg.Name, requestInfo.SiteID, requestInfo.OID, categorized).Inc()
		return nil, err
	}

	res, err := v.respUnmarshalStrategy.UnmarshalResponse(ctx, restResp.Body)
	if err != nil {
		telemetry.Metrics.RestApiAnomalyTotal.WithLabelValues(v.cfg.Name, requestInfo.SiteID, requestInfo.OID, err.Error()).Inc()
		return nil, err
	}

	products := make([]ProductInfo, 0, len(res))

	for _, ele := range res {
		trackParams := url.Params{
			ProductURL: ele.ProductURL,
			ClickID:    req.ClickID,
			UserID:     req.UserID,
			OS:         req.OS,
		}
		productURL, err := v.trackingURLStrategy.GenerateURL(v.cfg.Tracking, trackParams)
		if err != nil {
			return nil, err
		}

		products = append(products, ProductInfo{
			ProductID: ele.ProductID,
			Url:       productURL,
			Image:     ele.ProductImage,
			Price:     ele.ProductPrice,
			SalePrice: ele.ProductSalePrice,
			Currency:  ele.ProductCurrency,
		})
	}

	return products, nil
}

func categorizeError(restResp *httpkit.Response, err error) string {
	if err == nil {
		return ""
	}
	if isTimeoutError(err) {
		return errNetworkTimeout
	}
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) ||
		strings.Contains(err.Error(), "EOF") || strings.Contains(err.Error(), "connection reset") {
		return errRemoteConnectionReset
	}
	if restResp != nil {
		return errInvalidHTTPStatus + strconv.Itoa(restResp.StatusCode)
	}
	return errUnknownNetworkError
}

func isTimeoutError(err error) bool {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}
