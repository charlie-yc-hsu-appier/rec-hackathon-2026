package header

type Default struct{}

func (s *Default) GenerateHeaders(params Params) map[string]string {
	return nil
}
