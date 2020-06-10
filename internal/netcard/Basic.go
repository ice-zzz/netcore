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
package netcard

import (
	"errors"
	"fmt"
	"net"

	"github.com/google/gopacket/pcap"
)

// 获取网卡的IPv4地址
func findDeviceIpv4(device pcap.Interface) (string, error) {
	for _, addr := range device.Addresses {
		if ipv4 := addr.IP.To4(); ipv4 != nil {
			return ipv4.String(), nil
		}
	}
	return "", errors.New("no IPv4 Found")
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
					// fmt.Println("found one")
					// fmt.Println(i.HardwareAddr.String())
					return i.HardwareAddr.String(), nil
				}
			}
		}
	}
	return "", errors.New(fmt.Sprintf("no device has given ip: %s", ip))
}

// 得到所有具有IPv4地址的网卡
func GetNetCardsWithIPv4Addr() ([]NetCard, error) {
	ncs := make([]NetCard, 0, 4)
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, errors.New("find device error")
	}
	for _, d := range devices {
		var nc NetCard
		ip, e := findDeviceIpv4(d)
		if e != nil {
			continue
		}
		mac, e := findMacAddrByIp(ip)
		if e != nil {
			continue
		}
		nc.ipv4 = ip
		nc.mac = mac
		nc.name = d.Name
		ncs = append(ncs, nc)
	}
	return ncs, nil
}
