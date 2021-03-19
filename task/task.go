package task

import (
	"github.com/robfig/cron"
	"log"
	"mgd-check/core"
)

// 上班打卡任务
func startWorkCheck() {
	go core.RangeAllRegisterInfo(func(key, value interface{}) bool {
		context := &core.MgdContext{
			Info: value.(core.CheckInfo),
		}
		err := context.Login()
		if err != nil {
			log.Println(err)
		}
		err = context.StartWork()
		if err != nil {
			log.Println(err)
		}
		return true
	})
}

// 下班打卡任务
func endWorkCheck() {
	go core.RangeAllRegisterInfo(func(key, value interface{}) bool {
		context := &core.MgdContext{
			Info: value.(core.CheckInfo),
		}
		err := context.Login()
		if err != nil {
			log.Println(err)
		}
		err = context.EndWork()
		if err != nil {
			log.Println(err)
		}
		return true
	})
}

func InitTask() {
	// 上下班打卡定时任务
	c := cron.New()
	_ = c.AddFunc("0 0 7 * * ?", startWorkCheck)
	_ = c.AddFunc("0 0 17 * * ?", endWorkCheck)
	c.Start()
}

