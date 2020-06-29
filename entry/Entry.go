/**
 *     ______                 __
 *    /\__  _\               /\ \
 *    \/_/\ \/     ___     __\ \ \         __      ___     ___     __
 *       \ \ \    / ___\ / __ \ \ \  __  / __ \  /  _  \  / ___\ / __ \
 *        \_\ \__/\ \__//\  __/\ \ \_\ \/\ \_\ \_/\ \/\ \/\ \__//\  __/
 *        /\_____\ \____\ \____\\ \____/\ \__/ \_\ \_\ \_\ \____\ \____\
 *        \/_____/\/____/\/____/ \/___/  \/__/\/_/\/_/\/_/\/____/\/____/
 *
 *
 *                                                                    @寒冰
 *                                                            www.icezzz.cn
 *                                                     hanbin020706@163.com
 */
package entry

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/ice-zzz/netcore/service"
)

type Entry struct {
	service.Entity
	exitChannel chan os.Signal
	services    map[string]service.Service
}

func (e *Entry) AddService(service service.Service) {
	if _, ok := e.services[service.GetServiceName()]; ok {
		service.Stop()
	}
	e.services[service.GetServiceName()] = service
}

func Create() (entry *Entry) {
	entry = &Entry{}

	numCPU := runtime.NumCPU()
	if numCPU < 2 {
		numCPU = 1
	} else {
		numCPU = numCPU - 1
	}
	runtime.GOMAXPROCS(numCPU)

	entry.services = make(map[string]service.Service)
	return entry
}

func (e *Entry) Start() {
	for _, srv := range e.services {
		if srv.IsRunning() == false {
			go srv.Start()
			fmt.Printf("服务 < %s > 已启动,端口:%d \n", srv.GetServiceName(), srv.GetPort())
			srv.SetRunningStatus(true)
		}
	}

	// logger.Info("服务器启动完成... \n")
	// 好了我累了,休息了

}

func (e *Entry) GetService(serviceName string) service.Service {
	if s, ok := e.services[serviceName]; ok {
		return s
	}
	return nil
}

func (e *Entry) ExitSignalMonitor() {
	e.exitChannel = make(chan os.Signal, 1)
	signal.Notify(e.exitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-e.exitChannel

}

func (e *Entry) Stop() {
	// TODO 保存工作
	e.exitChannel <- syscall.SIGINT
}
