package url

import (
	"rec-vendor-api/internal/config"
)

type Params struct {
	UserID          string
	ClickID         string
	ImgWidth        int
	ImgHeight       int
	WebHost         string
	BundleID        string
	AdType          int
	PartnerID       string
	ClientIP        string
	KeetaCampaignID string
	Latitude        string
	Longitude       string
	SubID           string

	// tracking
	ProductURL string
}

//go:generate mockgen -source=./interface.go -destination=./interface_mock.go -package=url
type Strategy interface {
	GenerateURL(urlPattern config.URLPattern, params Params) (string, map[string]string, error)
}
