package tracker

type Params struct {
	TrackingURL string
	ProductURL  string
	ClickID     string
}

//go:generate mockgen -source=./interface.go -destination=./interface_mock.go -package=tracker

type Strategy interface {
	GenerateTrackingURL(params Params) string
}
