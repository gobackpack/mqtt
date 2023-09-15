package main

import (
	"context"
	"fmt"
	"github.com/gobackpack/mqtt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	// listen for errors
	go func() {
		defer logrus.Warn("errors listener finished")

		for {
			select {
			case <-hubCtx.Done():
				return
			case err, ok := <-pub1.OnError:
				if !ok {
					return
				}
				logrus.Error(err)
			case err, ok := <-pub2.OnError:
				if !ok {
					return
				}
				logrus.Error(err)
			}
		}
	}()

	// pub
	wg := sync.WaitGroup{}
	delta := 1000
	wg.Add(delta * 2)

	for i := 0; i < delta; i++ {
		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()
			pub1.Publish("mytopic", mqtt.DefaultPubQoS, []byte(fmt.Sprintf("message %d", i)))
		}(i, &wg)

		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()
			pub2.Publish("mytopic2", mqtt.DefaultPubQoS, []byte(fmt.Sprintf("message %d", i)))
		}(i, &wg)
	}

	wg.Wait()

	logrus.Warn("publisher finished...")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
