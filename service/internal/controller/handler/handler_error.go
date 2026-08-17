package handler

import (
	"github.com/golang/protobuf/proto"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"net/http"
)

func newErrorHandler() ErrorHandlerUseCase {
	return ErrorHandlerUseCase{}
}

type ErrorHandlerUseCase struct {
}

type Error struct {
	StatusCode int
	InnerError *zqbapis.Error
}

func (eh ErrorHandlerUseCase) Marshal(err error) []byte {
	result, marshalErr := proto.Marshal(eh.Parse(err).InnerError)
	if marshalErr != nil {
		return []byte{}
	}

	return result
}

func (eh ErrorHandlerUseCase) Parse(err error) *Error {
	e, ok := errs.Parse(err)
	if !ok {
		tmp := errs.CommonUnknownError
		return &Error{
			StatusCode: http.StatusInternalServerError,
			InnerError: &zqbapis.Error{
				Code:    string(tmp.GetCode()),
				Message: tmp.GetMessage(),
			},
		}
	}

	return &Error{
		StatusCode: e.GetStatusCode(),
		InnerError: &zqbapis.Error{
			Code:    string(e.GetCode()),
			Message: e.GetMessage(),
		},
	}
}

func (eh ErrorHandlerUseCase) IsDBError(err error) bool {
	selfDefilneError, ok := errs.Parse(err)
	if !ok {
		selfDefilneError = errs.CommonUnknownError
	}

	return selfDefilneError.GetClass() == errs.CodeDB
}
