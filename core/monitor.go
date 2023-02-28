package core

import (
	"github.com/shirou/gopsutil/load"
	"ops_agent/conf"
	"ops_agent/pkg"

	"time"
)

var (
	MonitorChan chan pkg.Store_Monitor_Data
)

type Monitor struct {
}

func (m *Monitor) InitMonitorChan() {
	MonitorChan = make(chan pkg.Store_Monitor_Data, 100) // 通道初始化(带缓冲区)
}

func RunMonitor() {
	exec := Monitor{}
	for {
		go exec.MonitorUptime()
		time.Sleep(time.Second * 1)
	}
}

func (m *Monitor) MonitorUptime() {
	data := make(map[string]interface{})
	var list_data []map[string]interface{}

	loads, _ := load.Avg()
	data["load1"] = loads.Load1
	data["load5"] = loads.Load5
	data["load15"] = loads.Load15

	list_data = append(list_data, data)
	commit_data := pkg.Store_Monitor_Data{
		Id:   Asset_ID,
		Type: Cpu,
		Data: list_data,
	}
	conf.Logger.Info("获取负载信息成功 %s", list_data)

	// channel存储数据
	MonitorChan <- commit_data
	time.Sleep(time.Second * 10)
}
