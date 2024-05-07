package mq

import (
	"fmt"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func TestXX(xx *testing.T) {
	server := "tcp://192.168.0.15:1883"
	clientID := "edge-1"
	username := "am"
	password := "123456"
	conf := MQConfig{
		Broker:   server,
		ClientID: clientID,
		Username: username,
		Password: password,
	}
	c := NewMQTTServer(conf)
	c.Start()
	defer c.Stop()

	s := Subscribe{
		Topic: "am/test",
		Qos:   0,
		Callback: func(c mqtt.Client, m mqtt.Message) {
			fmt.Println("am/testxxxxxxxxxxxxxx", string(m.Payload()))
		},
	}
	c.Subscribe(s)
	for {

		if err := c.Publish("am/test", 0, false, []byte("xx123456xx")); err == nil {
			fmt.Println("发布成功")
		} else {
			fmt.Println("发布失败")
		}
		time.Sleep(2 * time.Second)
	}
}
