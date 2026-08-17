package eventbus

import (
	"sync"
)

func newSubscribeRepository() ISubscribeRepository {
	return &subscribeRepository{
		topic:   make(map[string]map[EntryID]CallbackFunc),
		entryID: make(map[EntryID]string),
	}
}

type ISubscribeRepository interface {
	Store(topic string, entryID EntryID, callback CallbackFunc)
	Load(topic string) ([]CallbackFunc, bool)
	LoadTopic(entryID EntryID) (string, bool)
	Delete(entryID EntryID)
}

type subscribeRepository struct {
	topic   map[string]map[EntryID]CallbackFunc // key: topic
	entryID map[EntryID]string                  // key: EntryID, val: topic
	mx      sync.RWMutex
}

func (sr *subscribeRepository) Store(topic string, entryID EntryID, callback CallbackFunc) {
	sr.mx.Lock()

	callbacks, ok := sr.topic[topic]
	if !ok {
		callbacks = make(map[EntryID]CallbackFunc)
		sr.topic[topic] = callbacks
	}

	callbacks[entryID] = callback
	sr.entryID[entryID] = topic

	sr.mx.Unlock()
}

func (sr *subscribeRepository) Load(topic string) ([]CallbackFunc, bool) {
	sr.mx.RLock()
	defer sr.mx.RUnlock()

	callbacks, ok := sr.topic[topic]
	if !ok {
		return nil, false
	}

	result := make([]CallbackFunc, 0, len(callbacks))
	for _, callback := range callbacks {
		result = append(result, callback)
	}

	return result, true
}

func (sr *subscribeRepository) LoadTopic(entryID EntryID) (string, bool) {
	sr.mx.RLock()
	defer sr.mx.RUnlock()

	topic, ok := sr.entryID[entryID]

	return topic, ok
}

func (sr *subscribeRepository) Delete(entryID EntryID) {
	sr.mx.Lock()
	defer sr.mx.Unlock()

	topic, ok := sr.entryID[entryID]
	if !ok {
		return
	}

	callbacks, ok := sr.topic[topic]
	if !ok {
		return
	}

	delete(callbacks, entryID)
	delete(sr.entryID, entryID)
}
