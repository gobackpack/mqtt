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
	hubCancelled, err := hub.Connect(hubCtx)
	if err != nil {
		logrus.Fatal(err)
	}

	// sub
	hub.OnMessage = make(chan []byte)
	hub.OnError = make(chan error)
	subCtx, subCancel := context.WithCancel(hubCtx)
	subCancelled := hub.Subscribe(subCtx, "mytopic")

	go func(subCancel context.CancelFunc) {
		defer subCancel()

		mCount := 0
		for {
			select {
			case msg := <-hub.OnMessage:
				logrus.Info("message received: ", string(msg))
				
				mCount++
				if mCount == 100 { // we decide when to stop subscription
					return
				}
				break
			case err := <-hub.OnError:
				logrus.Error(err)
				break
			}
		}
	}(subCancel)

	// pub
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

	<-subCancelled
	close(subCancelled)

	hubCancel()
	<-hubCancelled
	close(hubCancelled)
}
