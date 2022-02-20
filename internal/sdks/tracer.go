package sdks

import (
	"github.com/Lofanmi/yam/api"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/log"
)

func NewTracerNoop(c api.TracerConfig) (opentracing.Tracer, func()) {
	return opentracing.NoopTracer{}, func() {}
}

func NewTracer(c api.TracerConfig) (opentracing.Tracer, func()) {
	cfg := &config.Configuration{
		ServiceName: c.ServiceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "127.0.0.1:6831",
		},
	}
	tracer, closer, err := cfg.NewTracer(
		config.Logger(log.NullLogger),
		config.MaxTagValueLength(4096),
	)
	// tracer, closer, err := cfg.NewTracer(config.Logger(log.StdLogger))
	if err != nil {
		panic(err)
	}
	return tracer, func() { _ = closer.Close() }
}
