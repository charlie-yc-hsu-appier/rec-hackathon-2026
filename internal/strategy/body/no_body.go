package body

type NoBody struct{}

func (s *NoBody) GenerateBody(_ Params) any {
	return nil
}
