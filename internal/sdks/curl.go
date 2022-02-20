package sdks

import (
	"context"
	"net/http"
	"time"

	"github.com/Lofanmi/yam/api"
	"github.com/opentracing/opentracing-go"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

const (
	defaultSchemeHttp = "http"
	defaultPort       = "80"
	defaultPath       = "/"
)

type curlContext struct {
	ctx               context.Context
	begin             time.Time
	end               time.Time
	url               string
	method            string
	scheme            string
	host              string
	port              string
	path              string
	query             string
	body              string
	bodySize          int
	contentType       string
	httpCode          int
	httpVersion       string
	totalTime         float64
	nameLookupTime    float64
	connectTime       float64
	preTransferTime   float64
	startTransferTime float64
	remoteIP          string
	remotePort        int
	headers           []string
}

func NewCurlContext() api.Curl {
	return &curlContext{
		method:      http.MethodGet,
		scheme:      defaultSchemeHttp,
		port:        defaultPort,
		path:        defaultPath,
		httpVersion: semconv.HTTPFlavorHTTP10.Value.AsString(),
		headers:     make([]string, 0, 8),
	}
}

func (c *curlContext) Begin() time.Time {
	return c.begin
}

func (c *curlContext) End() time.Time {
	return c.end
}

func (c *curlContext) Url() string {
	return c.url
}

func (c *curlContext) Method() string {
	return c.method
}

func (c *curlContext) Scheme() string {
	return c.scheme
}

func (c *curlContext) Host() string {
	return c.host
}

func (c *curlContext) Port() string {
	return c.port
}

func (c *curlContext) Path() string {
	return c.path
}

func (c *curlContext) Query() string {
	return c.query
}

func (c *curlContext) Body() string {
	return c.body
}

func (c *curlContext) BodySize() int {
	return c.bodySize
}

func (c *curlContext) ContentType() string {
	return c.contentType
}

func (c *curlContext) HttpCode() int {
	return c.httpCode
}

func (c *curlContext) HttpVersion() string {
	return c.httpVersion
}

func (c *curlContext) TotalTime() float64 {
	return c.totalTime
}

func (c *curlContext) NameLookupTime() float64 {
	return c.nameLookupTime
}

func (c *curlContext) ConnectTime() float64 {
	return c.connectTime
}

func (c *curlContext) PreTransferTime() float64 {
	return c.preTransferTime
}

func (c *curlContext) StartTransferTime() float64 {
	return c.startTransferTime
}

func (c *curlContext) RemoteIP() string {
	return c.remoteIP
}

func (c *curlContext) RemotePort() int {
	return c.remotePort
}

func (c *curlContext) Headers() []string {
	return c.headers
}

func (c *curlContext) SetBegin(v time.Time) {
	c.begin = v
}

func (c *curlContext) SetEnd(v time.Time) {
	c.end = v
}

func (c *curlContext) SetMethod(v string) {
	c.method = v
}

func (c *curlContext) SetUrl(v string) {
	c.url = v
}

func (c *curlContext) SetScheme(v string) {
	if v == "" {
		c.scheme = defaultSchemeHttp
	} else {
		c.scheme = v
	}
}

func (c *curlContext) SetHost(v string) {
	c.host = v
}

func (c *curlContext) SetPort(v string) {
	if v == "" {
		c.port = defaultPort
	} else {
		c.port = v
	}
}

func (c *curlContext) SetPath(v string) {
	if v == "" {
		c.path = defaultPath
	} else {
		c.path = v
	}
}

func (c *curlContext) SetQuery(v string) {
	c.query = v
}

func (c *curlContext) SetBody(v string) {
	c.body = v
}

func (c *curlContext) SetBodySize(v int) {
	c.bodySize = v
}

func (c *curlContext) SetContentType(v string) {
	c.contentType = v
}

func (c *curlContext) SetHttpCode(v int) {
	c.httpCode = v
}

func (c *curlContext) SetHttpVersion(v string) {
	c.httpVersion = v
}

func (c *curlContext) SetTotalTime(v float64) {
	c.totalTime = v
}

func (c *curlContext) SetNameLookupTime(v float64) {
	c.nameLookupTime = v
}

func (c *curlContext) SetConnectTime(v float64) {
	c.connectTime = v
}

func (c *curlContext) SetPreTransferTime(v float64) {
	c.preTransferTime = v
}

func (c *curlContext) SetStartTransferTime(v float64) {
	c.startTransferTime = v
}

func (c *curlContext) SetRemoteIP(v string) {
	c.remoteIP = v
}

func (c *curlContext) SetRemotePort(v int) {
	c.remotePort = v
}

func (c *curlContext) AppendHeader(header string) {
	c.headers = append(c.headers, header)
}

func (c *curlContext) BeginTrace(request api.PHPRequest, id int32) {
	ctx := request.GetContext()
	span, ctx := opentracing.StartSpanFromContext(ctx, c.Host())
	xRequestID := GetTraceIDFromContext(ctx)
	if len(xRequestID) > 0 {
		span.SetTag(HeaderXRequestID, xRequestID)
		if GetTraceIDFromSpanContext(span.Context()) != xRequestID {
			UpdateSpanTraceID(span, xRequestID)
		}
	}
	metadata := make(http.Header)
	SetMetadataRequestID(xRequestID, metadata)
	if err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(metadata)); err == nil {
		for key, _ := range metadata {
			c.AppendHeader(key + ": " + metadata.Get(key))
		}
		api.CurlInjectTraceHeader(id, c.Headers())
	}
	c.ctx = SpanWithContext(ctx, span)
}

func (c *curlContext) EndTrace(request api.PHPRequest, id int32) {
	span := SpanFromContext(c.ctx)
	defer span.Finish()
	span.SetTag(string(semconv.HTTPMethodKey), c.Method())
	span.SetTag(string(semconv.HTTPURLKey), c.Url())
	span.SetTag(string(semconv.HTTPHostKey), c.Host())
	span.SetTag(string(semconv.HTTPSchemeKey), c.Scheme())
	span.SetTag(string(semconv.HTTPStatusCodeKey), c.HttpCode())
	span.SetTag(string(semconv.HTTPFlavorKey), c.HttpVersion())
	span.SetTag(string(semconv.NetPeerIPKey), c.RemoteIP())
	switch c.Method() {
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		fallthrough
	case http.MethodPatch:
		span.SetTag(string(semconv.HTTPRequestContentLengthUncompressedKey), c.BodySize())
	}
}
