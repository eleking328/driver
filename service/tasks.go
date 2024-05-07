package service

import (
	"github.com/eleking328/driver-sdk/api"
	"github.com/eleking328/driver-sdk/common/log"
	"github.com/eleking328/driver-sdk/datasource"
	"github.com/eleking328/driver-sdk/tasks"
)

type TasksService struct {
	driver  api.Driver
	taskMgr *tasks.TasksMgr
}

var Tasks *TasksService

func NewTaskService(driverImp api.Driver) *TasksService {
	Tasks = &TasksService{
		taskMgr: tasks.NewTaskManager(EdgeMQ.FuncNotify, driverImp),
		driver:  driverImp,
	}
	return Tasks
}

func (t TasksService) Start() {
	go t.taskMgr.Run()
}
func (TasksService) Stop() {

}
func (t TasksService) DeleteTask(datasouceId int64) {
	t.taskMgr.DeleteTask(datasouceId)
}
func (t TasksService) PutDataSource(info datasource.DataSourceInfo) {
	if len(info.DataPoint) == 0 {
		log.Debugf("Datasource[%d]没有点位配置\n", info.CloudID)
		return
	}
	if t.taskMgr.DeleteTask(info.CloudID) {
		log.Debugf("更新Datasource[%d]配置\n", info.CloudID)
	}
	err := t.driver.CreateChannel(info.CloudID, []byte(info.Properties))
	if err != nil {
		log.Errorf("创建数据源通道失败：%v", err)
		return
	}
	points, err := t.driver.ParsePointsProperties(info.CloudID, []byte(info.DataPoint))
	if err != nil {
		panic(err)
	}
	if points == nil {
		t.driver.SubPointsData(info.CloudID, EdgeMQ.FuncNotify)
		return
	}
	for interval, points := range points {
		t.taskMgr.AddTask(interval, points, info.CloudID)
	}

}
