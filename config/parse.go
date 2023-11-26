package config

import (
	"github.com/BurntSushi/toml"
)

type General struct {
	SecretKey   string
	Environment string
	Port        int32
	HomeAddress string
}

type Mysql struct {
	Host     string
	Port     int32
	User     string
	Password string
	Db       string
}

type Prometheus struct {
	Host string
	Port int32
}

type DingWebhook struct {
	Webhook string
}

type GoInception struct {
	Host     string
	Port     int32
	User     string
	Password string
}

type Log struct {
	Name  string
	Level string
}

type Config struct {
	General     General
	Mysql       Mysql
	Prometheus  Prometheus
	DingWebhook DingWebhook
	GoInception GoInception
	Log         Log
}

var Conf Config

func InitConfig(configFile string) {
	if _, err := toml.DecodeFile(configFile, &Conf); err != nil {
		panic("ERROR occurred:" + err.Error())
	}
}
