package body

import (
	"strconv"
	"strings"
)

type PostFormat struct{}

type postFormatRequestBody struct {
	App       app       `json:"app"`
	Device    device    `json:"device"`
	Imp       imp       `json:"imp"`
	Affiliate affiliate `json:"affiliate"`
}

type app struct {
	BundleID string `json:"bundleId"`
}

type device struct {
	ID  string `json:"id"`
	Lmt int    `json:"lmt"`
}

type imp struct {
	ImageSize string `json:"imageSize"`
}

type affiliate struct {
	SubID string `json:"subId"`
}

func (pf *PostFormat) GenerateBody(params Params) interface{} {
	body := postFormatRequestBody{
		App: app{
			BundleID: params.BundleID,
		},
		Device: device{
			ID:  strings.ToLower(params.UserID),
			Lmt: 0,
		},
		Imp: imp{
			ImageSize: strconv.Itoa(params.ImgWidth) + "x" + strconv.Itoa(params.ImgHeight),
		},
		Affiliate: affiliate{
			SubID: params.SubID,
		},
	}

	return body
}
