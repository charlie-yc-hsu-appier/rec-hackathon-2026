package header

//go:generate mockgen -source=./interface.go -destination=./interface_mock.go -package=header

type Strategy interface {
	GenerateHeaders(params Params) map[string]string
}
