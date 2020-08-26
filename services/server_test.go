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
	"fmt"
	"testing"
)

type TestEchoCore struct {
	*EchoCore
}

func (te *TestEchoCore) Stop() {
	fmt.Println("自己实现的哦")
}

func TestEchoCore_Start(t *testing.T) {
	echo := &TestEchoCore{
		EchoCore: &EchoCore{HttpService: make(map[string]*HttpService),
			Socket:    make(map[string]*Socket),
			WebSocket: make(map[string]*WebSocket),
			Mode:      "debug",
			AppName:   "Test"},
	}
	echo.Write()
	echo.Start()
	echo.Stop()
}
