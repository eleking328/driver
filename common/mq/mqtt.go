package mq

import (
	"git.cddpi.com/iot/iot-edge-driver/common/log"

	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQConfig
//
// mq配置
type MQConfig struct {
	Broker   string `json:"broker"`   //broker addr
	ClientID string `json:"clientId"` //客户端编号
	Username string `json:"username"` //用户名
	Password string `json:"password"` //密码
}
type Subscribe struct {
	Topic    string
	Qos      byte
	Callback mqtt.MessageHandler
}
type MQTT struct {
	client     mqtt.Client
	subscribes map[string]Subscribe
	config     MQConfig
}

func NewMQTTServer(conf MQConfig) *MQTT {
	return &MQTT{
		config:     conf,
		subscribes: make(map[string]Subscribe),
	}
}

func (p *MQTT) connectHandler(client mqtt.Client) {
	log.Debugf("broker=%s connected", p.config.Broker)
	//重新订阅
	for topic, item := range p.subscribes {
		client.Subscribe(topic, item.Qos, item.Callback)
	}
}

func (p *MQTT) connectLostHandler(client mqtt.Client, err error) {
	log.Debugf("Connect lost: %v", err)
}

// Publish 发布消息
func (p *MQTT) Publish(topic string, qos byte, retained bool, payload []byte) error {
	if token := p.client.Publish(topic, qos, retained, payload); token.Wait() && token.Error() == nil {
		return nil
	} else {
		return token.Error()
	}
}

// Subscribe 订阅消息
func (p *MQTT) Subscribe(items ...Subscribe) {
	if p.client != nil && p.client.IsConnected() && p.client.IsConnectionOpen() {
		//先退订
		for _, item := range items {
			p.client.Unsubscribe(item.Topic)
		}
		//再订阅
		for _, item := range items {
			p.subscribes[item.Topic] = item
			p.client.Subscribe(item.Topic, item.Qos, item.Callback)
		}

	} else {
		for _, item := range items {
			p.subscribes[item.Topic] = item
		}
	}
}

// Unsubscribe 退订
func (p *MQTT) Unsubscribe(topics ...string) {
	p.client.Unsubscribe(topics...)
}
func (p *MQTT) Stop() {
	if p == nil {
		return
	}
	if p.client == nil {
		return
	}
	if !p.client.IsConnected() {
		return
	}
	p.client.Disconnect(1000)
}
func (p *MQTT) Start() *MQTT {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(p.config.Broker)
	opts.SetClientID(p.config.ClientID)
	opts.SetUsername(p.config.Username)
	opts.SetPassword(p.config.Password)
	//设置自动重联
	opts.SetAutoReconnect(true)
	//5秒心跳
	opts.SetKeepAlive(5 * time.Second)
	//opts.SetDefaultPublishHandler(p.messageHandler)
	opts.OnConnect = p.connectHandler
	opts.OnConnectionLost = p.connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	p.client = client
	log.Debugf(" mqtt server started	: %s", p.config.Broker)
	return p
}
