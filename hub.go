package mqtt

import (
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

type Hub struct {
	conn *connection
}

type Subscriber struct {
	OnMessage chan []byte
	OnError   chan error
	Finished  chan bool
}

type Publisher struct {
	OnError chan error
	publish chan *frame
}

type frame struct {
	topic   string
	payload []byte
}

func NewHub(conf *Config) *Hub {
	return &Hub{
		conn: newConnection(conf),
	}
}

func (hub *Hub) Connect(ctx context.Context) error {
	if err := hub.conn.connect(); err != nil {
		return err
	}

	go func(ctx context.Context) {
		defer logrus.Warn("hub closed MQTT connection")

		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	return nil
}

func (hub *Hub) Publish(topic string, message []byte, pub *Publisher) {
	pub.publish <- &frame{
		topic:   topic,
		payload: message,
	}
}

func (hub *Hub) Subscribe(ctx context.Context, topic string) *Subscriber {
	sub := &Subscriber{
		OnMessage: make(chan []byte),
		OnError:   make(chan error),
		Finished:  make(chan bool),
	}

	go func(ctx context.Context, sub *Subscriber) {
		defer func() {
			close(sub.Finished)
		}()

		if token := hub.conn.subscribe(topic, func(mqttClient mqtt.Client, message mqtt.Message) {
			sub.OnMessage <- message.Payload()
		}); token.Wait() && token.Error() != nil {
			sub.OnError <- token.Error()
		}

		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}(ctx, sub)

	return sub
}

func (hub *Hub) Publisher(ctx context.Context) *Publisher {
	pub := &Publisher{
		OnError: make(chan error),
		publish: make(chan *frame),
	}

	go func(ctx context.Context, pub *Publisher) {
		for {
			select {
			case fr := <-pub.publish:
				if token := hub.conn.publish(fr.topic, fr.payload); token.Wait() && token.Error() != nil {
					pub.OnError <- token.Error()
				}
			case <-ctx.Done():
				return
			}
		}
	}(ctx, pub)

	return pub
}
