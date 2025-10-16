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

	// decide to use upper case or lower case user id
	OS string
}

//go:generate mockgen -source=./interface.go -destination=./interface_mock.go -package=url

// Strategy defines the contract for URL generation.
// Implementations of GenerateURL should return the URL (with query parameters) as a string,
type Strategy interface {
	GenerateURL(urlPattern config.URLPattern, params Params) (string, error)
}
