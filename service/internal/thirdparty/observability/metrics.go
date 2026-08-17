package observability

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/dig"

	"github.com/james730922/wallet/service/internal/utils/errs"
)

var (
	requestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "wallet",
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total HTTP requests handled by route and status.",
	}, []string{"method", "route", "status"})
	requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "wallet",
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "HTTP request latency by route.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "route"})
	requestErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "wallet",
		Subsystem: "http",
		Name:      "errors_total",
		Help:      "HTTP error responses by route and status.",
	}, []string{"method", "route", "status"})
	requestInFlight = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "wallet",
		Subsystem: "http",
		Name:      "requests_in_flight",
		Help:      "Current number of HTTP requests being served.",
	})
	walletTransactionFailures = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "wallet",
		Subsystem: "transaction",
		Name:      "failures_total",
		Help:      "Wallet transaction failures by operation and bounded reason.",
	}, []string{"operation", "reason"})
	scanPayTransactionFailures = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "wallet",
		Subsystem: "scanpay",
		Name:      "transaction_failures_total",
		Help:      "ScanPay transaction failures by stage and bounded reason.",
	}, []string{"stage", "reason"})
)

type metricsSet struct {
	dig.In

	DB      *gorm.DB
	DBSlave *gorm.DB `name:"dbSlave"`
	Redis   *redis.Client
}

type Metrics struct {
	registry  *prometheus.Registry
	db        *gorm.DB
	dbSlave   *gorm.DB
	redis     *redis.Client
	dbPool    *prometheus.GaugeVec
	dbWait    *prometheus.GaugeVec
	redisPool *prometheus.GaugeVec
}

func NewMetrics(set metricsSet) *Metrics {
	m := &Metrics{
		registry: prometheus.NewRegistry(),
		db:       set.DB,
		dbSlave:  set.DBSlave,
		redis:    set.Redis,
		dbPool: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "wallet", Subsystem: "db", Name: "pool_connections",
			Help: "Database pool connections by role and state.",
		}, []string{"role", "state"}),
		dbWait: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "wallet", Subsystem: "db", Name: "pool_wait",
			Help: "Database pool cumulative wait count and seconds.",
		}, []string{"role", "metric"}),
		redisPool: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "wallet", Subsystem: "redis", Name: "pool",
			Help: "Redis pool connections and cumulative events.",
		}, []string{"metric"}),
	}
	m.registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		requestTotal,
		requestDuration,
		requestErrors,
		requestInFlight,
		walletTransactionFailures,
		scanPayTransactionFailures,
		m.dbPool,
		m.dbWait,
		m.redisPool,
	)
	return m
}

func (m *Metrics) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/metrics" {
			ctx.Next()
			return
		}

		started := time.Now()
		requestInFlight.Inc()
		defer requestInFlight.Dec()
		ctx.Next()

		route := ctx.FullPath()
		if route == "" {
			route = "unmatched"
		}
		status := strconv.Itoa(ctx.Writer.Status())
		requestTotal.WithLabelValues(ctx.Request.Method, route, status).Inc()
		requestDuration.WithLabelValues(ctx.Request.Method, route).Observe(time.Since(started).Seconds())
		if ctx.Writer.Status() >= http.StatusBadRequest {
			requestErrors.WithLabelValues(ctx.Request.Method, route, status).Inc()
		}
	}
}

func (m *Metrics) Handler() http.Handler {
	handler := promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.updatePoolMetrics()
		handler.ServeHTTP(w, r)
	})
}

func (m *Metrics) Livez(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "alive"})
}

func (m *Metrics) Readyz(ctx *gin.Context) {
	dbReady := m.pingDB()
	redisReady := m.pingRedis()
	status := http.StatusOK
	state := "ready"
	if !dbReady || !redisReady {
		status = http.StatusServiceUnavailable
		state = "not_ready"
	}
	ctx.JSON(status, gin.H{
		"status": state,
		"dependencies": gin.H{
			"db":    readinessState(dbReady),
			"redis": readinessState(redisReady),
		},
	})
}

func (m *Metrics) pingDB() bool {
	ctx, cancel := contextWithReadinessTimeout()
	defer cancel()
	return m.db.DB().PingContext(ctx) == nil
}

func (m *Metrics) pingRedis() bool {
	ctx, cancel := contextWithReadinessTimeout()
	defer cancel()
	return m.redis.WithContext(ctx).Ping().Err() == nil
}

func (m *Metrics) updatePoolMetrics() {
	m.updateDBPool("primary", m.db)
	m.updateDBPool("replica", m.dbSlave)
	stats := m.redis.PoolStats()
	for metric, value := range map[string]uint32{
		"hits_total": stats.Hits, "misses_total": stats.Misses, "timeouts_total": stats.Timeouts,
		"connections": stats.TotalConns, "idle_connections": stats.IdleConns, "stale_connections_total": stats.StaleConns,
	} {
		m.redisPool.WithLabelValues(metric).Set(float64(value))
	}
}

func (m *Metrics) updateDBPool(role string, db *gorm.DB) {
	stats := db.DB().Stats()
	for state, value := range map[string]int{
		"max_open": stats.MaxOpenConnections,
		"open":     stats.OpenConnections,
		"in_use":   stats.InUse,
		"idle":     stats.Idle,
	} {
		m.dbPool.WithLabelValues(role, state).Set(float64(value))
	}
	m.dbWait.WithLabelValues(role, "count").Set(float64(stats.WaitCount))
	m.dbWait.WithLabelValues(role, "seconds").Set(stats.WaitDuration.Seconds())
}

func readinessState(ready bool) string {
	if ready {
		return "ready"
	}
	return "unavailable"
}

func RecordWalletTransactionFailure(operation string, err error) {
	walletTransactionFailures.WithLabelValues(operation, errorClass(err)).Inc()
}

func RecordScanPayTransactionFailure(stage string, err error) {
	scanPayTransactionFailures.WithLabelValues(stage, errorClass(err)).Inc()
}

func errorClass(err error) string {
	switch {
	case errors.Is(err, errs.WalletMemberUpdateBalanceIsNegative):
		return "insufficient_balance"
	case errors.Is(err, errs.CommonAmountDecimalPlacesError), errors.Is(err, errs.CommonRequestParamInvalid):
		return "validation"
	case errors.Is(err, errs.WalletMemberSignValidateFailed), errors.Is(err, errs.OrderScanPaySignValidateFailed):
		return "integrity"
	case errors.Is(err, errs.DBOperationFailed), errors.Is(err, errs.DBUpdateFailed), errors.Is(err, errs.DBInsertFailed):
		return "database"
	default:
		return "system"
	}
}
