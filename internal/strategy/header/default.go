package header

type Params struct {
	UserID  string
	ClickID string
}

type Default struct{}

func (s *Default) GenerateHeaders(params Params) map[string]string {
	return nil
}
