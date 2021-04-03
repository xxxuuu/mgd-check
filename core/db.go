package core

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
)

var memDb MemDb
var mongoDb MongoDb

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
	// 通知邮箱
	NoticeEmail string
}

type Db interface {
	// 录入信息
	Register(c CheckInfo)
	// 遍历所有录入信息
	RangeAllRegisterInfo(handle func(key, value interface{}) bool)
}

type MemDb struct {
	db sync.Map
	init sync.Once
}

func (m *MemDb) initMemDb() {
	m.db = sync.Map{}
}

func (m *MemDb) Register(c CheckInfo) {
	m.init.Do(m.initMemDb)
	m.db.Store(c.Phone, c)
}

func (m *MemDb) RangeAllRegisterInfo(handle func(key, value interface{}) bool) {
	m.init.Do(m.initMemDb)
	m.db.Range(handle)
}

type MongoDb struct {
	init sync.Once
	client *mongo.Client
}

func (r *MongoDb) initClient() {
	host := fmt.Sprintf("mongodb://%s", viper.GetString("db.host"))
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(host))
	if err != nil {
		log.Println(err)
	}
	r.client = client
}

func (r *MongoDb) Register(c CheckInfo) {
	r.init.Do(r.initClient)
	collection := r.client.Database("mgd").Collection("checkinfo")

	where := bson.D{{"phone", c.Phone}}
	var result CheckInfo
	err := collection.FindOne(context.Background(), where).Decode(&result)
	if err == nil {
		_, _ = collection.DeleteOne(context.Background(), where)
	}
	_, _ = collection.InsertOne(context.Background(), c)
}

func (r *MongoDb) RangeAllRegisterInfo(handle func(key, value interface{}) bool) {
	r.init.Do(r.initClient)
	collection := r.client.Database("mgd").Collection("checkinfo")

	cur, err := collection.Find(context.Background(), bson.D{{}}, options.Find())
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.Background()) {
		var elem CheckInfo
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
		}
		handle(elem.Phone, elem)
	}

	if err := cur.Err(); err != nil {
		log.Println(err)
	}

	_ = cur.Close(context.Background())
}

func GetDb() Db {
	config := viper.Get("db.type")
	if config == nil {
		return &memDb
	}
	switch config {
	case "":
	case "mem":
		return &memDb
	case "mongodb":
		return &mongoDb
	}
	panic("db config error")
}
