package task

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"log"
	"mgd-check/core"
)

type Result struct {
	code int
	msg string
}

const (
	SUCCESS = iota
	FAIL
)

const (
	START_WORK = iota
	END_WORK
)

func SendEmail(to string, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", viper.GetString("email.username"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "mgd-check 打卡通知")
	m.SetBody("text/html", body)

	host := viper.GetString("email.host")
	port := viper.GetInt("email.port")
	username := viper.GetString("email.username")
	password := viper.GetString("email.password")
	d := gomail.NewDialer(host, port, username, password)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("给 %s 的邮件发送失败：%v\n", to, err)
		return
	}
	log.Printf("给 %s 的邮件发送成功", to)
}

// 打卡结果通知
func notice(info core.CheckInfo, result Result) {
	var msg string
	if result.code == SUCCESS {
		msg = fmt.Sprintf("账号 %s 打卡成功", info.Phone)
	} else {
		msg = fmt.Sprintf("账号 %s 打卡失败:%v", info.Phone, result.msg)
	}
	log.Println(msg)

	if viper.Get("email.enable") != nil && info.NoticeEmail != "" {
		go SendEmail(info.NoticeEmail, msg)
	}
}

func workCheck(checkType int) {
	go core.GetDb().RangeAllRegisterInfo(func(key, value interface{}) bool {
		context := &core.MgdContext{
			Info: value.(core.CheckInfo),
		}
		err := context.Login()
		if err != nil {
			notice(context.Info, Result{code: FAIL, msg: fmt.Sprintf("%v", err)})
			return true
		}

		if checkType == START_WORK {
			err = context.StartWork()
		} else if checkType == END_WORK {
			err = context.EndWork()
		}
		if err != nil {
			notice(context.Info, Result{code: FAIL, msg: fmt.Sprintf("%v", err)})
		}

		notice(context.Info, Result{code: SUCCESS, msg: "ok"})
		return true
	})
}


func InitTask() {
	// 上下班打卡定时任务
	c := cron.New()
	_ = c.AddFunc("0 0 7 * * ?", func() {
		workCheck(START_WORK)
	})
	_ = c.AddFunc("0 0 17 * * ?", func() {
		workCheck(END_WORK)
	})
	c.Start()
}

