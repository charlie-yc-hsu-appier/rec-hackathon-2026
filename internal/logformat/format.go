package logformat

import (
	"rec-vendor-api/internal/telemetry"

	"github.com/plaxieappier/rec-go-kit/logkit"
	log "github.com/sirupsen/logrus"
)

type LogFormat struct {
	logkit.BaseLogFormat

	SiteId string `json:"sid"`
	OID    string `json:"oid"`
	Vendor string `json:"vendor"`
	SubID  string `json:"subid"`
}

func (l *LogFormat) PrepareFormat(entry *log.Entry) any {
	requestInfo := telemetry.RequestInfoFromContext(entry.Context)

	return LogFormat{
		BaseLogFormat: l.BaseLogFormat.PrepareFormat(entry).(logkit.BaseLogFormat),
		SiteId:        requestInfo.SiteID,
		OID:           requestInfo.OID,
		Vendor:        requestInfo.Vendor,
		SubID:         requestInfo.SubID,
	}
}
