package controller

import (
	"net/http"
	"net/url"
	"rec-vendor-api/internal/config"

	"github.com/gin-gonic/gin"
)

type VendorInfo struct {
	VendorKey   string `json:"vendor_key"`
	RequestHost string `json:"request_host"`
}

type vendorManager struct {
	vendors []VendorInfo
}

func NewVendorManager(cfg config.VendorConfig) *vendorManager {
	vendors := make([]VendorInfo, 0, len(cfg.Vendors))
	for _, vendor := range cfg.Vendors {
		requestHost := ""
		if parsedURL, err := url.Parse(vendor.Request.URL); err == nil {
			requestHost = parsedURL.Host
		}
		vendors = append(vendors, VendorInfo{
			VendorKey:   vendor.Name,
			RequestHost: requestHost,
		})
	}

	return &vendorManager{
		vendors: vendors,
	}
}

// Vendors godoc
// @Summary 	Vendors
// @Description Usage for getting vendor information
// @Produce 	json
// @Success 	200 {array} VendorInfo
// @Router 		/vendors [get]
func (vm *vendorManager) GetVendors(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, vm.vendors)
}
