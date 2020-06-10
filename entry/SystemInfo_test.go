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
	"encoding/json"
	"log"
	"testing"
	"time"
)

func TestSYSTEM_Start(t *testing.T) {
	s := &SYSTEM{}
	go s.Start()
	time.Sleep(time.Second * 3)
	sjson, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(sjson))
	for {

		log.Printf("Total: %d MB    Used: %d MB    Free: %d MB    Percent: %f%%   \n", s.MEM.Total, s.MEM.Used, s.MEM.Free, s.MEM.UsedPercent)
		time.Sleep(time.Second)

	}
}
