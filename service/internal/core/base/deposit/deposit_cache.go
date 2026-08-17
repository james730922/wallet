package deposit

import (
	"fmt"
	"time"

	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
)

const depositDedupeTTL = 30 * time.Second

const releaseDepositClaimScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
end
return 0`

func newDepositCache() iDepositCache {
	return &depositCacheUseCase{}
}

type iDepositCache interface {
	claimDepositWithSameMemberAndAmount(cond *model.Deposit, owner string) (bool, error)
	releaseDepositWithSameMemberAndAmount(cond *model.Deposit, owner string)
}

type depositCacheUseCase struct{}

func (dc *depositCacheUseCase) claimDepositWithSameMemberAndAmount(cond *model.Deposit, owner string) (bool, error) {
	return packet.Redis.SetNX(depositDedupeKey(cond), owner, depositDedupeTTL).Result()
}

func (dc *depositCacheUseCase) releaseDepositWithSameMemberAndAmount(cond *model.Deposit, owner string) {
	if err := packet.Redis.Eval(
		releaseDepositClaimScript,
		[]string{depositDedupeKey(cond)},
		owner,
	).Err(); err != nil {
		logger.ApLog().Warnf("release deposit dedupe claim failed: memberID=%d err=%v", cond.MemberID, err)
	}
}

func depositDedupeKey(cond *model.Deposit) string {
	return fmt.Sprintf("deposit:dedupe:%d:%s", cond.MemberID, cond.Amount.StringFixed(2))
}
