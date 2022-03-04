package mqtt

import "context"

type Hub struct {
	conn    *Connection
	publish chan *Frame
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
