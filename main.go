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

	"github.com/ice-zzz/netcore/services/http"
	"github.com/panjf2000/gnet"
)

func main() {
	port := 9999

	tcpServer := &http.HttpServer{}

	gnet.Serve(tcpServer,
		fmt.Sprintf("tcp://0.0.0.0:%d", port),
		gnet.WithMulticore(true),
		gnet.WithCodec(&http.HttpCode{}),
	)

}
