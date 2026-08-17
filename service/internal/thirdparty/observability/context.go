package observability

import (
	"context"

	"github.com/james730922/wallet/service/internal/utils/conf"
)

func contextWithReadinessTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), conf.Observability().GetReadinessTimeout())
}
