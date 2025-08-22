package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type VendorController struct {
	// TODO
}

// TODO: add swagger documentation
func NewVendorController() *VendorController {
	return &VendorController{}
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

	// TODO:
	// response, err := vendorClient.GetRecommendationItems(ctx.Request.Context(), serviceReq)

	ctx.JSON(http.StatusOK, "message: success")
}
