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
	h           = flag.Bool("h", false, "显示帮助")
	daemon      = flag.Bool("d", false, "Daemon模式")
)

func Entry() {
	flag.Parse()
	args := os.Args[1:]
	if len(args) <= 2 {
		usage()
		os.Exit(1)
	} else {
		switch args[0] {
		case "start":

			break
		case "stop":

			break
		case "restart":

			break
		default:
			usage()
			os.Exit(1)
		}
	}

	if *h {
		flag.Usage()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `icecore version: icecore/0.2.1
Usage: icecore start|stop|restart [-hvVtTq] [-s signal]

Options:
`)
	flag.PrintDefaults()
}

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
