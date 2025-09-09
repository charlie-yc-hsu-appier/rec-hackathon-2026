package vendor

import "rec-vendor-api/internal/strategy/requester"

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
	ClientIP        string
}

func (r Request) toRequesterParams(url string) requester.Params {
	return requester.Params{
		RequestURL:      url,
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
