package controller

import (
	"errors"
	"fmt"
	"net/http"

	"rec-vendor-api/internal/service"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type VendorController struct {
	vendorRegistry map[string]service.Client
}

func NewVendorController(vendorRegistry map[string]service.Client) *VendorController {
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

func (c *VendorController) Recommend(ctx *gin.Context) {
	var req RecommendQuery
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("fail to bind query parameter, uri: %s", ctx.Request.RequestURI)
		handleBadRequest(ctx, err)
		return
	}

	serviceReq := service.Request{
		UserID:    req.UserID,
		ClickID:   req.ClickID,
		ImgWidth:  req.ImgWidth,
		ImgHeight: req.ImgHeight,
	}

	vendorClient := c.vendorRegistry[req.VendorKey]
	if vendorClient == nil {
		log.WithContext(ctx).Errorf("Invalid vendor key: %s", req.VendorKey)
		handleBadRequest(ctx, errors.New(fmt.Sprintf("Vendor key '%s' not found", req.VendorKey)))
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
