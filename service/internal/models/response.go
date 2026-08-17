package models

import "github.com/james730922/wallet/service/internal/utils/errs"

type CommonResponseStatus int

const (
	CommonStatusFailed = iota
	CommonStatusSuccess
)

//ResponseWithData all GET methods should return this struct to client
type ResponseWithData struct {
	Paging *PagingResult `json:"paging,omitempty"`
	Data   interface{}   `json:"data"`
}

//ResponseWithStatus all POST, PUT, PATCH, DELETE methods should return this struct to client
type ResponseWithStatus struct {
	ID     int64                `json:"id,string"`
	Status CommonResponseStatus `json:"status"`
}

type WsResponseWithStatus struct {
	Status CommonResponseStatus `json:"status"`
	Topic  string               `json:"topic"`
	Data   interface{}          `json:"data,omitempty"`
	Error  *WsResponseError     `json:"error,omitempty"`
}

func (resp *WsResponseWithStatus) SetError(err error) {
	e, ok := errs.Parse(err)
	resp.Error = &WsResponseError{}

	if ok {
		resp.Error.Code = string(e.GetCode())
		resp.Error.Message = e.GetMessage()
	} else {
		resp.Error.Message = err.Error()
	}
}

type WsResponseError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
