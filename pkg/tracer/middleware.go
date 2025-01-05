package tracer

import (
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strings"
)

func TcpMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parentTrace := strings.Split(r.Header.Get(headerB3), "-")
		ctx := r.Context()

		if len(parentTrace) > 2 {
			traceID, err := trace.TraceIDFromHex(parentTrace[0])
			if err == nil {
				var traceFlags trace.TraceFlags
				if parentTrace[2] == samplingState1 {
					traceFlags = trace.FlagsSampled
				}
				spanContext := trace.NewSpanContext(trace.SpanContextConfig{
					TraceID:    traceID,
					TraceFlags: traceFlags,
				})

				ctx = trace.ContextWithSpanContext(ctx, spanContext)
			}
		}

		ctx, span := StartTrace(ctx, r.Method+" "+r.URL.Path)
		defer span.End()
		span.SetStatus(codes.Ok, "http request")
		span.SetAttributes(NewHttpAttributes(tcp, r)...)

		h.ServeHTTP(w, r.WithContext(ctx))
		return
	})
}
