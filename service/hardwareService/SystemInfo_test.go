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
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/google/gopacket/pcap"
)

var (
	downStreamDataSize = 0 // 单位时间内下行的总字节数
	upStreamDataSize   = 0 // 单位时间内上行的总字节数

)

func TestCreateNetInfo(t *testing.T) {

	// devices, err := pcap.FindAllDevs()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// // Find exact device
	// // 根据网卡名称从所有网卡中取到精确的网卡
	// var device pcap.Interface
	// for _, d := range devices {
	// 	if d.Name == "en0" {
	// 		device = d
	// 	}
	// }
	//
	// // 根据网卡的ipv4地址获取网卡的mac地址，用于后面判断数据包的方向
	// macAddr, err := findMacAddrByIp(findDeviceIpv4(device))
	// if err != nil {
	// 	panic(err)
	// }

	// a := CreateNetInfo()
	// b := a["en0"]
	//
	// fmt.Printf("Chosen device's IPv4: %s\n", b.Addrs)
	// fmt.Printf("Chosen device's MAC: %s\n", b.Hardwareaddr)
	// go aaaa(b.Hardwareaddr)

	// fmt.Println("========",b.al.Nc.GetName(),b.Hardwareaddr)

	for {

		fmt.Println(GetCpuPercent())

		time.Sleep(time.Second * 1)
	}

}

// 获取网卡的IPv4地址
func findDeviceIpv4(device pcap.Interface) string {
	for _, addr := range device.Addresses {
		if ipv4 := addr.IP.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}
	panic("device has no IPv4")
}

// 根据网卡的IPv4地址获取MAC地址
// 有此方法是因为gopacket内部未封装获取MAC地址的方法，所以这里通过找到IPv4地址相同的网卡来寻找MAC地址
func findMacAddrByIp(ip string) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(interfaces)
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			panic(err)
		}

		for _, addr := range addrs {
			if a, ok := addr.(*net.IPNet); ok {
				if ip == a.IP.String() {
					return i.HardwareAddr.String(), nil
				}
			}
		}
	}
	return "", errors.New(fmt.Sprintf("no device has given ip: %s", ip))
}
