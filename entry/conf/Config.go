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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/ice-zzz/netcore/easygo/tools"
	network "github.com/ice-zzz/netcore/net"

	"github.com/BurntSushi/toml"
)

var (
	confpath = "./config"
)

// *************************************************************
// *   平台配置
// *************************************************************

type PlatformConfig struct {
	Sys       SystemConfig
	WebSocket network.WebSocketOption
	Http      network.HttpOption
}

type SystemConfig struct {
	NumCPU int `toml:"num_cpu"`
}

func ReadPlatformConfig() (*PlatformConfig, error) {
	var file *os.File

	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	platformconfig := fmt.Sprintf("%s/%s", confpath, "config.toml")
	plconfig := &PlatformConfig{}

	if isExists, _ := tools.PathExists(platformconfig); isExists {
		file, _ := os.Open(platformconfig)
		confBytes, _ := ioutil.ReadAll(file)
		if _, err := toml.Decode(string(confBytes), plconfig); err != nil {
			return plconfig, err
		}
		return plconfig, nil
	}
	return plconfig, errors.New("Config file is not exists! ")

}

func ReadConfigWithFile(path string) ([]byte, error) {

	configPath := fmt.Sprintf("%s/%s.toml", confpath, path)

	if isExists, _ := tools.PathExists(configPath); isExists {
		file, _ := os.Open(configPath)
		confBytes, _ := ioutil.ReadAll(file)
		return confBytes, nil
	}
	return nil, errors.New("Config file is not exists! ")

}

func CreatConfig() {
	var file *os.File

	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	platformconfig := "./config.toml"
	plconfig := &PlatformConfig{}

	file, _ = os.OpenFile(platformconfig, syscall.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err := toml.NewEncoder(file).Encode(plconfig); err != nil {
		fmt.Println(err.Error())
	}

}
