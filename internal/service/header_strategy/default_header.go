package header

type Params struct {
	UserID  string
	ClickID string
}

type NoHeaderStrategy struct{}

func (s *NoHeaderStrategy) GenerateHeaders(params Params) map[string]string {
	return nil
}
