package driver

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.cddpi.com/iot/iot-edge-driver/api"
	"git.cddpi.com/iot/iot-edge-driver/common/log"
	"git.cddpi.com/iot/iot-edge-driver/config"
	"git.cddpi.com/iot/iot-edge-driver/service"
)

func StartDriver(d api.Driver) {
	//read config
	c, err := config.InitAppConfig()
	if err != nil {
		log.Infof("init config error - %s", err.Error())
		return
	}
	c.Driver = d
	log.SetLog(c.Debug, c.Log)
	//xlog.Infof("config-%+v", c)
	//启动边缘MQTT服务
	mqttEdge := service.NewEdgeMQTT(c)
	defer mqttEdge.Stop()
	go mqttEdge.Start()

	taskMgr := service.NewTaskService(c.Driver)
	defer taskMgr.Stop()
	go taskMgr.Start()

	//结束
	ExitListen(context.WithCancel(context.Background()))
}

// ExitListen 退出监听器
func ExitListen(ctx context.Context, cancel context.CancelFunc) {
	sigs := make(chan os.Signal)
	//设置要接收的信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Kill, os.Interrupt)
	select {
	case s := <-sigs:
		//外部信号
		cancel()
		//设置延时退出
		log.Info("外部信号", s)
		time.Sleep(1 * time.Second)
	case <-ctx.Done():
		//内部结束
		//直接退出
		log.Info("内部信号", ctx.Err())
	}
	log.Info("进程终止")
	os.Exit(1)
}
