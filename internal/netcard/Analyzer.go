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
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Analyzer struct {
	Nc                 *NetCard  // 网卡
	Running            bool      // 是否在运行
	stop               chan bool // 停止信号
	allCanceled        chan bool // 所有goroutines是否中止
	downStreamDataSize int       // 单位时间内下行的总字节数
	upStreamDataSize   int       // 单位时间内上行的总字节数
	upSpeed            float32   // 转化后的上行速度
	downSpeed          float32   // 转化后的下行速度
}

func (anl *Analyzer) Init() {
	anl.stop = make(chan bool)
	anl.allCanceled = make(chan bool)
	anl.Running = false
}

func (anl *Analyzer) Capture() {
	anl.Running = true
	handler, err := pcap.OpenLive(anl.Nc.name, 1024, true, 30*time.Second)
	if err != nil {
		panic(err)
	}
	defer handler.Close()
	ctx, cancel := context.WithCancel(context.Background())
	// 开启子线程，每一秒计算一次该秒内的数据包大小平均值，并将下载、上传总量置零
	go anl.monitor(ctx)

	// 开始抓包
	packetSource := gopacket.NewPacketSource(handler, handler.LinkType())
	// 这种方式从channel中读数据很有意思
	for packet := range packetSource.Packets() {
		// 只获取以太网帧
		ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
		if ethernetLayer != nil {
			ethernet := ethernetLayer.(*layers.Ethernet)
			// 如果封包的目的MAC是本机则表示是下行的数据包，否则为上行
			if ethernet.DstMAC.String() == anl.Nc.mac {
				anl.downStreamDataSize += len(packet.Data()) // 统计下行封包总大小
			} else {
				anl.upStreamDataSize += len(packet.Data()) // 统计上行封包总大小
			}
		}
		select {
		case <-anl.stop:
			cancel()
			return
		default:
			continue
		}
	}
}

// 每一秒计算一次该秒内的数据包大小平均值，并将下载、上传总量置零
func (anl *Analyzer) monitor(ctx context.Context) {
	for {
		os.Stdout.WriteString(fmt.Sprintf("\rDown:%.2fkb/s \t Up:%.2fkb/s", float32(anl.downStreamDataSize)/1024/1, float32(anl.upStreamDataSize)/1024/1))
		anl.downSpeed = float32(anl.downStreamDataSize) / 1024 / 1
		anl.upSpeed = float32(anl.upStreamDataSize) / 1024 / 1
		anl.downStreamDataSize = 0
		anl.upStreamDataSize = 0

		time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			anl.Running = false
			anl.allCanceled <- true
			return
		default:
			continue
		}
	}
}

func (anl *Analyzer) GetUpSpeed() float32 {
	return anl.upSpeed
}

func (anl *Analyzer) GetDownSpeed() float32 {
	return anl.downSpeed
}

func (anl *Analyzer) Stop() {
	anl.stop <- true
	<-anl.allCanceled
}
