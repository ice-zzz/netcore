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
	"fmt"
	"testing"
)

func TestConfig_Read(t *testing.T) {
	c := &Config{}
	if err := c.Read("./config.toml"); err != nil {
		fmt.Printf("%s", err.Error())
	}
	fmt.Println(c)
}

func TestConfig_Write(t *testing.T) {
	c := &Config{}

	(*c)["aaa"] = 11123123123
	(*c)["bbb"] = 11123123123
	(*c)["fff"] = 11123123123
	(*c)["666"] = 11123123123
	(*c)["123123"] = 11123123123

	if err := c.Write("./config.toml"); err != nil {
		fmt.Printf("%s", err.Error())
	}

	if err := c.Read("./config.toml"); err != nil {
		fmt.Printf("%s", err.Error())
	}
	fmt.Println(c)
}
