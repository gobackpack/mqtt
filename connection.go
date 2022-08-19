package mqtt

import (
	mqttLib "github.com/eclipse/paho.mqtt.golang"
)

type connection struct {
	conf   *Config
	client mqttLib.Client
}

func newConnection(conf *Config) *connection {
	conn := &connection{
		conf: conf,
	}

	opts := mqttLib.NewClientOptions()

	broker := conn.conf.Host + ":" + conn.conf.Port

	opts.AddBroker(broker)
	opts.SetClientID(conn.conf.ClientID)
	opts.SetUsername(conn.conf.Username)
	opts.SetPassword(conn.conf.Password)
	opts.SetCleanSession(conn.conf.CleanSession)
	opts.SetAutoReconnect(conn.conf.AutoReconnect)
	opts.SetKeepAlive(conn.conf.KeepAlive)
	opts.SetMessageChannelDepth(conn.conf.MsgChanDept)

	conn.client = mqttLib.NewClient(opts)

	return conn
}

func (conn *connection) connect() error {
	if token := conn.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (conn *connection) publish(topic string, qos byte, payload []byte) mqttLib.Token {
	return conn.client.Publish(topic, qos, conn.conf.Retained, payload)
}

func (conn *connection) subscribe(topic string, qos byte, callback func(mqttClient mqttLib.Client, message mqttLib.Message)) mqttLib.Token {
	return conn.client.Subscribe(topic, qos, callback)
}
