package utils

import (
	"encoding/base64"
	"sort"
	"strings"
)

func EncodeClickID(clickID string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(clickID))
	return strings.TrimRight(encoded, "=")
}

// GetSortedStringKeys returns the sorted keys of a string map in ascending order.
func GetSortedStringKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
