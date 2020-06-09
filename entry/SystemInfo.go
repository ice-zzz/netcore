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
	"github.com/ice-zzz/netcore/internal/netcard"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type SYSTEM struct {
	CPU  []*CpuInfo `toml:"cpu"`
	DISK *DiskInfo  `toml:"disk"`
	NET  []*NetInfo `toml:"net"`
	HOST *HostInfo  `toml:"host"`
	MEM  *MemInfo   `toml:"mem"`
}

type CpuInfo struct {
	ModelName string  `toml:"model_name"`
	Cores     int32   `toml:"cores"`
	Mhz       float64 `toml:"mhz"`
	CacheSize int32   `toml:"cache_size"`
}

func GetCpuInfo() []*CpuInfo {
	cpus, _ := cpu.Info()
	cinfos := make([]*CpuInfo, len(cpus))
	for k, v := range cpus {
		cinfos[k] = &CpuInfo{
			ModelName: v.ModelName,
			Cores:     v.Cores,
			Mhz:       v.Mhz,
			CacheSize: v.CacheSize,
		}
	}
	return cinfos
}

type MemInfo struct {
	Total int `toml:"total"`
}

func GetMemInfo() *MemInfo {
	v, _ := mem.VirtualMemory()
	return &MemInfo{
		Total: int(v.Total) / 1024 / 1024,
	}
}

type NetInfo struct {
	Name         string `toml:"name"`
	Hardwareaddr string `toml:"hardwareaddr"`
	Addrs        string `toml:"addrs"`
}

func GetNetInfo() []*NetInfo {
	nets := make([]*NetInfo, 0)
	v, _ := netcard.GetNetCardsWithIPv4Addr()

	for _, vv := range v {
		nets = append(nets, &NetInfo{
			Name:         vv.GetName(),
			Hardwareaddr: vv.GetMacAddr(),
			Addrs:        vv.GetIPv4Addr(),
		})
	}

	return nets
}

type HostInfo struct {
	Hostname        string `toml:"hostname"`
	OS              string `toml:"os"`
	Platform        string `toml:"platform"`
	PlatformVersion string `toml:"platform_version"`
	KernelVersion   string `toml:"kernel_version"`
	KernelArch      string `toml:"kernel_arch"`
	Hostid          string `toml:"hostid"`
}

func GetHostInfo() *HostInfo {
	v, _ := host.Info()
	return &HostInfo{
		Hostname:        v.Hostname,
		OS:              v.OS,
		Platform:        v.Platform,
		PlatformVersion: v.PlatformVersion,
		KernelVersion:   v.KernelVersion,
		KernelArch:      v.KernelArch,
		Hostid:          v.HostID,
	}
}

type DiskInfo struct {
	Fstype      string  `toml:"fstype"`
	Total       uint64  `toml:"total"`
	Free        uint64  `toml:"free"`
	Used        uint64  `toml:"used"`
	UsedPercent float64 `toml:"used_percent"`
}

func GetDiskInfo() *DiskInfo {
	v, _ := disk.Usage("/")
	return &DiskInfo{
		Fstype:      v.Fstype,
		Total:       v.Total,
		Free:        v.Free,
		Used:        v.Used,
		UsedPercent: v.UsedPercent,
	}
}
