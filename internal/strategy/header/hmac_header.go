package header

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
)

type HmacHeader struct {
	AccessKey string
	SecretKey string
	Clock     Clock
}

func (h *HmacHeader) GenerateHeaders(params Params) map[string]string {
	datetimeGMT := h.Clock.getDatetimeGMT()
	parsedURL, _ := url.Parse(params.RequestURL)

	path := parsedURL.Path
	query := parsedURL.Query().Encode()

	message := datetimeGMT + "GET" + path + query
	accessKey := h.AccessKey
	secret := h.SecretKey

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	signature := hex.EncodeToString(mac.Sum(nil))

	return map[string]string{
		"Authorization": fmt.Sprintf("CEA algorithm=HmacSHA256, access-key=%s, signed-date=%s, signature=%s", accessKey, datetimeGMT, signature),
	}
}
