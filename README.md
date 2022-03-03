## ENV variables

| ENV            | Default value         |
|:---------------|:---------------------:|
| MQTT_HOST      | localhost             |
| MQTT_PORT      | 1883                  |
| MQTT_USERNAME  | guest                 |
| MQTT_PASSWORD  | guest                 |
| MQTT_CLIENT_ID | uuid.New().String()   |

## Usage

* **Create mqtt config**
```go
mqttConfig := mqtt.NewConfig()
```

* **Optionally, customize config values (*these are defaults*)**
```go
mqttConfig.KeepAlive = time.Second * 15
mqttConfig.CleanSession = true
mqttConfig.AutoReconnect = true
mqttConfig.MsgChanDept = 100
mqttConfig.PubQoS = 0
mqttConfig.SubQoS = 0
```

* **Create mqtt connection**
```go
mqttConn, err := mqtt.NewConnection(mqttConfig)
if err != nil {
    return err
}
```

* **Publish payload to mqtt topic**
```go
if token := mqttConn.Publish("my/topic", []byte("message")); token.Wait() && token.Error() != nil {
    log.Print("mqtt publish error: ", token.Error())
}
```

* **Subscribe to mqtt topic**
```go
if token := mqttConn.Subscribe("my/topic", func(mqttClient mqtt.Client, message mqtt.Message) {
    log.Print(message)
}); token.Wait() && token.Error() != nil {
    log.Print("mqtt subscribe error: ", token.Error())
}
```