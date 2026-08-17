package tools

import (
	"sort"
	"strings"
	"time"

	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

const (
	DateFormatLayout       = "2006-01-02"
	DateHourFormatLayout   = "2006-01-02 15"
	DateMinuteFormatLayout = "2006-01-02 15:04"
	DateSecondFormatLayout = "2006-01-02 15:04:05"
	offsetUTCMinus4        = -4 * 60 * 60
	EasternTimeZone        = "(美东)"
)

// 解析RFC3339時間
func ParseRFC3339(timeRFC3339 string) (*time.Time, error) {
	t, err := time.Parse(time.RFC3339, timeRFC3339)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.CommonRequestParamInvalid
	}

	t = t.UTC()

	return &t, nil
}

// 解析RFC3339時間
// 因用於報表輸出，時間為空是合理狀態，故不打印錯誤訊息直接回空
func ParseRFC3339IgnoreNull(timeRFC3339 string) *time.Time {
	if timeRFC3339 == "" {
		return nil
	}

	t, err := time.Parse(time.RFC3339, timeRFC3339)
	if err != nil {
		return nil
	}

	t = t.UTC()

	return &t
}

// 轉換為美東時區日期
func ETZoneDate(tp time.Time) (year int, month time.Month, day int) {
	utcMinus4Zone := time.FixedZone("", offsetUTCMinus4)
	return tp.In(utcMinus4Zone).Date()
}

func ETZoneTimeFormat(tp time.Time) string {
	utcMinus4Zone := time.FixedZone("", offsetUTCMinus4)
	return tp.In(utcMinus4Zone).Format(DateSecondFormatLayout)
}

// ETZoneFormatWithDefaultValue
// returns "" when time got nil value
func ETZoneFormatWithDefaultValue(tp *time.Time) string {
	if tp != nil {
		utcMinus4Zone := time.FixedZone("", offsetUTCMinus4)
		return tp.In(utcMinus4Zone).Format(DateSecondFormatLayout)
	}

	return ""
}

func ETZoneTime(tp time.Time) time.Time {
	utcMinus4Zone := time.FixedZone("", offsetUTCMinus4)
	return tp.In(utcMinus4Zone)
}

func NowInETZone() time.Time {
	return ETZoneTime(time.Now())
}

func TheLatestTime(times ...time.Time) time.Time {
	sort.Slice(times, func(i, j int) bool {
		if times[i].After(times[j]) {
			return true
		}
		return false
	})

	return times[0]
}

func RFC3339ToDateLayout(date string) string {
	result := strings.Split(date, "T")

	if len(result) != 2 {
		logger.ApLog().Errorf("err:%v, date:%s\n", errs.CommonRequestParamInvalid, date)
		return ""
	}

	return result[0]
}

func ETZoneLocation() *time.Location {
	return time.FixedZone("", offsetUTCMinus4)
}
