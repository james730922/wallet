package eventbus

import (
	"encoding/binary"
	"errors"
	"sync"
)

func WithInternalEngine() IEngine {
	engine := &InternalEngine{
		messageQueue: make(chan []byte, 128),
		topic:        make(map[string]chan<- []byte),
		stop:         make(map[string]chan struct{}),
	}

	engine.run()

	return engine
}

type InternalEngine struct {
	messageQueue chan []byte
	topic        map[string]chan<- []byte // key: topic, val: message chan
	stop         map[string]chan struct{} // key: topic, val: stop signal
	mx           sync.RWMutex
}

func (ie *InternalEngine) run() *InternalEngine {
	go ie.pipe()
	return ie
}

func (ie *InternalEngine) Publish(topic string, message []byte) error {
	ie.messageQueue <- ie.pack(topic, message)
	return nil
}

func (ie *InternalEngine) Subscribe(topic string) (<-chan []byte, error) {
	ie.mx.Lock()
	defer ie.mx.Unlock()

	if _, ok := ie.topic[topic]; ok {
		return nil, ErrDuplicateSubscripts
	}

	mq := make(chan []byte, 128)
	ie.topic[topic] = mq

	stop := make(chan struct{}, 1)
	ie.stop[topic] = stop

	msg := make(chan []byte, 1)

	go func() {
		defer close(msg)
		defer func() {
			ie.mx.Lock()
			delete(ie.topic, topic)
			ie.mx.Unlock()
		}()
		for {
			select {
			case m := <-mq:
				msg <- m
			case <-stop:
				return
			}
		}
	}()

	return msg, nil
}

func (ie *InternalEngine) Unsubscribe(topic string) {
	ie.mx.Lock()
	defer ie.mx.Unlock()

	stop, ok := ie.stop[topic]
	if !ok {
		return
	}

	select {
	case stop <- struct{}{}:
		delete(ie.stop, topic)
	default:
	}
}

func (ie *InternalEngine) pipe() {
	for mq := range ie.messageQueue {
		topic, msg, err := ie.unpack(mq)
		if err != nil {
			continue
		}

		ie.mx.RLock()
		ch, ok := ie.topic[topic]
		ie.mx.RUnlock()

		if !ok {
			continue
		}
		ch <- msg
	}
}

func (ie *InternalEngine) pack(topic string, message []byte) []byte {
	topicBuf := []byte(topic)
	//msgBuf := []byte(message)
	msgBuf := message

	lenBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(lenBuf, uint16(len(topicBuf)))

	lenBufSize := len(lenBuf)
	topicBufSize := len(topicBuf)
	msgBufSize := len(msgBuf)
	size := lenBufSize + topicBufSize + msgBufSize
	buf := make([]byte, size)
	copy(buf[:lenBufSize], lenBuf)
	copy(buf[lenBufSize:lenBufSize+topicBufSize], topicBuf)
	copy(buf[lenBufSize+topicBufSize:], msgBuf)

	return buf
}

func (ie *InternalEngine) unpack(buf []byte) (topic string, message []byte, err error) {
	if len(buf) < 2 {
		return "", []byte{}, errors.New("illegal packet size ")
	}

	topicSize := binary.LittleEndian.Uint16(buf[:2])
	topic = string(buf[2 : 2+topicSize])
	message = buf[2+topicSize:]

	return topic, message, nil
}
