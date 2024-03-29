package mqtt

import (
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

const (
	DefaultPubQoS = 0
	DefaultSubQoS = 0
)

type Hub struct {
	conn *connection
}

type Subscriber struct {
	OnMessage chan []byte
	OnError   chan error
}

type Publisher struct {
	OnError chan error
	publish chan *frame
}

type frame struct {
	topic   string
	qos     int
	payload []byte
}

func NewHub(conf *Config) *Hub {
	return &Hub{
		conn: newConnection(conf),
	}
}

// Connect to MQTT server
func (hub *Hub) Connect(ctx context.Context) error {
	if err := hub.conn.connect(); err != nil {
		return err
	}

	go func(hub *Hub) {
		defer logrus.Warn("hub closed MQTT connection")

		for {
			select {
			case <-ctx.Done():
				hub.conn.client.Disconnect(1000)
				return
			}
		}
	}(hub)

	return nil
}

// Subscribe will create MQTT subscriber and listen for messages.
// Messages and errors are sent to OnMessage and OnError channels.
func (hub *Hub) Subscribe(ctx context.Context, topic string, qos int) *Subscriber {
	sub := &Subscriber{
		OnMessage: make(chan []byte),
		OnError:   make(chan error),
	}

	go func(hub *Hub, sub *Subscriber, topic string, qos int) {
		defer func() {
			close(sub.OnMessage)
			close(sub.OnError)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				if token := hub.conn.subscribe(topic, byte(qos), func(mqttClient mqtt.Client, message mqtt.Message) {
					if payload := message.Payload(); len(payload) > 0 {
						sub.OnMessage <- payload
					}
				}); token.Wait() && token.Error() != nil {
					sub.OnError <- token.Error()
				}
			}
		}
	}(hub, sub, topic, qos)

	return sub
}

// Publisher will create MQTT publisher and private listener for messages to be published.
// All messages to be published are sent through private publish channel.
// Errors will be sent to OnError channel.
func (hub *Hub) Publisher(ctx context.Context) *Publisher {
	pub := &Publisher{
		OnError: make(chan error),
		publish: make(chan *frame),
	}

	go func(hub *Hub, pub *Publisher) {
		defer close(pub.OnError)

		for {
			select {
			case fr := <-pub.publish:
				if token := hub.conn.publish(fr.topic, byte(fr.qos), fr.payload); token.Wait() && token.Error() != nil {
					pub.OnError <- token.Error()
				}
			case <-ctx.Done():
				return
			}
		}
	}(hub, pub)

	return pub
}

// Publish message to topic through private pub.publish channel.
// Thread-safe.
func (pub *Publisher) Publish(topic string, qos int, message []byte) {
	pub.publish <- &frame{
		topic:   topic,
		qos:     qos,
		payload: message,
	}
}
