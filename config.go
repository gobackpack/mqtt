package mqtt

import (
	"github.com/google/uuid"
	"os"
	"strings"
	"time"
)

type Config struct {
	Host          string
	Port          string
	Username      string
	Password      string
	ClientID      string
	CleanSession  bool
	AutoReconnect bool
	Retained      bool
	KeepAlive     time.Duration
	MsgChanDept   uint
}

func NewConfig() *Config {
	host := os.Getenv("MQTT_HOST")
	if strings.TrimSpace(host) == "" {
		host = "localhost"
	}

	port := os.Getenv("MQTT_PORT")
	if strings.TrimSpace(port) == "" {
		port = "1883"
	}

	username := os.Getenv("MQTT_USERNAME")
	if strings.TrimSpace(username) == "" {
		username = "guest"
	}

	password := os.Getenv("MQTT_PASSWORD")
	if strings.TrimSpace(password) == "" {
		password = "guest"
	}

	clientId := os.Getenv("MQTT_CLIENT_ID")
	if strings.TrimSpace(clientId) == "" {
		clientId = uuid.New().String()
	}

	return &Config{
		Host:          host,
		Port:          port,
		Username:      username,
		Password:      password,
		ClientID:      clientId,
		CleanSession:  true,
		AutoReconnect: true,
		Retained:      false,
		KeepAlive:     15 * time.Second,
		MsgChanDept:   100,
	}
}
