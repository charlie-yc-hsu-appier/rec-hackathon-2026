package tracker

//go:generate mockgen -source=./interface.go -destination=./interface_mock.go -package=tracker

type Strategy interface {
	GenerateTrackingURL(params Params) string
}
