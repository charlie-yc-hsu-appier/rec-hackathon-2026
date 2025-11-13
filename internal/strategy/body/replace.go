package body

import (
	"fmt"
	"rec-vendor-api/internal/strategy/utils"
	"strings"
)

type Replace struct{}

type replaceBody struct {
	App       replaceApp       `json:"app"`
	Device    replaceDevice    `json:"device"`
	Imp       replaceImp       `json:"imp"`
	Affiliate replaceAffiliate `json:"affiliate"`
	User      replaceUser      `json:"user"`
}

type replaceApp struct {
	BundleID string `json:"bundleId"`
}

type replaceDevice struct {
	ID  string `json:"id"`
	Lmt int    `json:"lmt"`
}

type replaceImp struct {
	ImageSize string `json:"imageSize"`
}

type replaceAffiliate struct {
	SubID    string `json:"subId"`
	SubParam string `json:"subParam"`
}

type replaceUser struct {
	Puid string `json:"puid"`
}

func (s *Replace) GenerateBody(params Params) any {
	clickIDBase64 := utils.EncodeClickID(params.ClickID)
	body := replaceBody{
		App: replaceApp{
			BundleID: params.BundleID,
		},
		Device: replaceDevice{
			ID:  strings.ToLower(params.UserID),
			Lmt: 0,
		},
		Imp: replaceImp{
			ImageSize: fmt.Sprintf("%dx%d", params.ImgWidth, params.ImgHeight),
		},
		Affiliate: replaceAffiliate{
			SubID:    params.SubID,
			SubParam: clickIDBase64,
		},
		User: replaceUser{
			Puid: clickIDBase64,
		},
	}

	return body
}
