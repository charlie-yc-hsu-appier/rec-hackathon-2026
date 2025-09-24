package logformat

import (
	"rec-vendor-api/internal/telemetry"

	"github.com/plaxieappier/rec-go-kit/logkit"
	log "github.com/sirupsen/logrus"
)

type LogFormat struct {
	logkit.BaseLogFormat

	SiteId    string `json:"sid"`
	OID       string `json:"oid"`
	VendorKey string `json:"vendor_key"`
	SubID     string `json:"subid"`
	BidObjID  string `json:"bid_obj_id"`
	ReqID     string `json:"request_id"`
	TraceID   string `json:"trace_id"`
}

func (l *LogFormat) PrepareFormat(entry *log.Entry) any {
	requestInfo := telemetry.RequestInfoFromContext(entry.Context)

	return LogFormat{
		BaseLogFormat: l.BaseLogFormat.PrepareFormat(entry).(logkit.BaseLogFormat),
		SiteId:        requestInfo.SiteID,
		OID:           requestInfo.OID,
		VendorKey:     requestInfo.VendorKey,
		SubID:         requestInfo.SubID,
		BidObjID:      requestInfo.BidObjID,
		ReqID:         requestInfo.ReqID,
		TraceID:       requestInfo.TraceID,
	}
}
