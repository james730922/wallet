package eventbus

import "errors"

type CallbackFunc func(data []byte)

type EntryID int

var (
	ErrIllegalTopic        = errors.New("illegal topic, con't empty")
	ErrDuplicateSubscripts = errors.New("duplicate subscripts")
)

var (
	errRedigoConnectionClosed = errors.New("redigo: connection closed")
)
