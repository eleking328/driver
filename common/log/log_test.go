package log

import (
	"sync"
	"testing"
	"time"
)

func TestInfo(t *testing.T) {
	Info("Info")
}

func TestInfof(t *testing.T) {
	Infof("Infof %d", 123)
}

func TestDebug(t *testing.T) {
	Debug("Debug")
}

func TestDebugf(t *testing.T) {
	Debugf("Debugf %d", 123)
}

func TestClose(t *testing.T) {
	SetLog(true, "")
	var l sync.WaitGroup
	l.Add(5)
	r := true
	go func() {
		for r {
			Info("info x")
			Infof("%s xx", "infof")
			Debug("debug xxx")
			Debugf("%s xxx", "debugf")
			time.Sleep(1 * time.Second)
			l.Done()
		}
	}()
	l.Wait()
	r = false
	Dispose()
}
