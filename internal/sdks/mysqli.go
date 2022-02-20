package sdks

import (
	"context"
	"time"

	"github.com/Lofanmi/yam/api"
	"github.com/opentracing/opentracing-go"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

type mysqliContext struct {
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

func NewMySQLiContext() api.DB {
	return &mysqliContext{}
}

func (m *mysqliContext) Begin() (t time.Time) {
	return m.begin
}

func (m *mysqliContext) End() (t time.Time) {
	return m.end
}

func (m *mysqliContext) Host() string {
	return m.host
}

func (m *mysqliContext) Port() string {
	return m.port
}

func (m *mysqliContext) Addr() string {
	return m.host + ":" + m.port
}

func (m *mysqliContext) Operation() string {
	return m.operation
}

func (m *mysqliContext) Statement() string {
	return m.statement
}

func (m *mysqliContext) User() string {
	return m.user
}

func (m *mysqliContext) Database() string {
	return m.database
}

func (m *mysqliContext) SetBegin(t time.Time) {
	m.begin = t
}

func (m *mysqliContext) SetEnd(t time.Time) {
	m.end = t
}

func (m *mysqliContext) SetHost(v string) {
	m.host = v
}

func (m *mysqliContext) SetPort(v string) {
	m.port = v
}

func (m *mysqliContext) SetOperation(v string) {
	m.operation = v
}

func (m *mysqliContext) SetStatement(v string) {
	m.statement = v
}

func (m *mysqliContext) SetUser(v string) {
	m.user = v
}

func (m *mysqliContext) SetDatabase(v string) {
	m.database = v
}

func (m *mysqliContext) BeginTrace(request api.PHPRequest, id uintptr) {
	ctx := request.GetContext()
	span, ctx := opentracing.StartSpanFromContext(ctx, m.Addr())
	xRequestID := GetTraceIDFromContext(ctx)
	if len(xRequestID) > 0 {
		span.SetTag(HeaderXRequestID, xRequestID)
		if GetTraceIDFromSpanContext(span.Context()) != xRequestID {
			UpdateSpanTraceID(span, xRequestID)
		}
	}
	m.ctx = SpanWithContext(ctx, span)
}

func (m *mysqliContext) EndTrace(request api.PHPRequest, id uintptr) {
	span := SpanFromContext(m.ctx)
	defer span.Finish()
	span.SetTag(string(semconv.DBSystemKey), semconv.DBSystemMySQL.Value.AsString())
	span.SetTag(string(semconv.DBConnectionStringKey), m.Addr())
	span.SetTag(string(semconv.DBOperationKey), m.Operation())
	span.SetTag(string(semconv.DBStatementKey), m.Statement())
	span.SetTag("db.type", "php:mysqli")
}
