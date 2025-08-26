package header

type NilHeader struct{}

func (s *NilHeader) GenerateHeaders(params Params) map[string]string {
	return nil
}
