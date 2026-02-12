package header

type AdpopcornHeader struct {
	UserAgent string
}

func (h *AdpopcornHeader) GenerateHeaders(_ Params) map[string]string {
	headers := map[string]string{
		"User-Agent": h.UserAgent,
	}
	return headers
}
