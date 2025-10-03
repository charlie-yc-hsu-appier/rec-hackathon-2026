package vendor

import (
	"context"
	"time"

	"rec-vendor-api/internal/config"
	"rec-vendor-api/internal/strategy/header"
	"rec-vendor-api/internal/strategy/unmarshaler"
	"rec-vendor-api/internal/strategy/url"
	"rec-vendor-api/internal/telemetry"

	"github.com/plaxieappier/rec-go-kit/httpkit"
)

type vendorClient struct {
	cfg                   config.Vendor
	client                httpkit.Client
	timeout               time.Duration
	headerStrategy        header.Strategy
	requestURLStrategy    url.Strategy
	respUnmarshalStrategy unmarshaler.Strategy
	trackingURLStrategy   url.Strategy
}

//go:generate mockgen -source=./client.go -destination=./client_mock.go -package=vendor
type Client interface {
	GetUserRecommendationItems(ctx context.Context, req Request) ([]ProductInfo, error)
}

func NewClient(cfg config.Vendor, client httpkit.Client, timeout time.Duration,
	headerStrategy header.Strategy, requestURLStrategy url.Strategy,
	respUnmarshalStrategy unmarshaler.Strategy, trackingURLStrategy url.Strategy) Client {
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
	requestInfo := telemetry.RequestInfoFromContext(ctx)

	requestURL, err := v.requestURLStrategy.GenerateURL(v.cfg.Request, req.toURLParams())
	if err != nil {
		return nil, err
	}
	restReq := httpkit.NewRequest(requestURL)

	headerParams := header.Params{RequestURL: requestURL, UserID: req.UserID}
	headers := v.headerStrategy.GenerateHeaders(headerParams)
	restReq = restReq.PatchHeaders(headers)

	restReq = restReq.SetMetrics(
		telemetry.Metrics.RestApiDurationSeconds.WithLabelValues(v.cfg.Name, requestInfo.SiteID, requestInfo.OID),
		telemetry.Metrics.RestApiErrorTotal.WithLabelValues(v.cfg.Name, requestInfo.SiteID, requestInfo.OID),
	)

	restResp, err := v.client.Get(ctx, restReq, v.timeout, []int{200})
	if err != nil {
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
