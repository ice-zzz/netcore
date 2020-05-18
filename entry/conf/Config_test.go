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
package conf

import (
	"fmt"
	"os"
	"syscall"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestReadPlatformConfig(t *testing.T) {
	platformconfig := "./config.toml"
	plconfig := &PlatformConfig{}

	file, _ := os.OpenFile(platformconfig, syscall.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err := toml.NewEncoder(file).Encode(plconfig); err != nil {
		fmt.Println(err.Error())
	}

}
