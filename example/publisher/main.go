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
	defer hubCancel()

	if err := hub.Connect(hubCtx); err != nil {
		logrus.Fatal(err)
	}

	pub1 := hub.Publisher(hubCtx)
	pub2 := hub.Publisher(hubCtx)

	go func(ctx context.Context) {
		for {
			select {
			case err := <-pub1.OnError:
				logrus.Error(err)
			case err := <-pub2.OnError:
				logrus.Error(err)
			case <-ctx.Done():
				return
			}
		}
	}(hubCtx)

	// pub
	wg := sync.WaitGroup{}
	wg.Add(200)

	for i := 0; i < 100; i++ {
		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()

			msg := []byte(fmt.Sprintf("message %d", i))
			hub.Publish("mytopic", msg, pub1)
		}(i, &wg)

		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()

			msg := []byte(fmt.Sprintf("message %d", i))
			hub.Publish("mytopic2", msg, pub2)
		}(i, &wg)
	}

	wg.Wait()
}
