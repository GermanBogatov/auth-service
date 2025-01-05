package tracer

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	tr "go.opentelemetry.io/otel/trace"
	"log"
	"net/url"
	"time"
)

func New(cfg *Config) (func(ctx context.Context), error) {
	urlEndpoint, err := url.Parse(cfg.Endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "fail parse endpoint")
	}

	var (
		exporter    *otlptrace.Exporter
		errExporter error
	)

	switch urlEndpoint.Scheme {
	case "http":
		exporter, errExporter = otlptrace.New(context.Background(), otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(urlEndpoint.Host),
			otlptracehttp.WithInsecure(),
		))
		if errExporter != nil {
			return nil, errors.Wrap(err, "fail init http exporter")
		}
	case "https":
		exporter, errExporter = otlptrace.New(context.Background(), otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(urlEndpoint.Host),
		))
		if errExporter != nil {
			return nil, errors.Wrap(err, "fail init https exporter")
		}
	default:
		return nil, fmt.Errorf("the URL scheme is not recognized. Pass the endpoint [%s] with http:// or https:// scheme", cfg.Endpoint)
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.TraceRatioFraction)),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			attribute.String("environment", cfg.Environment),
		)),
	)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(traceProvider)

	// Return func for graceful shutdown tracer
	return func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := traceProvider.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}, nil
}

func StartTrace(ctx context.Context, spanName string) (context.Context, tr.Span) {
	tp := otel.GetTracerProvider()
	t := tp.Tracer("")
	return t.Start(ctx, spanName)
}
