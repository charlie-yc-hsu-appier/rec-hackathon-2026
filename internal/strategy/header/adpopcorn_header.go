package header

type AdpopcornHeader struct {
	UserAgent   string
	ContentType string
}

func (h *AdpopcornHeader) GenerateHeaders(_ Params) map[string]string {
	headers := map[string]string{
		"User-Agent":   h.UserAgent,
		"Content-Type": h.ContentType,
	}
	return headers
}
