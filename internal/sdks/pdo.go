package sdks

import (
	"context"
	"time"

	"github.com/Lofanmi/yam/api"
	"github.com/opentracing/opentracing-go"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

type pdoContext struct {
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

func NewPdoContext() api.DB {
	return &pdoContext{}
}

func (p *pdoContext) Begin() (t time.Time) {
	return p.begin
}

func (p *pdoContext) End() (t time.Time) {
	return p.end
}

func (p *pdoContext) Host() string {
	return p.host
}

func (p *pdoContext) Port() string {
	return p.port
}

func (p *pdoContext) Addr() string {
	return p.host + ":" + p.port
}

func (p *pdoContext) Operation() string {
	return p.operation
}

func (p *pdoContext) Statement() string {
	return p.statement
}

func (p *pdoContext) User() string {
	return p.user
}

func (p *pdoContext) Database() string {
	return p.database
}

func (p *pdoContext) SetBegin(t time.Time) {
	p.begin = t
}

func (p *pdoContext) SetEnd(t time.Time) {
	p.end = t
}

func (p *pdoContext) SetHost(v string) {
	p.host = v
}

func (p *pdoContext) SetPort(v string) {
	p.port = v
}

func (p *pdoContext) SetOperation(v string) {
	p.operation = v
}

func (p *pdoContext) SetStatement(v string) {
	p.statement = v
}

func (p *pdoContext) SetUser(v string) {
	p.user = v
}

func (p *pdoContext) SetDatabase(v string) {
	p.database = v
}

func (p *pdoContext) BeginTrace(request api.PHPRequest, id uintptr) {
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

func (p *pdoContext) EndTrace(request api.PHPRequest, id uintptr) {
	span := SpanFromContext(p.ctx)
	defer span.Finish()
	span.SetTag(string(semconv.DBSystemKey), semconv.DBSystemMySQL.Value.AsString())
	span.SetTag(string(semconv.DBConnectionStringKey), p.Addr())
	span.SetTag(string(semconv.DBOperationKey), p.Operation())
	span.SetTag(string(semconv.DBStatementKey), p.Statement())
	span.SetTag("db.type", "php:pdo")
}
