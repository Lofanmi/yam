package php5

/*
#include "extension.h"
#include <stdlib.h>
#include <curl/curl.h>
#include <curl/easy.h>
*/
import "C"
import (
	"context"
	"log"
	"unsafe"

	"github.com/Lofanmi/yam/api"
	"github.com/Lofanmi/yam/internal/sdks"
	"github.com/Lofanmi/yam/internal/zend"
	"github.com/hashicorp/go-version"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

var (
	phpRequest api.PHPRequest
	hook       api.Hook
	span       opentracing.Span

	v200, v228, v400, v430 *version.Version
)

func initPHPRedisVersion() {
	v200, _ = version.NewVersion("2.0.0")
	v228, _ = version.NewVersion("2.2.8")
	v400, _ = version.NewVersion("4.0.0")
	v430, _ = version.NewVersion("4.3.0")
}

//export ModuleStartup
func ModuleStartup() int {

	initPHPRedisVersion()

	api.Globals = globals
	api.GlobalByKey = globalByKey
	api.CurlGetInfo = curlGetInfo
	api.CurlInjectTraceHeader = curlInjectTraceHeader
	api.RedisAddr = redisAddr
	return 0
}

//export ModuleShutdown
func ModuleShutdown() int {
	return 0
}

//export RequestStartup
func RequestStartup() int {
	phpRequest = sdks.NewPHPRequest()
	hook = sdks.NewHook(phpRequest)
	var (
		ctx        = context.Background()
		options    = []opentracing.StartSpanOption{opentracing.Tag{Key: string(ext.Component), Value: "HTTP"}}
		tracer     = opentracing.GlobalTracer()
		headers    = phpRequest.Headers()
		xRequestID = headers.Get(sdks.HeaderXRequestID)
	)
	// 检查有无上游调用
	if upstreamCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(headers)); err == nil {
		options = append(options, opentracing.ChildOf(upstreamCtx))
	}
	// 创建span
	span, ctx = opentracing.StartSpanFromContext(ctx, phpRequest.ServiceName(), options...)
	if traceID := sdks.GetTraceIDFromContext(ctx); len(traceID) > 0 && sdks.GetTraceIDFromSpanContext(span.Context()) != traceID {
		sdks.UpdateSpanTraceID(span, traceID)
	}
	if len(xRequestID) > 0 {
		span.SetTag(sdks.HeaderXRequestID, xRequestID)
	}
	// 如果 xRequestID 为空 或 无法设置span的traceID时 使用新的span自动生成的traceID
	xRequestID = sdks.UpdateSpanTraceID(span, xRequestID)
	// context 记录信息
	ctx = sdks.WithContext(ctx, xRequestID)
	ctx = sdks.SpanWithContext(ctx, span)

	phpRequest.SetContext(ctx)

	return 0
}

//export RequestShutdown
func RequestShutdown() int {
	span.Finish()
	phpRequest.Close()
	return 0
}

//export BeforeExecuteEx
func BeforeExecuteEx(opArray unsafe.Pointer, executeData unsafe.Pointer) {
	safeRun(func() { hook.BeforeExecuteEx(executeData) })
}

//export AfterExecuteEx
func AfterExecuteEx(opArray unsafe.Pointer, executeData unsafe.Pointer) {
	safeRun(func() { hook.AfterExecuteEx(executeData) })
}

//export BeforeExecuteInternal
func BeforeExecuteInternal(executeData unsafe.Pointer, returnValue unsafe.Pointer) {
	safeRun(func() { hook.BeforeExecuteInternal(executeData, returnValue) })
}

//export AfterExecuteInternal
func AfterExecuteInternal(executeData unsafe.Pointer, returnValue unsafe.Pointer) {
	safeRun(func() { hook.AfterExecuteInternal(executeData, returnValue) })
}

func safeRun(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	fn()
}

func globalByKey(typ api.GlobalTrackVars, key string, callback func(ret unsafe.Pointer)) {
	var s string
	switch typ {
	case api.GlobalTrackVarsPost:
		s = "_POST"
	case api.GlobalTrackVarsGet:
		s = "_GET"
	case api.GlobalTrackVarsCookie:
		s = "_COOKIE"
	case api.GlobalTrackVarsServer:
		s = "_SERVER"
	case api.GlobalTrackVarsEnv:
		s = "_ENV"
	case api.GlobalTrackVarsFiles:
		s = "_FILES"
	case api.GlobalTrackVarsRequest:
		s = "_REQUEST"
	default:
		return
	}
	sp := C.CString(s)
	sl := len(s)
	defer func() {
		C.free(unsafe.Pointer(sp))
	}()
	kp := C.CString(key)
	kl := len(key)
	defer func() {
		C.free(unsafe.Pointer(kp))
	}()
	C.yam_global_by_key(sp, C.int(sl), kp, C.int(kl), C.ulonglong(uint64(uintptr(unsafe.Pointer(&callback)))))
}

//export globalByKeyCallback
func globalByKeyCallback(ret unsafe.Pointer, fn unsafe.Pointer) {
	callback := *((*func(pointer unsafe.Pointer))(fn))
	callback(ret)
}

func globals(typ api.GlobalTrackVars, callback func(ret unsafe.Pointer)) {
	var s string
	switch typ {
	case api.GlobalTrackVarsPost:
		s = "_POST"
	case api.GlobalTrackVarsGet:
		s = "_GET"
	case api.GlobalTrackVarsCookie:
		s = "_COOKIE"
	case api.GlobalTrackVarsServer:
		s = "_SERVER"
	case api.GlobalTrackVarsEnv:
		s = "_ENV"
	case api.GlobalTrackVarsFiles:
		s = "_FILES"
	case api.GlobalTrackVarsRequest:
		s = "_REQUEST"
	default:
		return
	}
	sp := C.CString(s)
	sl := len(s)
	defer func() {
		C.free(unsafe.Pointer(sp))
	}()
	C.yam_globals(sp, C.int(sl), C.ulonglong(uint64(uintptr(unsafe.Pointer(&callback)))))
}

//export globalsCallback
func globalsCallback(ret unsafe.Pointer, fn unsafe.Pointer) {
	callback := *((*func(pointer unsafe.Pointer))(fn))
	callback(ret)
}

func curlGetInfo(ex unsafe.Pointer, callback func(ret unsafe.Pointer)) {
	C.yam_curl_getinfo(ex, C.ulonglong(uint64(uintptr(unsafe.Pointer(&callback)))))
}

//export curlGetInfoCallback
func curlGetInfoCallback(ret unsafe.Pointer, fn unsafe.Pointer) {
	callback := *((*func(pointer unsafe.Pointer))(fn))
	callback(ret)
}

func curlInjectTraceHeader(id int32, headers []string) {
	var sl *C.struct_curl_slist
	for _, _header := range headers {
		header := C.CString(_header)
		sl = C.curl_slist_append(sl, header)
		C.free(unsafe.Pointer(header))
	}
	if sl == nil {
		return
	}
	C.yam_curl_inject_trace_header(C.int(id), sl)
}

type redisSocket struct {
	Stream unsafe.Pointer
	Host   unsafe.Pointer
	Port   uint16
	Auth   unsafe.Pointer
}

func redisAddr(obj unsafe.Pointer) (host string, port int, ok bool) {
	redisVersion := C.GoString((*C.char)(C.yam_redis_version()))
	if len(redisVersion) < 1 {
		return
	}
	v, err := version.NewVersion(redisVersion)
	if err != nil {
		return
	}
	var p unsafe.Pointer
	if v.GreaterThanOrEqual(v200) && v.LessThanOrEqual(v228) {
		p = C.yam_get_redis_socket_gte200_lte228(obj)
		socket := (*redisSocket)(p)
		host, port, ok = C.GoString((*C.char)(socket.Host)), int(socket.Port), true
	} else if v.GreaterThanOrEqual(v400) && v.LessThanOrEqual(v430) {
		p = C.yam_get_redis_socket_gte400_lte430(obj)
		socket := (*redisSocket)(p)
		host, port, ok = ((*zend.String)(socket.Host)).String(), int(socket.Port), true
	} else {
		return
	}
	return
}
