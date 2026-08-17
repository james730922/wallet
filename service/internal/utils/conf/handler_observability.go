package conf

import (
	"os"
	"strings"
	"time"
)

type ObservabilityHandler struct{}

func (ObservabilityHandler) GetReadinessTimeout() time.Duration {
	milliseconds := zqbConf.v.GetInt("observability.readiness_timeout_ms")
	if milliseconds <= 0 {
		milliseconds = 1000
	}
	return time.Duration(milliseconds) * time.Millisecond
}

func (ObservabilityHandler) GetTracingEnabled() bool {
	return zqbConf.v.GetBool("observability.tracing.enable")
}

func (ObservabilityHandler) GetTracingServiceName() string {
	name := strings.TrimSpace(zqbConf.v.GetString("observability.tracing.service_name"))
	if name == "" {
		return "wallet-api"
	}
	return name
}

func (ObservabilityHandler) GetTracingSampleRatio() float64 {
	ratio := zqbConf.v.GetFloat64("observability.tracing.sample_ratio")
	if ratio < 0 || ratio > 1 {
		return 1
	}
	return ratio
}

func (ObservabilityHandler) GetOTLPEndpoint() string {
	if endpoint := strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")); endpoint != "" {
		return endpoint
	}
	if endpoint := strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")); endpoint != "" {
		return strings.TrimRight(endpoint, "/") + "/v1/traces"
	}
	return strings.TrimSpace(zqbConf.v.GetString("observability.tracing.otlp_http_endpoint"))
}
