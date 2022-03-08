### ENV variables

| ENV            | Default value         |
|:---------------|:---------------------:|
| MQTT_HOST      | localhost             |
| MQTT_PORT      | 1883                  |
| MQTT_USERNAME  | guest                 |
| MQTT_PASSWORD  | guest                 |
| MQTT_CLIENT_ID | uuid.New().String()   |

### Usage

```go
mqttConfig := mqtt.NewConfig()
hub := mqtt.NewHub(mqttConfig)

hubCtx, hubCancel := context.WithCancel(context.Background())
defer hubCancel()

hubFinished, err := hub.Connect(hubCtx)
if err != nil {
    logrus.Fatal(err)
}

// sub
subFinished, onMessage, onError := hub.Subscribe(hubCtx, "mytopic")

go func(ctx context.Context) {
    for {
        select {
        case msg := <-onMessage:
            logrus.Infof("[mytopic]: %s", string(msg))
            break
        case err := <-onError:
            logrus.Error(err)
            break
        case <-ctx.Done():
            return
        }
    }
}(hubCtx)

// pub
go hub.Publish("mytopic", []byte("message"))
```

* **Customize config values (*these are defaults*)**
```go
mqttConfig.KeepAlive = time.Second * 15
mqttConfig.CleanSession = true
mqttConfig.AutoReconnect = true
mqttConfig.MsgChanDept = 100
mqttConfig.PubQoS = 0
mqttConfig.SubQoS = 0
```