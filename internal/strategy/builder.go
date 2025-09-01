package strategy

import (
	"rec-vendor-api/internal/strategy/header"
	"rec-vendor-api/internal/strategy/requester"
	"rec-vendor-api/internal/strategy/tracker"
	"rec-vendor-api/internal/strategy/unmarshaler"
)

func BuildHeader(name, accessKey, secretKey string) header.Strategy {
	switch name {
	case "replace":
		return &header.HmacHeader{AccessKey: accessKey, SecretKey: secretKey, Clock: &header.ClockImpl{}}
	default:
		return &header.NoHeader{}
	}
}

func BuildRequester(name string) requester.Strategy {
	switch name {
	default:
		return &requester.Default{}
	}
}

func BuildUnmarshaler(name string) unmarshaler.Strategy {
	switch name {
	case "replace":
		return &unmarshaler.WrappedCoupangPartner{}
	default:
		return &unmarshaler.CoupangPartner{}
	}
}

func BuildTracker(name string) tracker.Strategy {
	switch name {
	default:
		return &tracker.Default{}
	}
}
