package eventbus

import (
	"bytes"
	"encoding/gob"
	"strings"
	"sync"
)

func New(engine IEngine) IEvent {
	return &Event{
		engine:              engine,
		seq:                 newSequence(),
		subscribeRepository: newSubscribeRepository(),
		receiveRepository:   newReceiveRepository(),
	}
}

type IEvent interface {
	Subscribe(topic string, callback CallbackFunc) (EntryID, error)
	Unsubscribe(entryID EntryID)
	Publish(topic string, message interface{}) error
}

type Event struct {
	engine              IEngine
	seq                 ISequence
	subscribeRepository ISubscribeRepository
	receiveRepository   IReceiveRepository
	mx                  sync.Mutex
}

func (ev *Event) Subscribe(topic string, callback CallbackFunc) (EntryID, error) {
	ev.mx.Lock()
	defer ev.mx.Unlock()

	entryID := ev.seq.Next()
	ev.subscribeRepository.Store(topic, entryID, callback)

	if err := ev.receive(topic); err != nil {
		return 0, err
	}

	return entryID, nil
}

func (ev *Event) Unsubscribe(entryID EntryID) {
	ev.mx.Lock()
	defer ev.mx.Unlock()

	topic, ok := ev.subscribeRepository.LoadTopic(entryID)
	if !ok {
		return
	}

	ev.subscribeRepository.Delete(entryID)

	callbacks, ok := ev.subscribeRepository.Load(topic)
	if ok && len(callbacks) == 0 {
		ev.receiveRepository.Delete(topic)
		ev.engine.Unsubscribe(topic)
	}
}

func (ev *Event) Publish(topic string, message interface{}) error {
	if strings.TrimSpace(topic) == "" {
		return ErrIllegalTopic
	}

	buf ,err := Marshal(message)
	if err!=nil {
		return err
	}

	ev.engine.Publish(topic, buf)

	return nil
}

func (ev *Event) receive(topic string) error {
	if _, ok := ev.receiveRepository.Load(topic); !ok {
		messageCh, err := ev.engine.Subscribe(topic)
		if err != nil {
			return err
		}
		ev.receiveRepository.Store(topic, newReceiver(topic, ev.subscribeRepository, messageCh))
	}

	return nil
}

func Marshal(data interface{}) ([]byte, error) {
	var buf bytes.Buffer

	err := gob.NewEncoder(&buf).Encode(data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}


func Unmarshal(data []byte, key interface{}) error {
	var buf bytes.Buffer

	buf.Write(data)

	err := gob.NewDecoder(&buf).Decode(key)
	if err != nil {
		return err
	}

	return nil
}
