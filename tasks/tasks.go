package tasks

import (
	"git.cddpi.com/iot/iot-edge-driver/api"
	"git.cddpi.com/iot/iot-edge-driver/common/log"
	"git.cddpi.com/iot/iot-edge-driver/datasource"
)

const (
	READ = iota
	WRITE
	ADDTASK
	DELTASK
)

type taskMsg struct {
	operation    int
	datasourceId int64
	task         *Task
}
type TasksMgr struct {
	pool        *Pool
	taskCache   map[int64][]*Task
	driver      api.Driver
	sendMqChan  chan *datasource.DataSourceNotify
	taskChannel chan taskMsg
}

func NewTaskManager(mqChan chan *datasource.DataSourceNotify, driverImp api.Driver) *TasksMgr {
	tm := &TasksMgr{}
	tm.driver = driverImp
	tm.pool = NewPool(tm.Execute)
	tm.sendMqChan = mqChan
	tm.taskChannel = make(chan taskMsg)
	tm.taskCache = make(map[int64][]*Task)
	return tm
}
func (t *TasksMgr) Run() {
	go t.pool.Run()
	for msg := range t.taskChannel {
		switch msg.operation {
		case ADDTASK:
			t.pool.FetchCoroutine(msg.task)
		case DELTASK:
			t.pool.DeleteTask(msg.datasourceId)
		default:
			continue
		}
	}

}
func (t *TasksMgr) DeleteTask(datasourceId int64) (ok bool) {
	_, ok = t.taskCache[datasourceId]
	if ok {
		t.taskChannel <- taskMsg{DELTASK, datasourceId, nil}
	}
	return
}

func (t *TasksMgr) AddTask(interval int, pointsIds []int64, datasourceId int64) {
	for _, pointId := range pointsIds {
		task := NewTask(interval, pointId, datasourceId)
		t.taskCache[datasourceId] = append(t.taskCache[datasourceId], task)
		t.taskChannel <- taskMsg{ADDTASK, datasourceId, task}
	}

}
func (t *TasksMgr) Execute(operation int, datasourceId int64, tasks []*Task) {
	switch operation {
	case READ:
		var pointsIds []int64
		for _, task := range tasks {
			pointsIds = append(pointsIds, task.pointId)
		}
		res, err := t.driver.Read(datasourceId, pointsIds)
		if err != nil {
			log.Errorf("read from driver err:%v", err)
			return
		}
		if res == nil {
			//log.Errorf("read from driver the res is nil")
			return
		}
		res.DataSourceId = datasourceId
		t.sendMqChan <- res
	case WRITE:
	default:
		return
	}
}
