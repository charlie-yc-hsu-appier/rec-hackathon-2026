package service

import (
	"context"
	"strconv"
	"time"

	"rec-vendor-api/internal/config"
	header "rec-vendor-api/internal/service/header_strategy"
	requester "rec-vendor-api/internal/service/request_strategy"
	trackurl "rec-vendor-api/internal/service/trackurl_strategy"
	unmarshaler "rec-vendor-api/internal/service/unmarshal_strategy"
	"rec-vendor-api/internal/telemetry"

	"bitbucket.org/plaxieappier/rec-go-kit/httpkit"
)

type vendorClient struct {
	cfg                   config.Vendor
	client                httpkit.Client
	timeout               time.Duration
	headerStrategy        HeaderStrategy
	requestURLStrategy    RequestURLStrategy
	respUnmarshalStrategy RespUnmarshalStrategy
	trackingURLStrategy   TrackingURLStrategy
}

type HeaderStrategy interface {
	GenerateHeaders(params header.Params) map[string]string
}

type RequestURLStrategy interface {
	GenerateRequestURL(params requester.Params) string
}

type RespUnmarshalStrategy interface {
	UnmarshalResponse(body []byte) (*[]unmarshaler.CoupangPartnerResp, error)
}

type TrackingURLStrategy interface {
	GenerateTrackingURL(params trackurl.Params) string
}

type Request struct {
	UserID    string
	ClickID   string
	ImgWidth  int
	ImgHeight int
}

type Client interface {
	GetUserRecommendationItems(ctx context.Context, req Request) (Response, error)
}

func NewVendorClient(cfg config.Vendor, client httpkit.Client, timeout time.Duration,
	headerStrategy HeaderStrategy, requestURLStrategy RequestURLStrategy,
	respUnmarshalStrategy RespUnmarshalStrategy, trackingURLStrategy TrackingURLStrategy) Client {
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

func (v *vendorClient) GetUserRecommendationItems(ctx context.Context, req Request) (Response, error) {
	reqParams := requester.Params{
		RequestURL: v.cfg.RequestURL,
		UserID:     req.UserID,
		ImgWidth:   req.ImgWidth,
		ImgHeight:  req.ImgHeight,
	}
	url := v.requestURLStrategy.GenerateRequestURL(reqParams)
	restReq := httpkit.NewRequest(url)

	headerParams := header.Params{
		UserID:  req.UserID,
		ClickID: req.ClickID,
	}
	headers := v.headerStrategy.GenerateHeaders(headerParams)
	restReq = restReq.PatchHeaders(headers)

	restReq = restReq.SetMetrics(
		telemetry.Metrics.RestApiDurationSeconds.WithLabelValues(v.cfg.Name),
		telemetry.Metrics.RestApiErrorTotal.WithLabelValues(v.cfg.Name),
	)

	restResp, err := v.client.Get(ctx, restReq, v.timeout, []int{200})
	if err != nil {
		return Response{}, err
	}

	cpResp, err := v.respUnmarshalStrategy.UnmarshalResponse(restResp.Body)
	if err != nil {
		return Response{}, err
	}

	productIDs := make([]string, 0, len(*cpResp))
	productPatch := make(map[string]ProductPatch, len(*cpResp))

	for _, ele := range *cpResp {
		productIDStr := strconv.Itoa(ele.ProductID)
		productIDs = append(productIDs, productIDStr)

		trackParams := trackurl.Params{
			TrackingURL: v.cfg.TrackingURL,
			ProductURL:  ele.ProductUrl,
			ClickID:     req.ClickID,
		}
		productUrl := v.trackingURLStrategy.GenerateTrackingURL(trackParams)

		productPatch[productIDStr] = ProductPatch{
			Url:   productUrl,
			Image: ele.ProductImage,
		}
	}

	return Response{
		ProductIDs:   productIDs,
		ProductPatch: productPatch,
	}, nil
}
