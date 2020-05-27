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
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ice-zzz/netcore/easygo/logs"
	"github.com/ice-zzz/netcore/entry/conf"
	network "github.com/ice-zzz/netcore/net"
	"github.com/shirou/gopsutil/load"
)

var (
	logger *logs.Logger
)

type Entry struct {
	exitChannel chan os.Signal
	pConfig     *conf.PlatformConfig

	httpService      *network.HttpServer
	webSocketService *network.WebSocketServer
}

func Create(path string) (entry *Entry, err error) {
	entry = &Entry{}
	entry.pConfig, err = conf.ReadPlatformConfig(path)
	numCPU := runtime.NumCPU()
	if entry.pConfig.Sys.NumCPU == 0 {
		if numCPU < 2 {
			numCPU = 1
		} else {
			numCPU = numCPU - 1
		}
	} else {
		numCPU = entry.pConfig.Sys.NumCPU
	}
	runtime.GOMAXPROCS(numCPU)

	if err != nil {
		return nil, err
	}
	// 初始化Entry日志
	logger = logs.New(logs.LogOption{
		WriteToFile: false,
		LogFilePath: "",
		ZipTime:     0,
	})
	// 初始化系统监控服务
	go func() {
		for {
			disk := GetDiskInfo()
			if disk.UsedPercent >= float64(85) {
				logger.Info("磁盘快满了")
			}
			v, _ := load.Avg()
			tCpus := float64(numCPU)

			if v.Load1 > tCpus && v.Load5 < tCpus && v.Load15 < tCpus {
				logger.Info("服务器波动, 短期堵塞预警")
			} else if v.Load1 > tCpus && v.Load5 > tCpus && v.Load15 < tCpus {
				logger.Info("服务器压力警告, 堵塞预警")
			} else if v.Load1 > tCpus && v.Load5 > tCpus && v.Load15 > tCpus {
				logger.Info("服务器严重压力警告, 已经堵塞很久了")
			} else if v.Load1 < tCpus && v.Load5 > tCpus && v.Load15 > tCpus {
				logger.Info("堵塞正在缓解,请保持关注")
			} else {
				logger.Info("服务器正常")
			}

			time.Sleep(15 * time.Second)
		}

	}()

	// 初始化http服务
	entry.httpService = network.CreateHttp(entry.pConfig.Http)

	// 初始化websocket服务
	entry.webSocketService = network.CreateWebSocket(entry.pConfig.WebSocket)

	return entry, nil
}

func (e *Entry) Start() {

	go e.httpService.Start()
	go e.webSocketService.Start()

	logger.Info("服务器启动完成... \n")
	// rendertext(writer, fmt.Sprintf("服务器启动完成"), ct.Green)
	// 好了我累了,休息了

}

func (e *Entry) GetHttp() *network.HttpServer {
	return e.httpService
}

func (e *Entry) GetWebSocket() *network.WebSocketServer {
	return e.webSocketService
}

func (e *Entry) ExitSignalMonitor() {
	e.exitChannel = make(chan os.Signal, 1)
	signal.Notify(e.exitChannel, syscall.SIGINT, syscall.SIGTERM)
	s := <-e.exitChannel
	logger.Errorf("接收到%s消息, 停止服务...\n", s.String())
}

func (e *Entry) Stop() {
	// TODO 保存工作
	logger.Info("正在退出...\n")
	e.exitChannel <- syscall.SIGINT
}
