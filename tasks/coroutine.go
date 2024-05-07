package tasks

import (
	"sync"
	"time"

	"github.com/eleking328/driver-sdk/common/log"

	"github.com/rs/xid"
)

type Coroutine struct {
	TaskCache map[int64][]*Task
	ticker    *time.Ticker
	//EntryChannel chan *Task
	execute      func(int, int64, []*Task)
	interval     int
	coId         string
	mu           sync.RWMutex
	total        int
	quit         chan struct{}
	recycleTimer <-chan time.Time
}

func NewCoroutine(interval int, execute func(int, int64, []*Task)) *Coroutine {
	return &Coroutine{
		TaskCache:    make(map[int64][]*Task),
		interval:     interval,
		ticker:       time.NewTicker(time.Duration(interval) * time.Millisecond),
		recycleTimer: make(<-chan time.Time),
		execute:      execute,
		coId:         xid.New().String(),
		mu:           sync.RWMutex{},
		total:        0,
		quit:         make(chan struct{}, 1),
	}
}

func (t *Coroutine) AddTask(task *Task) {
	t.mu.Lock()
	t.TaskCache[task.DatasourceID] = append(t.TaskCache[task.DatasourceID], task)
	t.total++
	t.mu.Unlock()
	//t.recycleTimer.Stop()
}
func (t *Coroutine) DelTask(datasourceId int64) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.total -= len(t.TaskCache[datasourceId])
	delete(t.TaskCache, datasourceId)
	return t.total
}
func (t *Coroutine) GetTaskCount() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.total
}
func (t *Coroutine) Stop() {
	t.ticker.Stop()
	t.quit <- struct{}{}
}
func (t *Coroutine) StartRecycle() chan struct{} {
	t.recycleTimer = time.After(30 * time.Second)
	return t.quit
}
func (t *Coroutine) Run() {
	for {
		select {
		case <-t.ticker.C:
			t.mu.RLock()
			var count = 0
			for datasourceId, tasks := range t.TaskCache {
				count += len(tasks)
				log.Debugf("携程【%d】【%s】执行datasource[%d]任务个数：【%d】", t.interval, t.coId, datasourceId, len(tasks))
				t.execute(READ, datasourceId, tasks)
			}
			t.mu.RUnlock()
			log.Debugf("携程【%d】【%s】执行总任务个数：【%d】", t.interval, t.coId, count)
		case <-t.recycleTimer:
			t.mu.RLock()
			left := t.total
			t.mu.RUnlock()
			log.Infof("携程【%d】【%s】开始回收资源，目前剩余任务数：【%d】", t.interval, t.coId, left)
			if left == 0 {
				t.Stop()
				return
			}

		}
	}
}
