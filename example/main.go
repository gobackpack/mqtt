package main

import (
	"fmt"
	"github.com/gobackpack/mqtt"
	"github.com/sirupsen/logrus"
	"log"
)

func main() {
	mqttConfig := mqtt.NewConfig()

	mqttConn, err := mqtt.NewConnection(mqttConfig)
	if err != nil {
		logrus.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		go func(i int) {
			msg := fmt.Sprintf("message %d", i)
			if token := mqttConn.Publish("my/topic", []byte(msg)); token.Wait() && token.Error() != nil {
				log.Print("mqtt publish error: ", token.Error())
			}
		}(i)
	}
}
