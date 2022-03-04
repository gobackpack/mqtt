package main

import (
	"context"
	"fmt"
	"github.com/gobackpack/mqtt"
	"github.com/sirupsen/logrus"
	"sync"
)

func main() {
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
}
