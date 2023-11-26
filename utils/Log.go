package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"horizon/config"
	"os"
)

func LogInit() {
	var file *os.File
	var err error
	if config.Conf.General.Environment == "dev" {
		file = os.Stdout
	} else {
		file, err = os.OpenFile(config.Conf.Log.Name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		logrus.SetOutput(file)
	}
	if err != nil {
		logrus.Fatalf("log init failed, err: %v", err.Error())
	}
	logrus.SetOutput(file)

	level, err := logrus.ParseLevel(config.Conf.Log.Level)
	if err != nil {
		logrus.Fatalf("log init failed, err: %v", err.Error())
	}
	logrus.SetLevel(level)

	// gin release mode
	gin.SetMode(gin.ReleaseMode)
	// gin logo handle
	gin.DisableConsoleColor() // 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	gin.DefaultWriter = file
}
