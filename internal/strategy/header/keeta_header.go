package header

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"

	"rec-vendor-api/internal/strategy/utils"
)

type KeetaHeader struct {
	SCaApp    string
	SCaSecret string
	Clock     Clock
}

func (s *KeetaHeader) GenerateHeaders(params Params) map[string]string {
	headers := map[string]string{
		// Note that in the example of vendor API document, the header key "IDFA" is all capitalized
		//
		// However since the rest client of Go automatically transformed the header key into camel-casing "Idfa"
		// when sending the request, if we use "IDFA" to generate the signature, the header key sent with the request would be
		// "Idfa", which violated the signature protocol.
		//
		// Due to the convention that HTTP headers is case-insensitive, the key is changed to "Idfa" so that the
		// key sent in header and used in signature is consistent
		"Idfa":           params.UserID,
		"S-Ca-App":       s.SCaApp,
		"S-Ca-Timestamp": s.Clock.getCurrentMilliTimestamp(),
	}

	parsedURL, _ := url.Parse(params.RequestURL)
	pathWithQuery := parsedURL.RequestURI()

	sCaSignature := genKeetaSignature(s.SCaSecret, http.MethodGet, headers, pathWithQuery)

	signedHeaders := utils.GetSortedStringKeys(headers)
	sCaSignatureHeaders := strings.Join(signedHeaders, ",")

	headers["S-Ca-Signature"] = sCaSignature
	headers["S-Ca-Signature-Headers"] = sCaSignatureHeaders

	return headers
}

func genKeetaSignature(secret string, method string, headers map[string]string, pathWithQuery string) string {
	var data strings.Builder

	// Method
	data.WriteString(method)
	data.WriteString("\n")

	// MD5
	// No additional handling needed besides writing a single "\n", as per keeta API docs
	data.WriteString("\n")

	// Headers
	for _, k := range utils.GetSortedStringKeys(headers) {
		data.WriteString(k)
		data.WriteString(":")
		data.WriteString(headers[k])
		data.WriteString("\n")
	}
	// Note that according to the "description" of vendor API docs,
	// we should write an additional "\n" here.
	// However the "signature example" provided in the vendor docs indicates that the "\n"
	// is not needed here.
	// data.WriteString("\n")

	// URL
	data.WriteString(pathWithQuery)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data.String()))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
