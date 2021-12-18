package api

import (
	"context"
	"net/http"
)

type PHPRequest interface {
	GetContext() context.Context
	SetContext(ctx context.Context)
	ServiceName() string
	Close()
	Globals(typ GlobalTrackVars) (m map[string]string)
	GlobalByKey(typ GlobalTrackVars, key string) (val string)
	Headers() http.Header
	Curl(id int32) Curl
	PHPRedis(id uintptr) DB
}
