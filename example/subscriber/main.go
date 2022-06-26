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
	sub1 := hub.Subscribe(hubCtx, "mytopic")
	sub2 := hub.Subscribe(hubCtx, "mytopic2")

	go func(ctx context.Context) {
		c1 := 0
		c2 := 0
		for {
			select {
			case msg := <-sub1.OnMessage:
				c1++
				logrus.Infof("[mytopic - %d]: %s", c1, string(msg))
				break
			case err := <-sub1.OnError:
				logrus.Error(err)
				break
			case msg := <-sub2.OnMessage:
				c2++
				logrus.Infof("[mytopic2 - %d]: %s", c2, string(msg))
				break
			case err := <-sub2.OnError:
				logrus.Error(err)
				break
			case <-ctx.Done():
				return
			}
		}
	}(hubCtx)

	logrus.Info("listening for messages...")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
