package service

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/eleking328/driver-sdk/common/log"
	"github.com/eleking328/driver-sdk/common/mq"
	"github.com/eleking328/driver-sdk/config"
	"github.com/eleking328/driver-sdk/datasource"

	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var total int32 = 0

const (
	UPDATE_ACTION = "update"
	DELETE_ACTION = "delete"
)

type EdgeMQTTService struct {
	driverId    string
	manageCode  string
	m           *mq.MQTT
	mq          mq.MQConfig
	quit        chan bool
	FuncNotify  chan *datasource.DataSourceNotify
	onlineTimer *time.Timer
}

// EdgeMQ EdgeMQ
var EdgeMQ *EdgeMQTTService

func NewEdgeMQTT(conf config.Config) *EdgeMQTTService {
	EdgeMQ = &EdgeMQTTService{
		driverId:   conf.DriverId,
		manageCode: conf.ManageCode,
		mq:         conf.EdgeMQ,
		quit:       make(chan bool),
		FuncNotify: make(chan *datasource.DataSourceNotify),
	}
	return EdgeMQ
}
func (p *EdgeMQTTService) Start() {
	p.m = mq.NewMQTTServer(p.mq)
	p.subscribe()
	p.m.Start()
	p.onlineTimer = p.online()
FOR:
	for {
		select {
		case <-p.quit:
			break FOR
		case m, ok := <-p.FuncNotify:
			if !ok {
				break FOR
			}
			go p.notifyEdgeFunc(m)
		}
	}
}

// 上线
func (p *EdgeMQTTService) online() *time.Timer {
	return time.AfterFunc(5*time.Second, func() {
		topic := fmt.Sprintf("edge/driver/%s/online/%s", p.manageCode, p.driverId)
		p.m.Publish(topic, 0, false, []byte("{}"))
	})
}
func (p *EdgeMQTTService) notifyEdgeFunc(item *datasource.DataSourceNotify) {
	if len(item.Func) == 0 {
		return
	}
	item.DriverId, _ = strconv.ParseInt(p.driverId, 10, 64)
	payload, err := json.Marshal(item)
	if err != nil {
		log.Debugf("notify error %v", err)
		return
	}
	topic := fmt.Sprintf("edge/func/%s/%s/%d", p.manageCode, p.driverId, item.DataSourceId)
	//log.Debug("topic", topic, "payload=>", string(payload))
	atomic.AddInt32(&total, 1)
	err = p.m.Publish(topic, 1, false, payload)
	if err != nil {
		log.Errorf("notify edge func err:%v", err)
	}

}

// 2022/3/18 topic增加action字段，区分更新和删除
// 设备产品变更 edge/product/{action}
// 设备设备变更 edge/device/{action}
// 设备通道变更 edge/channel/{action}
func (p *EdgeMQTTService) subscribe() {

	datasourceSub := mq.Subscribe{
		Topic:    fmt.Sprintf("edge/datasource/+/%s/%s", p.manageCode, p.driverId), //2022/3/14 modify
		Qos:      0,
		Callback: p.datasourceChange,
	}

	cmdSub := mq.Subscribe{
		Topic:    "edge/cmd/+/+", //edge/cmd/}/{datasourceId}
		Qos:      0,
		Callback: p.doCMD,
	}

	p.m.Subscribe(datasourceSub, cmdSub)
}

// parseCmmondTopic parse cmd topic
func parseCmmondTopic(topic string) (datasourceId int64, err error) {
	items := strings.Split(topic, "/")
	if len(items) < 4 {
		err = errors.New("错误的topic")
		return
	}
	datasourceId, err = strconv.ParseInt(items[3], 10, 64)
	return
}

func (p *EdgeMQTTService) doCMD(c mqtt.Client, m mqtt.Message) {
	log.Debugf("接收到来自topic[%s]的命令：%s", m.Topic(), string(m.Payload()))
	datasourceId, err := parseCmmondTopic(m.Topic())
	if err != nil {
		log.Errorf("命令反序列化失败：%v", err)
		return
	}
	// var cmdData = make(map[int64]interface{})
	// err = json.Unmarshal(m.Payload(), &cmdData)
	// if err != nil {
	// 	log.Errorf("命令反序列化失败：%v", err)
	// 	return
	// }

	Tasks.driver.Write(datasourceId, m.Payload())

}

func (p *EdgeMQTTService) Stop() {
	p.m.Stop()
}
func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}
func (p *EdgeMQTTService) datasourceChange(c mqtt.Client, m mqtt.Message) {
	log.Debugf("接收到来自topic[%s]的数据源配置：%s", m.Topic(), string(m.Payload()))
	if strings.Contains(m.Topic(), DELETE_ACTION) {
		datasourceId := BytesToInt64(m.Payload())
		Tasks.DeleteTask(datasourceId)
		return
	}
	var items []datasource.DataSourceInfo
	var err error
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()
	err = json.Unmarshal(m.Payload(), &items)
	if err != nil {
		err = errors.New("dataSourceChange Unmarshal error=" + err.Error())
		log.Errorf("配置反序列化失败：%v", err)
		return
	}
	driverId, _ := strconv.ParseInt(p.driverId, 10, 64)
	for _, item := range items {
		if item.DriverID != driverId {
			log.Errorf("不匹配得驱动信息")
			continue
		}
		Tasks.PutDataSource(item)
	}

}
