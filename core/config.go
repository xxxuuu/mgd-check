package core

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("配置文件不存在,将使用环境变量")
	}

	log.Println("配置读取完毕")
	log.Printf("DB Type:%v\n", viper.Get("db.type"))
	log.Printf("email host:%v\n", viper.Get("email.host"))
}