package configure

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	WeaponListCSV = "/CSO2-Server/assert/cstrike/scripts/item_list.csv"
	ExpLevelCSV   = "/CSO2-Server/assert/cstrike/scripts/exp_level.csv"
	UnlockCSV     = "/CSO2-Server/assert/cstrike/scripts/item_unlock.csv"
	BoxCSV        = "/CSO2-Server/assert/cstrike/scripts/supplyList.csv"
	VipCSV        = "/CSO2-Server/assert/cstrike/scripts/vip_info.csv"
)

type ItemData struct {
	ItemID      uint32
	Name        string
	Class       string
	Category    string
	BuyCategory string
}

type UnlockData struct {
	ItemID         uint32
	NextItemID     uint32
	ConditionFlag0 uint32
	Count0         uint32
	ConditionFlag1 uint32
	Count1         uint32
	ConditionFlag2 uint32
	Count2         uint32
}

type BoxData struct {
	BoxID      uint32
	Items      []BoxItem
	TotalValue int
}

type BoxItem struct {
	ItemID uint32
	Value  int
}

var (
	ItemList   = make(map[uint32]ItemData)
	UnlockList = make(map[uint32]UnlockData)
	BoxList    = make(map[uint32]BoxData)
)

func InitCSV(path string) {
	fmt.Println("Reading game data file ...")
	readWeaponList(path)
	readUnlockList(path)
	readBoxList(path)
}

func readWeaponList(path string) {
	//读取武器数据
	filepath := path + WeaponListCSV

	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if err == nil && len(record[1]) > 16 {
			id, err := strconv.Atoi(record[0])
			if err != nil {
				continue
			}
			itemd := ItemData{
				uint32(id),
				record[1][16:],
				record[4],
				record[5],
				record[6],
			}

			ItemList[itemd.ItemID] = itemd
		} else {
			continue
		}
	}

}

func readUnlockList(path string) {
	//读取武器解锁数据
	filepath := path + UnlockCSV

	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if err == nil {
			id, err := strconv.Atoi(record[0])
			if err != nil {
				continue
			}
			nextid, err := strconv.Atoi(record[1])
			if err != nil {
				continue
			}
			flag0, err := strconv.Atoi(record[5])
			if err != nil {
				continue
			}
			count0, err := strconv.Atoi(record[6])
			if err != nil {
				continue
			}
			flag1, err := strconv.Atoi(record[11])
			if err != nil {
				continue
			}
			count1, err := strconv.Atoi(record[12])
			if err != nil {
				continue
			}
			flag2, err := strconv.Atoi(record[17])
			if err != nil {
				continue
			}
			count2, err := strconv.Atoi(record[18])
			if err != nil {
				continue
			}
			itemd := UnlockData{
				uint32(id),
				uint32(nextid),
				uint32(flag0),
				uint32(count0),
				uint32(flag1),
				uint32(count1),
				uint32(flag2),
				uint32(count2),
			}

			UnlockList[itemd.ItemID] = itemd
		} else {
			continue
		}
	}
}

func readBoxList(path string) {
	//读取箱子数据
	filepath := path + BoxCSV

	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if err == nil {
			boxid, err := strconv.Atoi(record[0])
			if err != nil {
				continue
			}
			itemid, err := strconv.Atoi(record[1])
			if err != nil {
				continue
			}
			value, err := strconv.Atoi(record[2])
			if err != nil {
				continue
			}
			//保存当前物品数据
			if value <= 0 {
				fmt.Println("Warning ! illeagal value", value, "for item", itemid, "in box", boxid)
				continue
			}
			item := BoxItem{
				uint32(itemid),
				value,
			}
			if v, ok := BoxList[uint32(boxid)]; ok {
				//如果该box数据已经存在
				v.Items = append(v.Items, item)
				v.TotalValue += value
				BoxList[uint32(boxid)] = v
			} else {
				BoxList[uint32(boxid)] = BoxData{
					uint32(boxid),
					[]BoxItem{item},
					value,
				}
			}
		} else {
			continue
		}
	}
}
