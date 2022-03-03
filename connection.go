package mqtt

import (
	"errors"

	mqttLib "github.com/eclipse/paho.mqtt.golang"
)

type Connection struct {
	Config *Config
	Client mqttLib.Client
}

func NewConnection(config *Config) (*Connection, error) {
	conn := &Connection{
		Config: config,
	}

	opts := mqttLib.NewClientOptions()

	broker := conn.Config.Host + ":" + conn.Config.Port

	opts.AddBroker(broker)
	opts.SetClientID(conn.Config.ClientID)
	opts.SetUsername(conn.Config.Username)
	opts.SetPassword(conn.Config.Password)
	opts.SetCleanSession(conn.Config.CleanSession)
	opts.SetAutoReconnect(conn.Config.AutoReconnect)
	opts.SetKeepAlive(conn.Config.KeepAlive)
	opts.SetMessageChannelDepth(conn.Config.MsgChanDept)

	conn.Client = mqttLib.NewClient(opts)
	if token := conn.Client.Connect(); token.Wait() && token.Error() != nil {
		return nil, errors.New("MQTT client connection failed: " + token.Error().Error())
	}

	return conn, nil
}

// Publish payload to topic
func (conn *Connection) Publish(topic string, payload []byte) mqttLib.Token {
	return conn.Client.Publish(topic, byte(conn.Config.PubQoS), conn.Config.Retained, payload)
}

// Subscribe to topic
func (conn *Connection) Subscribe(topic string, callback func(mqttClient mqttLib.Client, message mqttLib.Message)) mqttLib.Token {
	return conn.Client.Subscribe(topic, byte(conn.Config.SubQoS), callback)
}
