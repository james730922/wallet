package eventbus

import "log"

func newReceiver(topic string, subscribeRepository ISubscribeRepository, messageCh <-chan []byte) IReceiver {
	return (&receiver{
		topic:               topic,
		subscribeRepository: subscribeRepository,
		messageCh:           messageCh,
	}).run()
}

type IReceiver interface {
}

type receiver struct {
	topic               string
	subscribeRepository ISubscribeRepository
	messageCh           <-chan []byte
}

func (rc *receiver) run() *receiver {
	go rc.receive()
	return rc
}

func (rc *receiver) receive() {
	for m := range rc.messageCh {
		go rc.doCallback(m)
	}
}

func (rc *receiver) doCallback(buf []byte) {
	callbacks, ok := rc.subscribeRepository.Load(rc.topic)
	if !ok {
		log.Printf("callback [%s] not exists", rc.topic)
		return
	}

	for _, f := range callbacks {
		go f(buf)
	}
}
