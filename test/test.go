package main

import (
	"fmt"
	"github.com/shirou/gopsutil/host"
	"time"
)

func main(){
	systemtime,_ := host.BootTime()
	t := time.Unix(int64(systemtime), 0)
	fmt.Println(t.Local().Format("2006-01-02 15:04:05"))

}
