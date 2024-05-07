package tasks

import (
	"git.cddpi.com/iot/iot-edge-driver/common/log"
	"time"
)

type Pool struct {
	cocache         map[int]map[string]*Coroutine
	execute         func(int, int64, []*Task)
	coChannel       chan *Coroutine
	datasouceMapCos map[int64][]*Coroutine
}

const MAX_COROUTINE_CAP = 1000

// 创建一个协程池
func NewPool(execute func(int, int64, []*Task)) *Pool {

	return &Pool{
		execute:         execute,
		cocache:         make(map[int]map[string]*Coroutine),
		coChannel:       make(chan *Coroutine),
		datasouceMapCos: make(map[int64][]*Coroutine),
	}
}

func (p *Pool) findIdleCoroutine(interval int) *Coroutine {
	cos, ok := p.cocache[interval]
	if !ok {
		return nil
	}
	for _, co := range cos {
		if co.GetTaskCount() < MAX_COROUTINE_CAP {
			return co
		}
	}
	return nil
}
func (p *Pool) FetchCoroutine(task *Task) {
	co := p.findIdleCoroutine(task.Interval)
	if co != nil {
		co.AddTask(task)
	} else {
		if p.cocache[task.Interval] == nil {
			p.cocache[task.Interval] = make(map[string]*Coroutine)
		}
		newco := NewCoroutine(task.Interval, p.execute)
		newco.AddTask(task)
		log.Infof("创建定时间隔【%d】的新携程:【 %s 】", task.Interval, newco.coId)

		p.cocache[task.Interval][newco.coId] = newco
		p.datasouceMapCos[task.DatasourceID] = append(p.datasouceMapCos[task.DatasourceID], newco)
		p.coChannel <- newco
	}
}

func (p *Pool) DeleteTask(datasourceId int64) {
	for _, co := range p.datasouceMapCos[datasourceId] {
		left := co.DelTask(datasourceId)
		if left == 0 {
			quit := co.StartRecycle()
			go func(co *Coroutine) {
				<-quit
				p.RecycleCo(co)
			}(co)
		}
	}
}
func (p *Pool) RecycleCo(co *Coroutine) {
	delete(p.cocache[co.interval], co.coId)
}

// 让协程池Pool开始工作
func (p *Pool) Run() {
	for {
		select {
		case co := <-p.coChannel:
			go co.Run()
		default:
			var total int = 0
			for _, cos := range p.cocache {
				total += len(cos)

			}
			log.Debugf("正在运行携程数量:[%d]", total)
			time.Sleep(1 * time.Second)
		}
	}
}
