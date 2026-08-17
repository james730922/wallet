package tools

import (
	"regexp"

	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

const (
	scanPayPasswdLen      = 6
	scanPayPasswdRegexNum = `[0-9]{6}`
)

func CheckSecurityPasswdRule(newPasswd string) error {
	/* PM提供的密碼規則
	長度6位數
	只能數字
	*/
	if len(newPasswd) != scanPayPasswdLen {
		return errs.RegistrationScanPayPasswdFmtErr
	}

	//至少一個數字
	match, err := regexp.MatchString(scanPayPasswdRegexNum, newPasswd)
	if err != nil {
		logger.ApLog().Error(err)
		return errs.CommonParseError
	}
	if !match {
		return errs.RegistrationScanPayPasswdFmtErr
	}

	return nil
}
