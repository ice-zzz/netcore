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
package hardwareService

import (
	"log"
	"testing"
	"time"
)

func TestSYSTEM_Start(t *testing.T) {
	s := &SYSTEM{}
	go s.Start()
	time.Sleep(time.Second * 3)
	for {
		up, down := s.GetNetSpeed("en0")
		log.Printf("Up: %s KB    Down: %s KB  \n", up, down)
		time.Sleep(1 * time.Second)

	}
}
