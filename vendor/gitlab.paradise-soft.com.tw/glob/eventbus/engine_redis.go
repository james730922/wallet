package eventbus

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/go-redis/redis"
)

func WithRedisEngine(addr, password string) IEngine {
	redisCli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})

	engine := &RedisEngine{
		redisCli: redisCli,
		topic:    make(map[string]chan<- []byte),
		stop:     make(map[string]chan struct{}),
	}

	engine.run()

	return engine
}

type RedisEngine struct {
	redisCli *redis.Client
	topic    map[string]chan<- []byte
	stop     map[string]chan struct{} // key: topic, val: stop signal
	mx       sync.Mutex
}

func (re *RedisEngine) run() *RedisEngine {
	go re.interruptNotify()
	return re
}

func (re *RedisEngine) Publish(topic string, message []byte) error {
	if err := re.redisCli.Publish(topic, message).Err(); err != nil {
		return err
	}

	return nil
}

// 訂閱主題
// 建立連線後會返回chan，失敗時返回error
// 如果chan 被關閉時，表示接收出了問題
func (re *RedisEngine) Subscribe(topic string) (<-chan []byte, error) {
	re.mx.Lock()
	defer re.mx.Unlock()

	if _, ok := re.stop[topic]; ok {
		return nil, ErrDuplicateSubscripts
	}

	subscribe := re.redisCli.Subscribe(topic)

	stop := make(chan struct{}, 1)
	re.stop[topic] = stop

	msg := make(chan []byte, 1)

	go func() {
		defer subscribe.Close()
		defer close(msg)

		for {
			select {
			case <-stop:
				return
			default:
				v, err := subscribe.ReceiveMessage()
				if err != nil {
					log.Printf("RedisEngine ReceiveMessage topic[%s]", topic)
				}
				msg <- []byte(v.Payload)
			}
		}
	}()

	return msg, nil
}

func (re *RedisEngine) Unsubscribe(topic string) {
	re.mx.Lock()
	defer re.mx.Unlock()

	stop, ok := re.stop[topic]
	if !ok {
		return
	}

	select {
	case stop <- struct{}{}:
		delete(re.stop, topic)
	default:
	}
}

func (re *RedisEngine) interruptNotify() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case <-interrupt:
		re.mx.Lock()
		for k, v := range re.stop {
			v <- struct{}{}
			delete(re.stop, k)
		}
		re.mx.Unlock()
	}
}
