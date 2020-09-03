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
package services

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/common-nighthawk/go-figure"
	"github.com/ice-zzz/netcore/services/http"
	"github.com/ice-zzz/netcore/services/websocket"
	"github.com/ice-zzz/netcore/utils"
)

type EchoConfig interface {
	Read() error
	Write() error
}

type Service interface {
	Start()
	Stop()
}

type Socket struct {
}

type EchoCore struct {
	HttpService map[string]*http.HttpService    `toml:"http_service"`
	Socket      map[string]*Socket              `toml:"socket"`
	WebSocket   map[string]*websocket.WebSocket `toml:"web_socket"`
	Mode        string                          `toml:"mode"`
	AppName     string                          `toml:"app_name"`
}

var configPath = "./config.toml"

func (echo *EchoCore) Start() {
	err := echo.Read()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	figure.NewColorFigure(echo.AppName, "big", "purple", false).Print()

}

func (echo *EchoCore) Stop() {
	panic("自己要实现哦")
}

func (echo *EchoCore) Read() error {
	var file *os.File
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	if isExists, _ := utils.PathExists(configPath); isExists {
		file, _ := os.Open(configPath)
		confBytes, _ := ioutil.ReadAll(file)
		if _, err := toml.Decode(string(confBytes), echo); err != nil {
			return err
		}
		return nil
	}
	return errors.New("Config file is not exists! ")

}

func (echo *EchoCore) Write() error {
	var file *os.File

	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()

	file, _ = os.OpenFile(configPath, syscall.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err := toml.NewEncoder(file).Encode(echo); err != nil {
		fmt.Println(err.Error())
	}
	return nil
}
