package mqtt

import (
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Hub struct {
	OnMessage chan []byte
	OnError   chan error

	conn    *connection
	publish chan *frame
}

type frame struct {
	topic   string
	payload []byte
}

func NewHub(conf *Config) *Hub {
	return &Hub{
		conn:    newConnection(conf),
		publish: make(chan *frame),
	}
}

func (hub *Hub) Connect(ctx context.Context) (chan bool, error) {
	if err := hub.conn.connect(); err != nil {
		return nil, err
	}

	cancelled := make(chan bool)

	go func(ctx context.Context) {
		defer func() {
			cancelled <- true
		}()

		for {
			select {
			case fr := <-hub.publish:
				hub.conn.publish(fr.topic, fr.payload)
				break
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	return cancelled, nil
}

func (hub *Hub) Publish(topic string, message []byte) {
	hub.publish <- &frame{
		topic:   topic,
		payload: message,
	}
}

func (hub *Hub) Subscribe(ctx context.Context, topic string) chan bool {
	cancelled := make(chan bool)

	go func() {
		defer func() {
			cancelled <- true
		}()

		go hub.listenForMessages(topic)

		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}()

	return cancelled
}

func (hub *Hub) listenForMessages(topic string) {
	if token := hub.conn.subscribe(topic, func(mqttClient mqtt.Client, message mqtt.Message) {
		hub.OnMessage <- message.Payload()
	}); token.Wait() && token.Error() != nil {
		hub.OnError <- token.Error()
	}
}
