package main

import (
	"context"
	"github.com/gobackpack/mqtt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
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

	// sub
	sub1 := hub.Subscribe(hubCtx, "mytopic", mqtt.DefaultSubQoS)
	sub2 := hub.Subscribe(hubCtx, "mytopic2", mqtt.DefaultSubQoS)

	// handle messages and errors for sub1
	go func(ctx context.Context, sub1 *mqtt.Subscriber) {
		defer logrus.Warn("sub1 message handler and error listener stopped")

		c := 0
		for {
			select {
			case msg, ok := <-sub1.OnMessage:
				if !ok {
					return
				}
				c++
				logrus.Infof("[mytopic - %d]: %s", c, string(msg))
			case err, ok := <-sub1.OnError:
				if !ok {
					return
				}
				logrus.Error(err)
			case <-ctx.Done():
				return
			}
		}
	}(hubCtx, sub1)

	// handle messages and errors for sub2
	go func(ctx context.Context, sub2 *mqtt.Subscriber) {
		defer logrus.Warn("sub2 message handler and error listener stopped")

		c := 0
		for {
			select {
			case msg, ok := <-sub2.OnMessage:
				if !ok {
					return
				}
				c++
				logrus.Infof("[mytopic2 - %d]: %s", c, string(msg))
			case err, ok := <-sub2.OnError:
				if !ok {
					return
				}
				logrus.Error(err)
			case <-ctx.Done():
				return
			}
		}
	}(hubCtx, sub2)

	logrus.Info("listening for messages...")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
