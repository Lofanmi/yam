package sdks

import (
	"net/textproto"
	"reflect"
	"unsafe"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

const (
	HeaderXRequestID = "X-Request-Id"
	HeaderXTraceID   = "X-Trace-Id"
)

var DefaultReqIDKeys = []string{HeaderXRequestID, HeaderXTraceID}

func SetMetadataRequestID(requestID string, metadata map[string][]string, extReqIDKeys ...string) {
	if len(requestID) == 0 {
		return
	}
	if len(extReqIDKeys) == 0 {
		extReqIDKeys = DefaultReqIDKeys
	}
	for _, key := range extReqIDKeys {
		textproto.MIMEHeader(metadata).Set(key, requestID)
	}
}

func GetRequestIDFromMetadata(metadata map[string][]string, extReqIDKeys ...string) string {
	if len(extReqIDKeys) == 0 {
		extReqIDKeys = DefaultReqIDKeys
	}
	var requestID string
	for _, key := range extReqIDKeys {
		if requestID = textproto.MIMEHeader(metadata).Get(key); len(requestID) > 0 {
			break
		}
	}
	return requestID
}

func GetTraceIDFromSpanContext(ctx opentracing.SpanContext) string {
	switch c := ctx.(type) {
	case jaeger.SpanContext:
		return c.TraceID().String()
	}
	return ""
}

// UpdateSpanTraceID
// 设置自定义traceID到context 成功返回 自定义ID 否则返回Span里的traceID
func UpdateSpanTraceID(span opentracing.Span, traceID string) string {
	if len(traceID) == 0 {
		return GetTraceIDFromSpanContext(span.Context())
	}
	var set bool
	switch sp := span.(type) {
	case *jaeger.Span:
		id, err := jaeger.TraceIDFromString(traceID)
		if err != nil {
			break
		}
		sp.Lock()
		setJaegerSpanContextTraceID(sp, id)
		sp.Unlock()
		set = true
	}
	if set {
		return traceID
	}
	return GetTraceIDFromSpanContext(span.Context())
}

func setJaegerSpanContextTraceID(sp *jaeger.Span, traceID jaeger.TraceID) {
	// traceIDPtr := reflect.ValueOf(sp).Elem().FieldByName("context").FieldByName("traceID").UnsafeAddr()
	// *(*jaeger.TraceID)(unsafe.Pointer(traceIDPtr)) = traceID

	// 写成同一行就不会有 warning 提示。
	*(*jaeger.TraceID)(unsafe.Pointer(reflect.ValueOf(sp).Elem().FieldByName("context").FieldByName("traceID").UnsafeAddr())) = traceID
}
