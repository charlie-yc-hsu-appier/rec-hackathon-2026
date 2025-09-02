package utils

import (
	"encoding/base64"
	"strings"
)

func EncodeClickID(clickID string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(clickID))
	return strings.TrimRight(encoded, "=")
}
