package sdks

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

type (
	spanKey struct{}
	idKey   struct{}
)

func SpanWithContext(ctx context.Context, span opentracing.Span) context.Context {
	return context.WithValue(ctx, spanKey{}, span)
}

func SpanFromContext(ctx context.Context) opentracing.Span {
	return ctx.Value(spanKey{}).(opentracing.Span)
}

func WithContext(ctx context.Context, id string) (nCtx context.Context) {
	if len(id) == 0 {
		id = NewIDFunc()
	}
	nCtx = context.WithValue(ctx, idKey{}, id)
	return
}

func GetTraceIDFromContext(ctx context.Context) string {
	id, ok := ctx.Value(idKey{}).(string)
	if !ok {
		return ""
	}
	return id
}
