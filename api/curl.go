package api

import (
	"time"
	"unsafe"
)

var (
	CurlGetInfo           func(ex unsafe.Pointer, callback func(ret unsafe.Pointer))
	Globals               func(typ GlobalTrackVars, callback func(ret unsafe.Pointer))
	GlobalByKey           func(typ GlobalTrackVars, key string, callback func(ret unsafe.Pointer))
	CurlInjectTraceHeader func(id int32, headers []string)
	RedisAddr             func(obj unsafe.Pointer) (host string, port int, ok bool)
)

type GlobalTrackVars int

const (
	GlobalTrackVarsPost    GlobalTrackVars = 0 // TRACK_VARS_POST
	GlobalTrackVarsGet     GlobalTrackVars = 1 // TRACK_VARS_GET
	GlobalTrackVarsCookie  GlobalTrackVars = 2 // TRACK_VARS_COOKIE
	GlobalTrackVarsServer  GlobalTrackVars = 3 // TRACK_VARS_SERVER
	GlobalTrackVarsEnv     GlobalTrackVars = 4 // TRACK_VARS_ENV
	GlobalTrackVarsFiles   GlobalTrackVars = 5 // TRACK_VARS_FILES
	GlobalTrackVarsRequest GlobalTrackVars = 6 // TRACK_VARS_REQUEST
)

type Curl interface {
	Begin() (t time.Time)
	End() (t time.Time)

	Url() string
	Method() string
	Scheme() string
	Host() string
	Port() string
	Path() string
	Query() string

	Body() string
	BodySize() int
	ContentType() string

	HttpCode() int
	HttpVersion() string

	TotalTime() float64
	NameLookupTime() float64
	ConnectTime() float64
	PreTransferTime() float64
	StartTransferTime() float64

	RemoteIP() string
	RemotePort() int

	Headers() []string

	SetBegin(v time.Time)
	SetEnd(v time.Time)

	SetUrl(v string)
	SetMethod(v string)
	SetScheme(v string)
	SetHost(v string)
	SetPort(v string)
	SetPath(v string)
	SetQuery(v string)

	SetBody(v string)
	SetBodySize(v int)
	SetContentType(v string)

	SetHttpCode(v int)
	SetHttpVersion(v string)

	SetTotalTime(v float64)
	SetNameLookupTime(v float64)
	SetConnectTime(v float64)
	SetPreTransferTime(v float64)
	SetStartTransferTime(v float64)

	SetRemoteIP(v string)
	SetRemotePort(v int)

	AppendHeader(header string)

	BeginTrace(request PHPRequest, id int32)
	EndTrace(request PHPRequest, id int32)
}
