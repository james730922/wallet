package observability

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/james730922/wallet/service/internal/utils/errs"
)

func TestMetricsMiddlewareRecordsRequestErrorAndLatency(t *testing.T) {
	gin.SetMode(gin.TestMode)
	m := &Metrics{}
	router := gin.New()
	router.Use(m.Middleware())
	router.GET("/failed", func(ctx *gin.Context) { ctx.Status(http.StatusServiceUnavailable) })

	requestBefore := testutil.ToFloat64(requestTotal.WithLabelValues(http.MethodGet, "/failed", "503"))
	errorBefore := testutil.ToFloat64(requestErrors.WithLabelValues(http.MethodGet, "/failed", "503"))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/failed", nil))

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", recorder.Code)
	}
	if got := testutil.ToFloat64(requestTotal.WithLabelValues(http.MethodGet, "/failed", "503")); got != requestBefore+1 {
		t.Fatalf("request counter = %v, want %v", got, requestBefore+1)
	}
	if got := testutil.ToFloat64(requestErrors.WithLabelValues(http.MethodGet, "/failed", "503")); got != errorBefore+1 {
		t.Fatalf("error counter = %v, want %v", got, errorBefore+1)
	}
}

func TestErrorClassUsesBoundedLabels(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{err: errs.WalletMemberUpdateBalanceIsNegative, want: "insufficient_balance"},
		{err: errs.OrderScanPaySignValidateFailed, want: "integrity"},
		{err: errs.DBOperationFailed, want: "database"},
		{err: errs.CommonServiceUnavailable, want: "system"},
	}
	for _, tt := range tests {
		if got := errorClass(tt.err); got != tt.want {
			t.Fatalf("errorClass(%v) = %q, want %q", tt.err, got, tt.want)
		}
	}
}

func TestReadinessStateDoesNotExposeErrors(t *testing.T) {
	if readinessState(true) != "ready" || readinessState(false) != "unavailable" {
		t.Fatal("unexpected readiness state")
	}
}
