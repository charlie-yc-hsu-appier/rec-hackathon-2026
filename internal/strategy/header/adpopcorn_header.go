package header

type AdpopcornHeader struct {
	UserAgent string
	ContnetType string
}

func (h *AdpopcornHeader) GenerateHeaders(_ Params) map[string]string {
	headers := map[string]string{
		"User-Agent": h.UserAgent,
		"Content-type": h.ContnetType
	}
	return headers
}
