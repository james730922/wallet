package models

type Topic int

const (
	TopicUnknown Topic = iota
	TopicRegister
	TopicValidate
)

func TopicName(topic Topic) string {
	m := map[Topic]string{
		TopicUnknown:  "unknown",
		TopicRegister: "register",
		TopicValidate: "validate",
	}
	return m[topic]
}

func GetTopic(topic string) Topic {
	m := map[string]Topic{
		"unknown":  TopicUnknown,
		"register": TopicRegister,
		"validate": TopicValidate,
	}

	if t, ok := m[topic]; ok {
		return t
	}
	return -1
}

type WsRequest struct {
	Topic string      `json:"topic"`
	Data  interface{} `json:"data"`
}
