package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/eleking328/driver-sdk/api"
	"github.com/eleking328/driver-sdk/common/mq"
)

func init() {
	log.Println("init config")
	//log
	log.SetFlags(log.Ldate | log.Ltime)
}

type AccrditConfig struct {
	License   string `json:"license"`
	PublicKey string `json:"publicKey"`
}

// Config config
type Config struct {
	Debug      bool        `json:"debug"`
	ManageCode string      `json:"manageCode"` //manage code
	DriverId   string      `json:"driverId"`
	Log        string      `json:"log"` //log path
	EdgeMQ     mq.MQConfig `json:"edgemq"`
	Driver     api.Driver
}

// InitAppConfig config
func InitAppConfig() (c Config, err error) {
	c = Config{
		Debug:      false,
		Log:        "",
		EdgeMQ:     mq.MQConfig{},
		ManageCode: os.Getenv("EDGE_CODE"),
		DriverId:   os.Getenv("DRIVER_ID"),
	}

	flag.BoolVar(&c.Debug, "debug", false, "is debug")

	flag.StringVar(&c.EdgeMQ.Broker, "broker", "tcp://192.168.0.15:1883", "mqtt服务端地址，例如：tcp://127.0.0.1:1883")
	flag.StringVar(&c.EdgeMQ.ClientID, "client", "", "mqtt客户端ID")
	flag.StringVar(&c.EdgeMQ.Username, "user", "test", "mqtt用户名")
	flag.StringVar(&c.EdgeMQ.Password, "password", "123456", "mqtt密码")
	flag.StringVar(&c.DriverId, "driverId", "", "驱动ID，来自于物联网云平台")
	flag.StringVar(&c.ManageCode, "manageCode", "", "所在边缘计算节点的设备code")
	flag.Parse()
	if c.ManageCode == "" {
		c.ManageCode = "U2kRJtfPwub2fhsjwbtBpu"
		c.DriverId = "185507477080005"
	}
	if len(c.EdgeMQ.ClientID) == 0 {
		c.EdgeMQ.ClientID = fmt.Sprintf("edge-%s-%s-%d", c.ManageCode, c.DriverId, time.Now().Unix())
	}
	return
}
