package event

import (
	"strconv"

	"gitlab.paradise-soft.com.tw/glob/eventbus"
)

type CallbackFunc = eventbus.CallbackFunc
type EntryID = eventbus.EntryID

type Topic int

const (
	TopicWalletMemberChange  Topic = 4
	TopicPaymentUpdateNotify Topic = 6
)

var e *event

func New() {
	e = &event{
		IEvent: eventbus.New(eventbus.WithInternalEngine()),
	}
}

func Event() *event {
	return e
}

type event struct {
	eventbus.IEvent
}

func (e *event) Subscribe(topic Topic, callback CallbackFunc) (EntryID, error) {
	return e.IEvent.Subscribe(strconv.Itoa(int(topic)), callback)
}

func (e *event) Publish(topic Topic, message interface{}) error {
	return e.IEvent.Publish(strconv.Itoa(int(topic)), message)
}

func Unmarshal(data []byte, key interface{}) error {
	return eventbus.Unmarshal(data, key)
}
