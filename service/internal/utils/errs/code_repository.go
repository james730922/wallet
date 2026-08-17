package errs

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

var (
	codeMap = newCodeRepository()
)

// newCodeRepository: registration error code, and check error code is unique.
func newCodeRepository() iCodeRepository {
	return &codeRepository{
		m: make(map[ErrorCode]string),
		c: make(map[int]struct{}),
	}
}

type iCodeRepository interface {
	Add(code ErrorCode, message string)
}

type codeRepository struct {
	m map[ErrorCode]string
	c map[int]struct{}
}

func (cr *codeRepository) Add(errorCode ErrorCode, message string) {
	_, code := parseCode(errorCode)

	if _, ok := cr.c[code]; ok {
		log.Panicf("error code duplicate definition, code: %d", code)
	}

	cr.c[code] = struct{}{}
	cr.m[errorCode] = message
}

func makeCode(status httpStatus, code int) ErrorCode {
	return ErrorCode(fmt.Sprintf("%d-%d", status, code))
}

func parseCode(errorCode ErrorCode) (httpStatus, code int) {
	tmp := strings.Split(string(errorCode), "-")
	if len(tmp) != 2 {
		log.Panicf("parseCode format error, errorCode: %s", errorCode)
	}

	s, err := strconv.Atoi(tmp[0])
	if err != nil {
		log.Panicf("parseCode format error, errorCode: %s", errorCode)
	}

	c, err := strconv.Atoi(tmp[1])
	if err != nil {
		log.Panicf("parseCode format error, errorCode: %s", errorCode)
	}

	return s, c
}
