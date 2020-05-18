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
package logs

import (
	"fmt"
	"os"
	"testing"

	"github.com/segmentio/ksuid"
)

func TestGenerateID(t *testing.T) {
	ar := make(map[string]int)
	for i := 0; i < 10000000; i++ {
		id := ksuid.New().String()
		if _, ok := ar[id]; ok {
			fmt.Printf("出现重复id: %s \n", id)
			return
		}
		ar[id] = 1
	}
	fmt.Printf("没有重复ID出现\n")

}

type aaa struct {
	eeee int
}

func TestObjectCpoy(t *testing.T) {
	a := &aaa{
		eeee: 123,
	}
	b := a
	a = &aaa{
		eeee: 777,
	}
	fmt.Printf("%d \n", b.eeee)

}

func TestFilePath(t *testing.T) {
	file, err := os.OpenFile("file.aaa", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	if file != nil {
		fmt.Printf("%s", file.Name())
		fmt.Printf("%s", os.Rename(file.Name(), "uuuu.ppp"))
	}

}
