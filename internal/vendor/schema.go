package vendor

import (
	"rec-vendor-api/internal/strategy/body"
	"rec-vendor-api/internal/strategy/url"
)

type Request struct {
	UserID          string `form:"user_id" binding:"required"`
	ClickID         string `form:"click_id" binding:"required"`
	ImgWidth        int    `form:"w" binding:"required"`
	ImgHeight       int    `form:"h" binding:"required"`
	WebHost         string `form:"web_host"`
	BundleID        string `form:"bundle_id"`
	AdType          int    `form:"adtype"`
	PartnerID       string `form:"partner_id"`
	KeetaCampaignID string `form:"k_campaign_id"`
	Latitude        string `form:"lat"`
	Longitude       string `form:"lon"`
	SubID           string `form:"subid"`
	OS              string `form:"os"`
	ClientIP        string
}

func (r Request) toURLParams() url.Params {
	return url.Params{
		UserID:          r.UserID,
		ClickID:         r.ClickID,
		ImgWidth:        r.ImgWidth,
		ImgHeight:       r.ImgHeight,
		WebHost:         r.WebHost,
		BundleID:        r.BundleID,
		AdType:          r.AdType,
		PartnerID:       r.PartnerID,
		ClientIP:        r.ClientIP,
		KeetaCampaignID: r.KeetaCampaignID,
		Latitude:        r.Latitude,
		Longitude:       r.Longitude,
		SubID:           r.SubID,
		OS:              r.OS,
	}
}

func (r Request) toBodyParams() body.Params {
	return body.Params{
		UserID:    r.UserID,
		ClickID:   r.ClickID,
		ImgWidth:  r.ImgWidth,
		ImgHeight: r.ImgHeight,
		BundleID:  r.BundleID,
		SubID:     r.SubID,
	}
}

type ProductInfo struct {
	ProductID string `json:"product_id"`
	Url       string `json:"url"`
	Image     string `json:"image"`
	Price     string `json:"price"`
	SalePrice string `json:"sale_price"`
	Currency  string `json:"currency"`
}
