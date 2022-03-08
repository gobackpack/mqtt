package mqtt

import (
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Hub struct {
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

	finished := make(chan bool)

	go func(ctx context.Context) {
		defer func() {
			finished <- true
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

	return finished, nil
}

func (hub *Hub) Publish(topic string, message []byte) {
	hub.publish <- &frame{
		topic:   topic,
		payload: message,
	}
}

func (hub *Hub) Subscribe(ctx context.Context, topic string) (chan bool, chan []byte, chan error) {
	finished := make(chan bool)
	onMessage := make(chan []byte)
	onError := make(chan error)

	go func(ctx context.Context, onMessage chan []byte, onError chan error) {
		defer func() {
			finished <- true
		}()

		if token := hub.conn.subscribe(topic, func(mqttClient mqtt.Client, message mqtt.Message) {
			onMessage <- message.Payload()
		}); token.Wait() && token.Error() != nil {
			onError <- token.Error()
		}

		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}(ctx, onMessage, onError)

	return finished, onMessage, onError
}
