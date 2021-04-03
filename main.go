package main

import (
	"log"
	"mgd-check/core"
	"mgd-check/task"
	"mgd-check/web"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	core.InitConfig()
	task.InitTask()
	web.InitWeb()
}