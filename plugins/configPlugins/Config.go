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
package configPlugins

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/ice-zzz/netcore/utils/filetools"
)

// *************************************************************
// *   平台配置
// *************************************************************

type Config map[string]interface{}

func (c *Config) Read(path string) error {
	var file *os.File
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	if path == "" {
		return errors.New("文件不存在! ")
	}
	platformconfig := fmt.Sprintf("%s", path)

	if isExists, _ := filetools.PathExists(platformconfig); isExists {
		file, _ := os.Open(platformconfig)
		confBytes, _ := ioutil.ReadAll(file)
		if _, err := toml.Decode(string(confBytes), c); err != nil {
			return err
		}
		return nil
	}
	return errors.New("Config file is not exists! ")

}

func (c *Config) Write(path string) error {
	var file *os.File

	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	if path == "" {
		path = "./config.toml"
	}
	platformconfig := fmt.Sprintf("%s", path)

	file, _ = os.OpenFile(platformconfig, syscall.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err := toml.NewEncoder(file).Encode(c); err != nil {
		fmt.Println(err.Error())
	}
	return nil
}
