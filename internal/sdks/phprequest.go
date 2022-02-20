package sdks

import "C"
import (
	"context"
	"net/http"
	"os"
	"strings"
	"unsafe"

	"github.com/Lofanmi/yam/api"
	"github.com/Lofanmi/yam/internal/zend"
	"github.com/opentracing/opentracing-go"
)

var _ api.PHPRequest = (*phpRequest)(nil)

type phpRequest struct {
	ctx                 context.Context
	serviceName         string
	tracerCloser        func()
	curls               map[int32]api.Curl
	phpRedis            map[uintptr]api.DB
	mysqli              map[uintptr]api.DB
	pdo                 map[uintptr]api.DB
	mysqliStmtMapping   map[uintptr]uintptr
	pdoStatementMapping map[uintptr]uintptr
}

func NewPHPRequest() api.PHPRequest {
	r := &phpRequest{
		curls:               make(map[int32]api.Curl),
		phpRedis:            make(map[uintptr]api.DB),
		mysqli:              make(map[uintptr]api.DB),
		pdo:                 make(map[uintptr]api.DB),
		mysqliStmtMapping:   make(map[uintptr]uintptr),
		pdoStatementMapping: make(map[uintptr]uintptr),
	}
	r.initTracer()
	return r
}

func (r *phpRequest) initTracer() {
	var tracer opentracing.Tracer
	// tracer, r.tracerCloser = NewTracerNoop(api.TracerConfig{ServiceName: r.ServiceName()})
	tracer, r.tracerCloser = NewTracer(api.TracerConfig{ServiceName: r.ServiceName()})
	opentracing.SetGlobalTracer(tracer)
}

func (r *phpRequest) GetContext() context.Context {
	return r.ctx
}

func (r *phpRequest) SetContext(ctx context.Context) {
	r.ctx = ctx
}

func (r *phpRequest) ServiceName() string {
	if r.serviceName == "" {
		host := r.GlobalByKey(api.GlobalTrackVarsServer, "HTTP_HOST")
		if host != "" {
			r.serviceName = host
		} else {
			separator := string(os.PathSeparator)
			scriptFilename := strings.TrimLeft(r.GlobalByKey(api.GlobalTrackVarsServer, "SCRIPT_FILENAME"), separator)
			r.serviceName = strings.ReplaceAll(strings.TrimSuffix(scriptFilename, ".php"), separator, "-")
		}
	}
	if r.serviceName == "" {
		hostname, _ := os.Hostname()
		r.serviceName = hostname + "(unknown-service)"
	}
	return r.serviceName
}

func (r *phpRequest) Close() {
	r.tracerCloser()
}

func (r *phpRequest) Globals(typ api.GlobalTrackVars) (m map[string]string) {
	api.Globals(typ, func(ret unsafe.Pointer) {
		zv := zend.NewZVal(ret)
		m = zv.AsArray().ToSSMap()
	})
	return
}

func (r *phpRequest) GlobalByKey(typ api.GlobalTrackVars, key string) (val string) {
	api.GlobalByKey(typ, key, func(s unsafe.Pointer) {
		val = C.GoString((*C.char)(s))
	})
	return
}

func (r *phpRequest) Headers() (header http.Header) {
	header = make(http.Header)
	for k, val := range r.Globals(api.GlobalTrackVarsServer) {
		if strings.HasPrefix(k, "HTTP_") {
			key := strings.ReplaceAll(strings.TrimPrefix(k, "HTTP_"), "_", "-")
			header.Set(key, val)
		}
	}
	return
}

func (r *phpRequest) Curl(id int32) api.Curl {
	if _, ok := r.curls[id]; !ok {
		r.curls[id] = NewCurlContext()
	}
	return r.curls[id]
}

func (r *phpRequest) PHPRedis(id uintptr) api.DB {
	if _, ok := r.phpRedis[id]; !ok {
		r.phpRedis[id] = NewPHPRedisContext()
	}
	return r.phpRedis[id]
}

func (r *phpRequest) MySQLi(id uintptr) api.DB {
	if _, ok := r.mysqli[id]; !ok {
		r.mysqli[id] = NewMySQLiContext()
	}
	return r.mysqli[id]
}

func (r *phpRequest) Pdo(id uintptr) api.DB {
	if _, ok := r.pdo[id]; !ok {
		r.pdo[id] = NewPdoContext()
	}
	return r.pdo[id]
}

func (r *phpRequest) MySQLiStmtMapping() map[uintptr]uintptr {
	return r.mysqliStmtMapping
}

func (r *phpRequest) PdoStatementMapping() map[uintptr]uintptr {
	return r.pdoStatementMapping
}
