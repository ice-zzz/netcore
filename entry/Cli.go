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
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/ice-zzz/netcore/entry/conf"
)

var (
	entry *Entry
)

func Cli() {
	var isInstall = flag.Bool("install", false, "install program")
	var isRun = flag.Bool("run", false, "run program")
	flag.Parse()
	if *isInstall {
		install()
	} else if *isRun {
		run()
	}
}

func GetEntry() *Entry {
	return entry
}

func install() {
	fmt.Printf("%s", "正在读取系统数据...")
	sys := &SYSTEM{
		CPU:  GetCpuInfo(),
		DISK: GetDiskInfo(),
		NET:  GetNetInfo(),
		HOST: GetHostInfo(),
		MEM:  GetMemInfo(),
	}
	sysInfo := "./sysreport.ice"
	file, _ := os.OpenFile(sysInfo, syscall.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err := toml.NewEncoder(file).Encode(sys); err != nil {
		fmt.Println(err.Error())
	}
	file.Close()
	fmt.Printf("%s", "正在创建配置文件. 如没有预配置,请填写配置文件,如有预配置文件请覆盖.")
	// 创建配置文件
	conf.CreatConfig()

	fmt.Printf("%s", "安装完成...")

}

func run() {

	if e, err := Create(); err != nil {
		entry = e
		fmt.Printf("启动异常: %s", err.Error())
	} else {
		e.Start()
		e.ExitSignalMonitor()
	}
}
