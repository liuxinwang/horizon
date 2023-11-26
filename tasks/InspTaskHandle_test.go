package tasks

import (
	"fmt"
	"horizon/config"
	"horizon/model"
	"horizon/utils"
	"testing"
)

func TestInspTaskRunning(t *testing.T) {
	InspTaskRunning("mysql-xx501o9101")
}

func TestMain(m *testing.M) {
	fmt.Println("begin")
	help := utils.HelpInit()
	config.InitConfig(help.ConfigFile)
	model.InitDb()
	m.Run()
	fmt.Println("end")
}
