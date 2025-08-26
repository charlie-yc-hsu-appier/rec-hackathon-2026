package vendor

import (
	"context"
	"strconv"
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
		ImgWidth:   req.ImgWidth,
		ImgHeight:  req.ImgHeight,
	}
	url := v.requestURLStrategy.GenerateRequestURL(reqParams)
	restReq := httpkit.NewRequest(url)

	headerParams := header.Params{}
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

	cpResp, err := v.respUnmarshalStrategy.UnmarshalResponse(restResp.Body)
	if err != nil {
		return nil, err
	}

	products := make([]ProductInfo, 0, len(*cpResp))

	for _, ele := range *cpResp {
		trackParams := tracker.Params{
			TrackingURL: v.cfg.TrackingURL,
			ProductURL:  ele.ProductURL,
			ClickID:     req.ClickID,
		}
		productUrl := v.trackingURLStrategy.GenerateTrackingURL(trackParams)

		products = append(products, ProductInfo{
			ProductID: strconv.Itoa(ele.ProductID),
			Url:       productUrl,
			Image:     ele.ProductImage,
		})
	}

	return products, nil
}
