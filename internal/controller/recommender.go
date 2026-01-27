package controller

import (
	"fmt"
	"net/http"

	"errors"
	controller_errors "rec-vendor-api/internal/controller/errors"
	"rec-vendor-api/internal/vendor"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Recommender struct {
	vendorRegistry map[string]vendor.Client
}

func NewRecommender(vendorRegistry map[string]vendor.Client) *Recommender {
	return &Recommender{
		vendorRegistry: vendorRegistry,
	}
}

// Recommend godoc
// @Summary      Get vendor recommendations
// @Description  Returns recommended products for a user from a vendor
// @Produce      json
// @Param        vendor_key  path  string true  "Vendor Key"
// @Param        user_id     query string true  "User ID"
// @Param        click_id    query string true  "Click ID"
// @Param        w           query int    true  "Image Width"
// @Param        h           query int    true  "Image Height"
// @Param        web_host    query string false "Web host domain"
// @Param        bundle_id   query string false "App bundle ID"
// @Param        adtype      query int    false "Ad Type (native → 3, else → 2)"
// @Param        partner_id  query string false "Partner ID"
// @Param        os          query string false "Operating System (android, ios)"
// @Success      200 {object} []vendor.ProductInfo
// @Failure      400 {object} map[string]string "Bad Request"
// @Failure      500 {object} map[string]string "Internal Error"
// @Router       /r/{vendor_key} [get]
func (c *Recommender) Recommend(ctx *gin.Context) {
	var req vendor.Request
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("fail to bind query parameter, uri: %s", ctx.Request.RequestURI)
		handleBadRequest(ctx, err)
		return
	}
	req.ClientIP = ctx.ClientIP()

	vendorKey := ctx.Param("vendor_key")
	vendorClient := c.vendorRegistry[vendorKey]
	if vendorClient == nil {
		log.WithContext(ctx).Errorf("Invalid vendor key: %s", vendorKey)
		handleBadRequest(ctx, fmt.Errorf("vendor key '%s' not supported", vendorKey))
		return
	}

	response, err := vendorClient.GetUserRecommendationItems(ctx, req)
	if err != nil {
		var badRequestErr *controller_errors.BadRequestError
		if errors.As(err, &badRequestErr) {
			log.WithContext(ctx).Errorf("VendorClient returned BadRequestError. err: %v", err)
			handleBadRequest(ctx, fmt.Errorf("VendorClient returned BadRequestError. err: %w", err))
			return
		}

		log.WithContext(ctx).Errorf("Fail to recommend any products. err: %v", err)
		handleInternalServerError(ctx, fmt.Errorf("fail to recommend any products for vendor %s. err: %w", vendorKey, err))
		return
	}
	ctx.JSON(http.StatusOK, response)
}
