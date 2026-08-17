package errs

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
)

func Parse(err error) (*Error, bool) {
	if e, ok := err.(*Error); ok {
		return e, true
	} else {
		return nil, false
	}
}

func new(status httpStatus, prefixCode PrefixCode, code int, message string) *Error {
	c, err := strconv.Atoi(fmt.Sprintf("%d%03d", prefixCode, code))
	if err != nil {
		panic(err)
	}

	errorCode := makeCode(status, c)
	codeMap.Add(errorCode, message)
	return &Error{
		statusCode: int(status),
		class:      prefixCode,
		code:       errorCode,
		message:    message,
		params:     make([]interface{}, 0, 0),
	}
}

type Error struct {
	statusCode int
	class      PrefixCode
	code       ErrorCode
	message    string
	params     []interface{}
}

func (c *Error) Error() string {
	if len(c.params) == 0 {
		return string(c.code) + ": " + c.message
	}
	return string(c.code) + ": " + fmt.Sprintf(c.message, c.params)
}

func (d *Error) GetClass() PrefixCode {
	return d.class
}

func (c *Error) GetCode() ErrorCode {
	return c.code
}

func (c *Error) GetStatusCode() int {
	return c.statusCode
}

func (c *Error) GetMessage() string {
	if len(c.params) == 0 {
		return c.message
	}
	return fmt.Sprintf(c.message, c.params...)
}

func (c *Error) Is(err error) bool {
	e, ok := Parse(err)
	if !ok {
		return false
	}
	return c.GetCode() == e.GetCode()
}

func (c *Error) WithParams(params ...interface{}) *Error {
	c.params = params
	return c
}

func ConvertDB(err error) error {
	if err == sql.ErrNoRows || err == gorm.ErrRecordNotFound {
		return DBNoRow
	}

	if e, ok := parseDBError(err); ok {
		return e
	}

	return err
}

var (
	dbErrorCode = map[int]error{
		1005: DBInsertFailed,    // CANNOT_CREATE_TABLE
		1006: DBInsertFailed,    // CANNOT_CREATE_DATABASE
		1007: DBInsertFailed,    // DATABASE_CREATE_EXISTS
		1040: DBOperationFailed, // TOO_MANY_CONNS
		1045: DBAccessDenied,    // ACCESS_DENIED
		1051: DBOperationFailed, // UNKNOWN_TABLE
		1062: DBInsertDuplicate, //DUPLICATE_ENTRY
	}
)

func parseDBError(err error) (error, bool) {
	s := strings.TrimSpace(err.Error())
	data := strings.Split(s, ":")
	if len(data) == 0 {
		return nil, false
	}

	numStr := strings.ToLower(data[0])
	numStr = strings.Replace(numStr, "error", "", -1)
	numStr = strings.TrimSpace(numStr)
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return nil, false
	}

	e, ok := dbErrorCode[num]
	if !ok {
		return nil, false
	}

	return e, true
}
