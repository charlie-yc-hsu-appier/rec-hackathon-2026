package strategy

import (
	"rec-vendor-api/internal/strategy/header"
	"rec-vendor-api/internal/strategy/requester"
	"rec-vendor-api/internal/strategy/tracker"
	"rec-vendor-api/internal/strategy/unmarshaler"
)

func BuildHeader(name string) header.Strategy {
	switch name {
	default:
		return &header.NilHeader{}
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
