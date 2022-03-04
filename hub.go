package mqtt

import (
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Hub struct {
	conn    *Connection
	publish chan *Frame

	OnMessage chan []byte
	OnError   chan error
}

type Frame struct {
	Topic   string
	Payload []byte
}

func NewHub(conf *Config) *Hub {
	return &Hub{
		conn:    NewConnection(conf),
		publish: make(chan *Frame),
	}
}

func (hub *Hub) Connect(ctx context.Context) (chan bool, error) {
	if err := hub.conn.Connect(); err != nil {
		return nil, err
	}

	cancelled := make(chan bool)

	go func(ctx context.Context) {
		defer func() {
			cancelled <- true
		}()

		for {
			select {
			case frame := <-hub.publish:
				hub.conn.Publish(frame.Topic, frame.Payload)
				break
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	return cancelled, nil
}

func (hub *Hub) Publish(topic string, message []byte) {
	hub.publish <- &Frame{
		Topic:   topic,
		Payload: message,
	}
}

func (hub *Hub) Subscribe(ctx context.Context, topic string) chan bool {
	cancelled := make(chan bool)

	go func() {
		defer func() {
			cancelled <- true
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				if token := hub.conn.Subscribe(topic, func(mqttClient mqtt.Client, message mqtt.Message) {
					hub.OnMessage <- message.Payload()
				}); token.Wait() && token.Error() != nil {
					hub.OnError <- token.Error()
				}
			}
		}
	}()

	return cancelled
}
