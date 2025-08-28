package header

type NoHeader struct{}

func (s *NoHeader) GenerateHeaders(_ Params) map[string]string {
	return map[string]string{}
}
