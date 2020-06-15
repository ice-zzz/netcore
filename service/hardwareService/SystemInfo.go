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
	"time"

	"github.com/ice-zzz/netcore/internal/netcard"
	"github.com/ice-zzz/netcore/service"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type SYSTEM struct {
	CPU      []*CpuInfo          `json:"cpu" toml:"cpu"`
	DISK     *DiskInfo           `json:"disk" toml:"disk"`
	NET      map[string]*NetInfo `json:"net" toml:"net"`
	HOST     *HostInfo           `json:"host" toml:"host"`
	MEM      *MemInfo            `json:"mem" toml:"mem"`
	exitChan chan struct{}       `json:"-" toml:"-"`
	service.Entity
}

func (s *SYSTEM) Start() {
	s.exitChan = make(chan struct{})
	s.CPU = GetCpuInfo()
	s.MEM = GetMemInfo()
	s.HOST = GetHostInfo()
	s.NET = CreateNetInfo()
	s.DISK = GetDiskInfo()

	for {
		select {
		case <-s.exitChan:
			return
		default:
			s.DISK = GetDiskInfo()
			s.MEM = GetMemInfo()
			time.Sleep(time.Second * 3)
		}

	}
}

func (s *SYSTEM) Stop() {
	s.exitChan <- struct{}{}
}

func (ni *NetInfo) GetNetSpeed() (upspeed float32, downspeed float32) {

	return ni.al.GetUpSpeed(), ni.al.GetDownSpeed()
}

type CpuInfo struct {
	ModelName string  `json:"model_name" toml:"model_name"`
	Cores     int32   `json:"cores" toml:"cores"`
	Mhz       float64 `json:"mhz" toml:"mhz"`
	CacheSize int32   `json:"cache_size" toml:"cache_size"`
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
	Total       int     `json:"total" toml:"total"`
	Used        int     `json:"used" toml:"used"`
	Free        int     `json:"free" toml:"free"`
	UsedPercent float64 `json:"used_percent" toml:"used_percent"`
}

func GetMemInfo() *MemInfo {
	v, _ := mem.VirtualMemory()
	mi := &MemInfo{}
	mi.Total = int(v.Total) / 1024 / 1024
	mi.Used = int(v.Used) / 1024 / 1024
	mi.Free = int(v.Free) / 1024 / 1024
	mi.UsedPercent = v.UsedPercent

	return mi
}

type NetInfo struct {
	Name               string `json:"name" toml:"name"`
	Hardwareaddr       string `json:"hardwareaddr" toml:"hardwareaddr"`
	Addrs              string `json:"addrs" toml:"addrs"`
	downStreamDataSize int
	upStreamDataSize   int
	downSpeed          float32
	upSpeed            float32
	al                 *netcard.Analyzer
}

func CreateNetInfo() map[string]*NetInfo {
	nets := make(map[string]*NetInfo)
	v, _ := netcard.GetNetCardsWithIPv4Addr()

	for _, vv := range v {
		if len(vv.GetMacAddr()) <= 0 {
			continue
		}
		vvv := vv
		analyzer := &netcard.Analyzer{}
		analyzer.Init()
		analyzer.Nc = &vvv
		nets[vv.GetName()] = &NetInfo{
			Name:         vv.GetName(),
			Hardwareaddr: vv.GetMacAddr(),
			Addrs:        vv.GetIPv4Addr(),
			al:           analyzer,
		}
		go analyzer.Capture()
	}

	return nets
}

type HostInfo struct {
	Hostname        string `json:"hostname" toml:"hostname"`
	OS              string `json:"os" toml:"os"`
	Platform        string `json:"platform" toml:"platform"`
	PlatformVersion string `json:"platform_version" toml:"platform_version"`
	KernelVersion   string `json:"kernel_version" toml:"kernel_version"`
	KernelArch      string `json:"kernel_arch" toml:"kernel_arch"`
	Hostid          string `json:"hostid" toml:"hostid"`
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
	Fstype      string  `json:"fstype" toml:"fstype"`
	Total       uint64  `json:"total" toml:"total"`
	Free        uint64  `json:"free" toml:"free"`
	Used        uint64  `json:"used" toml:"used"`
	UsedPercent float64 `json:"used_percent" toml:"used_percent"`
}

func GetDiskInfo() *DiskInfo {
	v, _ := disk.Usage("/")
	return &DiskInfo{
		Fstype:      v.Fstype,
		Total:       v.Total / 1024 / 1024,
		Used:        v.Used / 1024 / 1024,
		Free:        v.Free / 1024 / 1024,
		UsedPercent: v.UsedPercent,
	}
}
