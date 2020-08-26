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
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

type SavingTotal struct {
	title           *string
	count           *int
	percent         *float64
	totalbackcharge *float64
}

type SavingTotalDetail struct {
	backcharge    *float64
	savingpercent *float64
}

const (
	// 普通分组组数
	GroupNum = 20
	// 每组间隔
	GroupGap = 5
)

// 创建一个指定title的SavingTotal
func CreatSaveingTotal(title string) *SavingTotal {
	count := 0
	percent := 0.00
	totalbackcharge := 0.00
	return &SavingTotal{
		title:           &title,
		count:           &count,
		percent:         &percent,
		totalbackcharge: &totalbackcharge,
	}
}

// 转换分组数据
func DataGrouping(DataSource []*SavingTotalDetail) []*SavingTotal {
	// 输入数据
	output := DataSource
	if len(output) == 0 {
		return make([]*SavingTotal, 0)
	}
	// 初始化处理结果数组
	out := make([]*SavingTotal, 0)
	zero := CreatSaveingTotal("0")
	hundred := CreatSaveingTotal("100")
	none := CreatSaveingTotal("无法计算")
	for i := 0; i < GroupNum; i++ {
		// 格式化title
		title := fmt.Sprintf("%d--%d", i*GroupGap, (i+1)*GroupGap)
		out = append(out, CreatSaveingTotal(title))
	}
	var cur *SavingTotal
	// 遍历输入数据数组
	for _, detail := range output {
		// 优先处理掉值是nil的,这样后面计算分组不会报错
		if detail == nil {
			*none.count = *none.count + 1
			continue
		} else {
			// 做运算计算出当前的数值的百分比是属于哪一个分段的
			grouping := math.Floor(*detail.savingpercent / GroupGap)
			outindex := int(grouping)
			// 按照计算分组index的结果来确定当前处理的分组对象
			if grouping == 0.00 { // 0%分组
				cur = zero
			} else if grouping == 100 { // 100%分组
				cur = hundred
			} else {
				cur = out[outindex] // 普通分组
			}
		}
		// 计算count
		*cur.count = *cur.count + 1
		// 累计金额
		*cur.totalbackcharge = *cur.totalbackcharge + *detail.backcharge
	}

	// 处理普通分组, 单独统计占比这样节约算力
	for _, v := range out {
		percent := float64(*v.count) / float64(len(output))
		v.percent = &percent
	}
	// 处理0%分组
	zeroPercent := float64(*zero.count) / float64(len(output))
	zero.percent = &zeroPercent

	// 处理100%分组
	hundredPercent := float64(*hundred.count) / float64(len(output))
	hundred.percent = &hundredPercent

	// 处理无法计算分组
	nonePercent := float64(*none.count) / float64(len(output))
	none.percent = &nonePercent

	// 拼合数据
	outs := make([]*SavingTotal, 0)
	outs = append(outs, zero)
	outs = append(outs, out...)
	outs = append(outs, hundred)
	outs = append(outs, none)

	return outs
}

func CreateSavingTotalDetail() *SavingTotalDetail {
	backcharge := rand.Float64() * 100000
	savingpercent := rand.Float64() * 100
	return &SavingTotalDetail{
		backcharge:    &backcharge,
		savingpercent: &savingpercent,
	}

}

func TestEntry(t *testing.T) {

	a := SavingTotalDetail{} // 创建一个SavingTotalDetail对象 a
	b := &a                  // 创建 b指针 指向a
	c := &b                  // 创建c指针指向 b
	d := *b                  // 创建 d 获取b指针指向的对象的值
	e := *c                  // 创建 e 获取c指针指向的对象的值

	fmt.Printf("a:%v  \nb:%v   \nc:%v   \nd:%v   \ne:%v   \n", a, b, c, d, e)

	// return

	// fmt.Println(os.Args)
	// cmd := exec.Command(os.Args[0])
	// cmd.Start()
	// time.Sleep(time.Second * 5)
	// fmt.Println(cmd.Process.Pid)
	// syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)

}

func TestAAA(t *testing.T) {
	a := make(map[string][9]map[string]interface{})

	b := [...]map[string]interface{}{
		{
			"unitname":     "一单元",
			"room":         "2-1-101",
			"roomnum":      "01",
			"floor":        1,
			"type":         -5,
			"realunitflow": 0,
		},
		{
			"unitname":     "一单元",
			"room":         "2-1-102",
			"roomnum":      "02",
			"floor":        1,
			"type":         -5,
			"realunitflow": 0,
		},
		{
			"unitname":     "一单元",
			"room":         "2-1-103",
			"roomnum":      "03",
			"floor":        1,
			"type":         -5,
			"realunitflow": 0,
		},
		{
			"unitname":     "一单元",
			"room":         "2-1-104",
			"roomnum":      "04",
			"floor":        1,
			"type":         -5,
			"realunitflow": 0,
		},
		{
			"unitname":     "一单元",
			"room":         "2-1-201",
			"roomnum":      "01",
			"floor":        2,
			"type":         5,
			"realunitflow": 959662,
		},
		{
			"unitname":     "一单元",
			"room":         "2-1-202",
			"roomnum":      "02",
			"floor":        2,
			"type":         -5,
			"realunitflow": 0,
		},
		{
			"unitname":     "一单元",
			"room":         "2-1-203",
			"roomnum":      "03",
			"floor":        2,
			"type":         -5,
			"realunitflow": 0,
		},
		{
			"unitname":     "一单元",
			"room":         "2-1-204",
			"roomnum":      "04",
			"floor":        2,
			"type":         -5,
			"realunitflow": 0,
		},
		{
			"unitname":     "一单元",
			"room":         "2-1-301",
			"roomnum":      "01",
			"floor":        3,
			"type":         -5,
			"realunitflow": 0,
		}}

	a["floor"] = b

	c := make([]map[string]interface{}, 0)
	for _, v := range a["floor"] {
		if v["floor"].(int) > len(c) {
			tmp := make(map[string]interface{})
			tmp["floor"] = v["floor"]
			tmp["temp"] = make([]map[string]interface{}, 0)
			tmp["temp"] = append(tmp["temp"].([]map[string]interface{}), v)
			c = append(c, tmp)
		} else {
			tmp := c[v["floor"].(int)-1]
			tmp["temp"] = append(tmp["temp"].([]map[string]interface{}), v)
		}
	}
	bbb, _ := json.Marshal(c)
	fmt.Println(string(bbb))
}

type PPP struct {
	AAA *string
	BBB *int
	CCC *float64
}

func TestInit(t *testing.T) {
	tt := "123123"
	cc := 1
	u := PPP{
		AAA: &tt,
		BBB: &cc,
		CCC: nil,
	}

	form := make([]interface{}, 0)
	form = append(form, u)
	form = append(form, u)
	form = append(form, u)
	fmt.Println(reflectStruct2StringList(form))
}

// 遍历一个包含struct的数组,将该数组转换成字符串
// PS1: struct必须全是Public
// PS2: 数组是值数组不能是指针数组
func reflectStruct2StringList(form []interface{}) [][]interface{} {
	list := make([][]interface{}, 0)
	for _, v := range form {
		tmpStringList := make([]interface{}, 0)
		fromType := reflect.TypeOf(v)
		fromValue := reflect.ValueOf(v)
		fvNumField := fromType.NumField()
		for i := 0; i < fvNumField; i++ {
			if fromValue.Field(i).IsNil() {
				tmpStringList = append(tmpStringList, "-")
			} else {

				vt := fromValue.Field(i).Type().String()
				var vv interface{}
				switch vt {
				case "float64":
				case "float32":
					vv = fromValue.Field(i).Elem().Float()
				case "int":
				case "int8":
				case "int16":
				case "int32":
				case "int64":
					vv = fromValue.Field(i).Elem().Int()
				case "string":
					vv = fromValue.Field(i).Elem().String()
				}

				tmpStringList = append(tmpStringList, vv)
			}
		}
		list = append(list, tmpStringList)
	}
	return list
}

func TestStopaaa(t *testing.T) {
	fmt.Println(time.Now().UnixNano())
}
