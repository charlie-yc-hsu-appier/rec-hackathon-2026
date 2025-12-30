package strategy

import (
	"rec-vendor-api/internal/config"
	"rec-vendor-api/internal/strategy/body"
	"rec-vendor-api/internal/strategy/header"
	"rec-vendor-api/internal/strategy/unmarshaler"
	"rec-vendor-api/internal/strategy/url"
)

func BuildHeader(v config.Vendor) header.Strategy {
	switch v.Name {
	case "replace":
		return &header.ReplaceHeader{AccessKey: v.AccessKey, SecretKey: v.SecretKey, Clock: &header.ClockImpl{}}
	case "adpopcorn":
		return &header.AdpopcornHeader{UserAgent: v.UserAgent, ContentType: v.ContentType}
	case "keeta":
		return &header.KeetaHeader{SCaApp: v.SCaApp, SCaSecret: v.SCaSecret, Clock: &header.ClockImpl{}}
	default:
		return &header.NoHeader{}
	}
}

func BuildRequest(v config.Vendor) url.Strategy {
	switch v.Name {
	default:
		return &url.Default{}
	}
}

func BuildUnmarshaler(v config.Vendor) unmarshaler.Strategy {
	switch v.Name {
	case "adpopcorn":
		return &unmarshaler.Replace{}
	case "adpacker":
		return &unmarshaler.Adpacker{}
	case "keeta":
		return &unmarshaler.Keeta{}
	case "adforus":
		return &unmarshaler.Adforus{}
	case "replace":
		return &unmarshaler.Replace{}
	default:
		return &unmarshaler.CoupangPartner{}
	}
}

func BuildTracking(v config.Vendor) url.Strategy {
	switch v.Name {
	default:
		return &url.Default{}
	}
}

func BuildBody(v config.Vendor) body.Strategy {
	switch v.Name {
	case "replace":
		return &body.Replace{}
	case "adpopcorn":
		return &body.Adpopcorn{}
	default:
		return &body.NoBody{}
	}
}
