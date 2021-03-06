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

//
// func RunSMON() {
// 	go func() {
// 		for {
// 			disk := GetDiskInfo()
// 			if disk.UsedPercent >= float64(85) {
// 				logger.Info("磁盘快满了")
// 			}
// 			v, _ := load.Avg()
// 			tCpus := float64(runtime.NumCPU())
//
// 			if v.Load1 > tCpus && v.Load5 < tCpus && v.Load15 < tCpus {
// 				logger.Info("服务器波动\n")
// 			} else if v.Load1 > tCpus && v.Load5 > tCpus && v.Load15 < tCpus {
// 				logger.Info("服务器压力警告\n")
// 			} else if v.Load1 > tCpus && v.Load5 > tCpus && v.Load15 > tCpus {
// 				logger.Info("服务器严重压力警告\n")
// 			} else if v.Load1 < tCpus && v.Load5 > tCpus && v.Load15 > tCpus {
// 				logger.Info("堵塞正在缓解\n")
// 			}
//
// 			time.Sleep(15 * time.Second)
// 		}
//
// 	}()
// }
//
// func ExportSystemInfo() {
// 	fmt.Printf("%s", "正在读取系统数据...")
// 	sys := &SYSTEM{
// 		CPU:  GetCpuInfo(),
// 		DISK: GetDiskInfo(),
// 		NET:  GetNetInfo(),
// 		HOST: GetHostInfo(),
// 		MEM:  GetMemInfo(),
// 	}
// 	sysInfo := "./sysreport"
// 	file, _ := os.OpenFile(sysInfo, syscall.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
// 	defer func() {
// 		if file != nil {
// 			_ = file.Close()
// 		}
// 	}()
// 	if err := toml.NewEncoder(file).Encode(sys); err != nil {
// 		fmt.Println(err.Error())
// 	}
//
// 	fmt.Printf("%s", "正在创建配置文件. 如没有预配置,请填写配置文件,如有预配置文件请覆盖.")
// 	// 创建配置文件
// 	// conf.CreatConfig()
//
// 	fmt.Printf("%s", "安装完成...")
//
// }
