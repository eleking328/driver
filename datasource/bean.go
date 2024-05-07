package datasource

import (
	"encoding/json"
	"time"

	"github.com/eleking328/driver-sdk/common/log"
)

// DeviceChannle device_channel
type DeviceChannel struct {
	CloudID int64  `json:"cloudId"` //云端主键  0 表示未同步
	Config  string `json:"config"`  // 配置信息
}

// DeviceChannle device_channel
type DeviceProduct struct {
	ID         int64     `json:"id"`                  //自增涨主键
	CloudID    int64     `json:"cloudId"`             //云端主键  0 表示未同步
	Code       string    `json:"code"`                //code
	Name       string    `json:"name"`                //通道名称
	Secret     string    `json:"secret"`              //密钥
	Type       int       `json:"type"`                //产品类型
	Timeout    int64     `json:"timeout"`             //超时(秒)
	Agreement  int       `json:"agreement"`           //协议类型 1:MQTT;2:MODBUS;
	Func       string    `json:"funcs"`               // 功能点
	Event      string    `json:"events"`              // 事件
	Cmd        string    `json:"cmds"`                // 命令
	Extension  string    `json:"extension,omitempty"` //扩展信息  2021-06-08增加
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

func (p DeviceProduct) Functions() (items []Parameter) {
	if len(p.Func) == 0 {
		return
	}
	err := json.Unmarshal([]byte(p.Func), &items)
	if err != nil {
		log.Errorf("Unmarshal product=%s Func error %v", p.Code, err)
	}
	return
}
func (p DeviceProduct) Events() (items []Event) {
	if len(p.Event) == 0 {
		return
	}
	err := json.Unmarshal([]byte(p.Event), &items)
	if err != nil {
		log.Errorf("Unmarshal product=%s Event error %v", p.Code, err)
	}
	return
}
func (p DeviceProduct) Commands() (items []Command) {
	if len(p.Cmd) == 0 {
		return
	}
	err := json.Unmarshal([]byte(p.Cmd), &items)
	if err != nil {
		log.Errorf("Unmarshal product=%s Event error %v", p.Code, err)
	}
	return
}

// DeviceChannle device_channel
type DeviceInfo struct {
	ID         int64     `json:"id"`        //自增涨主键
	Code       string    `json:"code"`      //code
	Secret     string    `json:"secret"`    //密钥
	Type       int       `json:"type"`      //设备类型 1:MQTT;2:Modbus;3:Video;4:dpi-TCP;
	ProductID  int64     `json:"productId"` //产品
	ChannelID  int64     `json:"channelId"` //通道ID
	SlaveID    int       `json:"slaveId"`   //从机号
	Interval   int64     `json:"interval"`  //数据采集间隔 秒
	Config     string    `json:"config"`    //设备配置信息  2021-08-12增加
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

type DataSourceInfo struct {
	DataPoint  string `json:"datapoint"`
	Dsec       string `json:"describe"`
	DriverID   int64  `json:"driverId"`
	CloudID    int64  `json:"id"` //datasourc cloud id
	Name       string `json:"name"`
	Properties string `json:"properties"`
	CreateTime string `json:"createTime"`
}
type Device struct {
	Device  DeviceInfo    `json:"device"`
	Product DeviceProduct `json:"product"`
}

// DeviceFuncNotify Device 设备信息
type DataSourceNotify struct {
	DataSourceId int64                 `json:"dsCloudId"`
	DriverId     int64                 `json:"dCloudId"`
	Func         map[int64]interface{} `json:"func"` //功能点
	Timestamp    int64                 `json:"time"` //更新时间
}
