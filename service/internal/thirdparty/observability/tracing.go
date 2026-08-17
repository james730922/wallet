package observability

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/james730922/wallet/service/internal/utils/conf"
)

type Tracing struct {
	enabled  bool
	provider *sdktrace.TracerProvider
}

func NewTracing() (*Tracing, error) {
	config := conf.Observability()
	tracing := &Tracing{enabled: config.GetTracingEnabled()}
	if !tracing.enabled {
		return tracing, nil
	}

	providerOptions := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(resource.NewSchemaless(
			attribute.String("service.name", config.GetTracingServiceName()),
		)),
		sdktrace.WithSampler(sdktrace.ParentBased(
			sdktrace.TraceIDRatioBased(config.GetTracingSampleRatio()),
		)),
	}
	if endpoint := config.GetOTLPEndpoint(); endpoint != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(endpoint))
		if err != nil {
			return nil, err
		}
		providerOptions = append(providerOptions, sdktrace.WithBatcher(exporter))
	}

	tracing.provider = sdktrace.NewTracerProvider(providerOptions...)
	otel.SetTracerProvider(tracing.provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return tracing, nil
}

func (t *Tracing) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !t.enabled {
			ctx.Next()
			return
		}

		parent := otel.GetTextMapPropagator().Extract(
			ctx.Request.Context(),
			propagation.HeaderCarrier(ctx.Request.Header),
		)
		spanContext, span := otel.Tracer("wallet/http").Start(
			parent,
			ctx.Request.Method+" "+ctx.Request.URL.Path,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.request.method", ctx.Request.Method),
				attribute.String("url.path", ctx.Request.URL.Path),
			),
		)
		defer span.End()
		ctx.Request = ctx.Request.WithContext(spanContext)
		if traceID := span.SpanContext().TraceID(); traceID.IsValid() {
			ctx.Header("X-Trace-ID", traceID.String())
		}

		ctx.Next()
		route := ctx.FullPath()
		if route == "" {
			route = "unmatched"
		}
		span.SetName(ctx.Request.Method + " " + route)
		span.SetAttributes(
			attribute.String("http.route", route),
			attribute.Int("http.response.status_code", ctx.Writer.Status()),
		)
		if ctx.Writer.Status() >= http.StatusInternalServerError {
			span.SetStatus(codes.Error, http.StatusText(ctx.Writer.Status()))
		}
	}
}

func (t *Tracing) Shutdown(ctx context.Context) error {
	if t.provider == nil {
		return nil
	}
	return t.provider.Shutdown(ctx)
}
