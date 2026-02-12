package body

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVendorPostBody_GenerateBody(t *testing.T) {
	params := Params{
		UserID:    "testUser",
		SubID:     "testSubID",
		ImgWidth:  320,
		ImgHeight: 480,
		BundleID:  "testBundle",
	}

	bodyBuilder := &VendorPostBody{}
	body := bodyBuilder.GenerateBody(params)

	bodyBytes, err := json.Marshal(body)
	assert.NoError(t, err)

	var bodyData vendorPostBody
	err = json.Unmarshal(bodyBytes, &bodyData)
	assert.NoError(t, err)

	assert.Equal(t, params.BundleID, bodyData.App.BundleID)
	assert.Equal(t, strings.ToLower(params.UserID), bodyData.Device.ID)
	assert.Equal(t, 0, bodyData.Device.Lmt)
	assert.Equal(t, strconv.Itoa(params.ImgWidth)+"x"+strconv.Itoa(params.ImgHeight), bodyData.Imp.ImageSize)
	assert.Equal(t, params.SubID, bodyData.Affiliate.SubID)
}
