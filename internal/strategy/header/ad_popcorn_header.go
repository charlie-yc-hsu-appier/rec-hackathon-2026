package header

type AdPopcornHeader struct {
	UserAgent string
}

func (h *AdPopcornHeader) GenerateHeaders(_ Params) map[string]string {
	headers := map[string]string{
		"User-Agent": h.UserAgent,
	}
	return headers
}
