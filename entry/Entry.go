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
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	BuildTime   string
	GitBranch   string
	CommitId    string
	GoVersion   string
	exitChannel chan os.Signal
	version     = flag.Bool("v", false, "显示版本号")
)

func Init() {
	numCPU := runtime.NumCPU()
	if numCPU < 2 {
		numCPU = 1
	} else {
		numCPU = numCPU - 1
	}
	runtime.GOMAXPROCS(numCPU)

	flag.Parse()
	if *version {
		fmt.Printf("BuildTime: %s Branch: %s Commit: %s RuntimeVersion: %s", BuildTime, GitBranch, CommitId, GoVersion)
	}
}

func ExitSignalMonitor() {
	exitChannel = make(chan os.Signal, 1)
	signal.Notify(exitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-exitChannel
}

func Stop() {
	// TODO 保存工作
	exitChannel <- syscall.SIGINT
}
