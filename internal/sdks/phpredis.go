package sdks

import (
	"context"
	"time"

	"github.com/Lofanmi/yam/api"
	"github.com/opentracing/opentracing-go"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

type phpRedisContext struct {
	ctx       context.Context
	begin     time.Time
	end       time.Time
	host      string
	port      string
	operation string
	statement string
	user      string
	database  string
}

func NewPHPRedisContext() api.DB {
	return &phpRedisContext{}
}

func (p *phpRedisContext) Begin() (t time.Time) {
	return p.begin
}

func (p *phpRedisContext) End() (t time.Time) {
	return p.end
}

func (p *phpRedisContext) Host() string {
	return p.host
}

func (p *phpRedisContext) Port() string {
	return p.port
}

func (p *phpRedisContext) Addr() string {
	if p.host != "" && p.port != "" {
		return p.host + ":" + p.port
	}
	if p.host != "" {
		return p.host
	}
	return ""
}

func (p *phpRedisContext) Operation() string {
	return p.operation
}

func (p *phpRedisContext) Statement() string {
	return p.statement
}

func (m *phpRedisContext) User() string {
	return m.user
}

func (m *phpRedisContext) Database() string {
	return m.database
}

func (p *phpRedisContext) SetBegin(t time.Time) {
	p.begin = t
}

func (p *phpRedisContext) SetEnd(t time.Time) {
	p.end = t
}

func (p *phpRedisContext) SetHost(v string) {
	p.host = v
}

func (p *phpRedisContext) SetPort(v string) {
	p.port = v
}

func (p *phpRedisContext) SetOperation(v string) {
	p.operation = v
}

func (p *phpRedisContext) SetStatement(v string) {
	p.statement = v
}

func (m *phpRedisContext) SetUser(v string) {
	m.user = v
}

func (m *phpRedisContext) SetDatabase(v string) {
	m.database = v
}

func (p *phpRedisContext) BeginTrace(request api.PHPRequest, id uintptr) {
	ctx := request.GetContext()
	span, ctx := opentracing.StartSpanFromContext(ctx, p.Addr())
	xRequestID := GetTraceIDFromContext(ctx)
	if len(xRequestID) > 0 {
		span.SetTag(HeaderXRequestID, xRequestID)
		if GetTraceIDFromSpanContext(span.Context()) != xRequestID {
			UpdateSpanTraceID(span, xRequestID)
		}
	}
	p.ctx = SpanWithContext(ctx, span)
}

func (p *phpRedisContext) EndTrace(request api.PHPRequest, id uintptr) {
	span := SpanFromContext(p.ctx)
	if span == nil {
		return
	}
	defer span.Finish()
	span.SetTag(string(semconv.DBSystemKey), semconv.DBSystemRedis.Value.AsString())
	span.SetTag(string(semconv.DBConnectionStringKey), p.Addr())
	span.SetTag(string(semconv.DBOperationKey), p.Operation())
	span.SetTag(string(semconv.DBStatementKey), p.Statement())
}
