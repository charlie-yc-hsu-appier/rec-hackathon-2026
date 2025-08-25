package controller

import (
	"fmt"
	"net/http"

	"rec-vendor-api/internal/vendor"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type VendorController struct {
	vendorRegistry map[string]vendor.Client
}

func NewVendorController(vendorRegistry map[string]vendor.Client) *VendorController {
	return &VendorController{
		vendorRegistry: vendorRegistry,
	}
}

type RecommendQuery struct {
	VendorKey string `form:"vendor_key" binding:"required"`
	UserID    string `form:"user_id" binding:"required"`
	ClickID   string `form:"click_id"`
	ImgWidth  int    `form:"w" binding:"required"`
	ImgHeight int    `form:"h" binding:"required"`
}

// Recommend godoc
// @Summary      Get vendor recommendations
// @Description  Returns recommended products for a user from a vendor
// @Produce      json
// @Param        vendor_key  query string true  "Vendor Key"
// @Param        user_id     query string true  "User ID"
// @Param        click_id    query string false "Click ID"
// @Param        w           query int    true  "Image Width"
// @Param        h           query int    true  "Image Height"
// @Success      200 {object} vendor.Response
// @Failure      400 {object} map[string]string "Bad Request"
// @Failure      500 {object} map[string]string "Internal Error"
// @Router       /r [get]
func (c *VendorController) Recommend(ctx *gin.Context) {
	var req RecommendQuery
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("fail to bind query parameter, uri: %s", ctx.Request.RequestURI)
		handleBadRequest(ctx, err)
		return
	}

	serviceReq := vendor.Request{
		UserID:    req.UserID,
		ClickID:   req.ClickID,
		ImgWidth:  req.ImgWidth,
		ImgHeight: req.ImgHeight,
	}

	vendorClient := c.vendorRegistry[req.VendorKey]
	if vendorClient == nil {
		log.WithContext(ctx).Errorf("Invalid vendor key: %s", req.VendorKey)
		handleBadRequest(ctx, fmt.Errorf("Vendor key '%s' not supported", req.VendorKey))
		return
	}

	response, err := vendorClient.GetUserRecommendationItems(ctx.Request.Context(), serviceReq)
	if err != nil {
		log.WithContext(ctx).Errorf("Fail to recommend any products. err: %v", err)
		handleInternalServerError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
