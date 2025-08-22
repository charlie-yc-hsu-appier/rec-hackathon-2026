package requester

type Strategy interface {
	GenerateRequestURL(params Params) string
}
