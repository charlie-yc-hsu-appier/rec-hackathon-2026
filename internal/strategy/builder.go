package strategy

import (
	"rec-vendor-api/internal/config"
	"rec-vendor-api/internal/strategy/header"
	"rec-vendor-api/internal/strategy/requester"
	"rec-vendor-api/internal/strategy/tracker"
	"rec-vendor-api/internal/strategy/unmarshaler"
)

func BuildHeader(v config.Vendor) header.Strategy {
	switch v.Name {
	case "replace":
		return &header.ReplaceHeader{AccessKey: v.AccessKey, SecretKey: v.SecretKey, Clock: &header.ClockImpl{}}
	case "adpopcorn":
		return &header.AdPopcornHeader{UserAgent: v.UserAgent}
	case "keeta":
		return &header.KeetaHeader{SCaApp: v.SCaApp, SCaSecret: v.SCaSecret, Clock: &header.ClockImpl{}}
	default:
		return &header.NoHeader{}
	}
}

func BuildRequester(v config.Vendor) requester.Strategy {
	switch v.Name {
	case "inl_corp_0", "inl_corp_1", "inl_corp_2", "inl_corp_3", "inl_corp_4", "inl_corp_5":
		return &requester.InlCorp{SizeCodeMap: v.SizeCodeMap}
	case "keeta":
		return &requester.Keeta{SceneType: v.SceneType, Ver: v.Ver, ChannelToken: v.ChannelToken}
	default:
		return &requester.Default{}
	}
}

func BuildUnmarshaler(v config.Vendor) unmarshaler.Strategy {
	switch v.Name {
	case "replace", "adpopcorn":
		return &unmarshaler.WrappedCoupangPartner{}
	case "adpacker":
		return &unmarshaler.Adpacker{}
	default:
		return &unmarshaler.CoupangPartner{}
	}
}

func BuildTracker(v config.Vendor) tracker.Strategy {
	switch v.Name {
	default:
		return &tracker.Default{}
	}
}
