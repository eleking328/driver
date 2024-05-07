package log

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"git.cddpi.com/iot/iot-edge-driver/common/utils"
)

var (
	debug     = false
	path      = ""
	exit      = make(chan bool)
	d_channel = make(chan string)
	i_channel = make(chan string)
	e_channel = make(chan string)
	w_channel = make(chan string)
)

func init() {
	go func() {
		defer func() {
			fmt.Println("exit")
			close(d_channel)
			close(i_channel)
			close(e_channel)
			close(w_channel)
		}()
		checkFileName := time.NewTicker(1 * time.Second)
		fileName := utils.NewDateTime().FormatDate()
	L:
		for {
			select {
			case <-checkFileName.C:
				if len(path) > 0 {
					name := utils.NewDateTime().FormatDate()
					if fileName != name {
						//change output
						setOutput(name)
					}
				}
			case <-exit:
				break L
			case msg, open := <-d_channel:
				if debug && open {
					log.Println("[DEBUG]", msg)
				}
			case msg, open := <-i_channel:
				if open {
					log.Println("[INFO]", msg)
				}
			case msg, open := <-e_channel:
				if open {
					log.Println("[ERROR]", msg)
				}
			case msg, open := <-w_channel:
				if open {
					log.Println("[WARN]", msg)
				}
			}
		}
	}()
}

func SetLog(d bool, p string) {
	debug = d
	if path == p {
		return
	}
	if len(p) == 0 {
		return
	}
	path = p
	if len(path) > 0 {
		fileName := utils.NewDateTime().FormatDate()
		setOutput(fileName)
	}
}

func IsDebug() bool {
	return debug
}

func getPath() string {
	if strings.HasSuffix(path, "/") {
		return path
	}
	return path + "/"
}
func setOutput(name string) {
	//check path
	p := getPath()
	has, _ := utils.PathExists(p)
	if !has {
		os.MkdirAll(p, os.ModePerm)
	}
	name = fmt.Sprintf("%s%s.log", getPath(), name)
	file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Println("create error", err)
	}
}
func Info(v ...interface{}) {
	message := fmt.Sprint(v...)
	i_channel <- message
}
func Infof(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	i_channel <- message
}
func Debug(v ...interface{}) {
	message := fmt.Sprint(v...)
	d_channel <- message
}
func Debugf(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	d_channel <- message
}

func Error(v ...interface{}) {
	message := fmt.Sprint(v...)
	e_channel <- message
}
func Errorf(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	e_channel <- message
}
func Warn(v ...interface{}) {
	message := fmt.Sprint(v...)
	w_channel <- message
}
func Warnf(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	w_channel <- message
}
func Dispose() {
	exit <- true
}
