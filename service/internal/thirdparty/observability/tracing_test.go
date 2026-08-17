package observability

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestTracingMiddlewareContinuesW3CTrace(t *testing.T) {
	gin.SetMode(gin.TestMode)
	oldProvider := otel.GetTracerProvider()
	oldPropagator := otel.GetTextMapPropagator()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	defer func() {
		_ = provider.Shutdown(context.Background())
		otel.SetTracerProvider(oldProvider)
		otel.SetTextMapPropagator(oldPropagator)
	}()

	router := gin.New()
	router.Use((&Tracing{enabled: true, provider: provider}).Middleware())
	router.GET("/trace", func(ctx *gin.Context) { ctx.Status(http.StatusNoContent) })

	request := httptest.NewRequest(http.MethodGet, "/trace", nil)
	request.Header.Set("traceparent", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if got := recorder.Header().Get("X-Trace-ID"); got != "4bf92f3577b34da6a3ce929d0e0e4736" {
		t.Fatalf("X-Trace-ID = %q", got)
	}
}
