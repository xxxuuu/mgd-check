package core

import "sync"

var db = sync.Map{}

// 签到信息
type CheckInfo struct {
	// 手机（登录账号）
	Phone string
	// 密码
	Password string
	// 打卡备注
	Description string
	// 国家
	Country string
	// 省份
	Province string
	// 城市
	City string
	// 详细地址
	Address string
	// 纬度
	Latitude string
	// 经度
	Longitude string
}

// 录入信息
func Register(c CheckInfo) {
	db.Store(c.Phone, c)
}

// 遍历所有录入信息
func RangeAllRegisterInfo(handle func(key, value interface{}) bool) {
	db.Range(handle)
}