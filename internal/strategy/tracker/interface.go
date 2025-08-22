package tracker

type Strategy interface {
	GenerateTrackingURL(params Params) string
}
