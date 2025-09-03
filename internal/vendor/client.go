package vendor

import (
	"context"
	"time"

	"rec-vendor-api/internal/config"
	"rec-vendor-api/internal/strategy/header"
	"rec-vendor-api/internal/strategy/requester"
	"rec-vendor-api/internal/strategy/tracker"
	"rec-vendor-api/internal/strategy/unmarshaler"
	"rec-vendor-api/internal/telemetry"

	"github.com/plaxieappier/rec-go-kit/httpkit"
)

type vendorClient struct {
	cfg                   config.Vendor
	client                httpkit.Client
	timeout               time.Duration
	headerStrategy        header.Strategy
	requestURLStrategy    requester.Strategy
	respUnmarshalStrategy unmarshaler.Strategy
	trackingURLStrategy   tracker.Strategy
}

//go:generate mockgen -source=./client.go -destination=./client_mock.go -package=vendor
type Client interface {
	GetUserRecommendationItems(ctx context.Context, req Request) ([]ProductInfo, error)
}

func NewClient(cfg config.Vendor, client httpkit.Client, timeout time.Duration,
	headerStrategy header.Strategy, requestURLStrategy requester.Strategy,
	respUnmarshalStrategy unmarshaler.Strategy, trackingURLStrategy tracker.Strategy) Client {
	return &vendorClient{
		cfg:                   cfg,
		client:                client,
		timeout:               timeout,
		headerStrategy:        headerStrategy,
		requestURLStrategy:    requestURLStrategy,
		respUnmarshalStrategy: respUnmarshalStrategy,
		trackingURLStrategy:   trackingURLStrategy,
	}
}

func (v *vendorClient) GetUserRecommendationItems(ctx context.Context, req Request) ([]ProductInfo, error) {
	reqParams := requester.Params{
		RequestURL: v.cfg.RequestURL,
		UserID:     req.UserID,
		ClickID:    req.ClickID,
		ImgWidth:   req.ImgWidth,
		ImgHeight:  req.ImgHeight,
		WebHost:    req.WebHost,
		BundleID:   req.BundleID,
		AdType:     req.AdType,
	}
	url, err := v.requestURLStrategy.GenerateRequestURL(reqParams)
	if err != nil {
		return nil, err
	}
	restReq := httpkit.NewRequest(url)

	headerParams := header.Params{RequestURL: url}
	headers := v.headerStrategy.GenerateHeaders(headerParams)
	restReq = restReq.PatchHeaders(headers)

	restReq = restReq.SetMetrics(
		telemetry.Metrics.RestApiDurationSeconds.WithLabelValues(v.cfg.Name),
		telemetry.Metrics.RestApiErrorTotal.WithLabelValues(v.cfg.Name),
	)

	restResp, err := v.client.Get(ctx, restReq, v.timeout, []int{200})
	if err != nil {
		return nil, err
	}

	res, err := v.respUnmarshalStrategy.UnmarshalResponse(restResp.Body)
	if err != nil {
		telemetry.Metrics.RestApiAnomalyTotal.WithLabelValues(v.cfg.Name, err.Error()).Inc()
		return nil, err
	}

	products := make([]ProductInfo, 0, len(res))

	for _, ele := range res {
		trackParams := tracker.Params{
			TrackingURL: v.cfg.TrackingURL,
			ProductURL:  ele.ProductURL,
			ClickID:     req.ClickID,
			UserID:      req.UserID,
		}
		productUrl := v.trackingURLStrategy.GenerateTrackingURL(trackParams)

		products = append(products, ProductInfo{
			ProductID: ele.ProductID,
			Url:       productUrl,
			Image:     ele.ProductImage,
		})
	}

	return products, nil
}
