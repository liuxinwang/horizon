package tasks

import (
	"fmt"
	"horizon/config"
	"horizon/model"
	"testing"
)

func TestInspTaskRunning(t *testing.T) {
	InspTaskRunning("mysql-xx501o9101")
}

func TestMain(m *testing.M) {
	fmt.Println("begin")
	config.InitConfig()
	model.InitDb()
	m.Run()
	fmt.Println("end")
}
