package requester

type Params struct {
	RequestURL      string
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
}

//go:generate mockgen -source=./interface.go -destination=./interface_mock.go -package=requester
type Strategy interface {
	GenerateRequestURL(params Params) (string, error)
}
