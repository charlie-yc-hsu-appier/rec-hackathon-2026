package body

import (
	"fmt"
	"strings"
)

type Adpopcorn struct{}

type adpopcornBody struct {
	App       adpopcornApp       `json:"app"`
	Device    adpopcornDevice    `json:"device"`
	Imp       adpopcornImp       `json:"imp"`
	Affiliate adpopcornAffiliate `json:"affiliate"`
}

type adpopcornApp struct {
	BundleID string `json:"bundleId"`
}

type adpopcornDevice struct {
	ID  string `json:"id"`
	Lmt int    `json:"lmt"`
}

type adpopcornImp struct {
	ImageSize string `json:"imageSize"`
}

type adpopcornAffiliate struct {
	SubID string `json:"subId"`
}

func (s *Adpopcorn) GenerateBody(params Params) any {
	body := adpopcornBody{
		App: adpopcornApp{
			BundleID: "",
		},
		Device: adpopcornDevice{
			ID:  strings.ToLower(params.UserID),
			Lmt: 0,
		},
		Imp: adpopcornImp{
			ImageSize: fmt.Sprintf("%dx%d", params.ImgWidth, params.ImgHeight),
		},
		Affiliate: adpopcornAffiliate{
			SubID: params.SubID,
		},
	}

	return body
}
