package header

import (
	"strconv"
	"time"
)

type Params struct {
	RequestURL string
	UserID     string
	HTTPMethod string
}

//go:generate mockgen -source=./interface.go -destination=./interface_mock.go -package=header

type Strategy interface {
	GenerateHeaders(params Params) map[string]string
}

type Clock interface {
	getDatetimeGMT() string
	getCurrentMilliTimestamp() string
}

type ClockImpl struct{}

// Response format: yymmdd'T'HHMMSS'Z'
func (ClockImpl) getDatetimeGMT() string {
	return time.Now().UTC().Format("060102T150405Z")
}

func (ClockImpl) getCurrentMilliTimestamp() string {
	return strconv.FormatInt(time.Now().UnixMilli(), 10)
}
