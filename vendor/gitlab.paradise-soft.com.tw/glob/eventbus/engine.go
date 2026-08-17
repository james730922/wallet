package eventbus

type IEngine interface {
	Publish(topic string, message []byte) error
	Subscribe(topic string) (<-chan []byte, error)
	Unsubscribe(topic string)
}
