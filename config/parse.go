package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type General struct {
	SecretKey   string
	Environment string
	Port        int32
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

type Config struct {
	General    General
	Mysql      Mysql
	Prometheus Prometheus
}

var Conf Config

func InitConfig() {
	if _, err := toml.DecodeFile("./conf.toml", &Conf); err != nil {
		panic("ERROR occurred:" + err.Error())
	}
	fmt.Printf("%s (%s)\n", Conf.Mysql.User, Conf.Mysql.Password)
}
