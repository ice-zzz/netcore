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
package main

import (
	"fmt"

	"github.com/ice-zzz/netcore/entry"
)

func main() {

	if entry.Cli() == false {
		if e, err := entry.Create("./config.toml"); err != nil {
			fmt.Printf("启动异常: %s", err.Error())
		} else {
			e.Start()
			e.ExitSignalMonitor()
		}
	}

}
