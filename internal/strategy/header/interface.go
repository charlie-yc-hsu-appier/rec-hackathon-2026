package header

type Params struct {
	// Add any necessary fields for header generation
}

//go:generate mockgen -source=./interface.go -destination=./interface_mock.go -package=header

type Strategy interface {
	GenerateHeaders(params Params) map[string]string
}
