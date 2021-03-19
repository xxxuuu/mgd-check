package main

import (
	"log"
	"mgd-check/task"
	"mgd-check/web"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	task.InitTask()
	web.InitWeb()
}