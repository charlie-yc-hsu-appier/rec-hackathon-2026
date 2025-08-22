package header

type Strategy interface {
	GenerateHeaders(params Params) map[string]string
}
