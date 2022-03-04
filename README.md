## ENV variables

| ENV            | Default value         |
|:---------------|:---------------------:|
| MQTT_HOST      | localhost             |
| MQTT_PORT      | 1883                  |
| MQTT_USERNAME  | guest                 |
| MQTT_PASSWORD  | guest                 |
| MQTT_CLIENT_ID | uuid.New().String()   |

## Usage

```go
mqttConfig := mqtt.NewConfig()

hub := mqtt.NewHub(mqttConfig)

hubCtx, hubCancel := context.WithCancel(context.Background())
cancelled, err := hub.Connect(hubCtx)
if err != nil {
    logrus.Fatal(err)
}

wg := sync.WaitGroup{}
wg.Add(100)

for i := 0; i < 100; i++ {
    go func(i int, wg *sync.WaitGroup) {
        defer wg.Done()

        msg := []byte(fmt.Sprintf("message %d", i))
        hub.Publish("mytopic", msg)
    }(i, &wg)
}

wg.Wait()

hubCancel()
<-cancelled
close(cancelled)
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