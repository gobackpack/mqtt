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

	hubFinished, err := hub.Connect(hubCtx)
	if err != nil {
		logrus.Fatal(err)
	}

	// sub
	subFinished, onMessage1, onError1 := hub.Subscribe(hubCtx, "mytopic")
	subFinished2, onMessage2, onError2 := hub.Subscribe(hubCtx, "mytopic2")

	go func(ctx context.Context) {
		defer func() {
			subFinished <- true
			subFinished2 <- true
		}()

		c1 := 0
		c2 := 0
		for {
			select {
			case msg := <-onMessage1:
				c1++
				logrus.Infof("[mytopic - %d]: %s", c1, string(msg))
				break
			case err := <-onError1:
				logrus.Error(err)
				break
			case msg := <-onMessage2:
				c2++
				logrus.Infof("[mytopic2 - %d]: %s", c2, string(msg))
				break
			case err := <-onError2:
				logrus.Error(err)
				break
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
			hub.Publish("mytopic", msg)
		}(i, &wg)

		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()

			msg := []byte(fmt.Sprintf("message %d", i))
			hub.Publish("mytopic2", msg)
		}(i, &wg)
	}

	wg.Wait()

	<-subFinished
	close(subFinished)

	<-subFinished2
	close(subFinished2)

	<-hubFinished
	close(hubFinished)
}
