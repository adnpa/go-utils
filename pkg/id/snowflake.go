package id

import (
	"fmt"
	"log"
	"time"
)

// 类Snowflake算法生成的全局唯一id

var (
	sonyFlake     *sonyflake.Sonyflake // 实例
	sonyMachineID uint16               // 机器ID
)

func getMachineID() (uint16, error) { // 返回全局定义的机器ID
	return sonyMachineID, nil
}

func init() {
	var err error
	var sonyMachineID = conf.Cfg.MachineId
	if sonyMachineID == 0 {
		panic(fmt.Errorf("config not found: snowflake machine id"))
	}
	t, _ := time.Parse("2006-01-02", "2022-02-09") // 初始化一个开始的时间
	settings := sonyflake.Settings{                // 生成全局配置
		StartTime: t,
		MachineID: getMachineID, // 指定机器ID
	}
	sonyFlake, err = sonyflake.New(settings) // 用配置生成snowflake节点
	if err != nil {
		panic(err)
	}
}

// GetID 返回生成的id值
func GetID() (id uint64, err error) { // 拿到sonyFlake节点生成id值
	if sonyFlake == nil {
		err = fmt.Errorf("snoy flake not inited")
		log.Println("not init")
		return
	}
	id, err = sonyFlake.NextID()
	return
}
