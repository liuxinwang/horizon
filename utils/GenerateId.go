package utils

import (
	"horizon/model"
	"math/rand"
	"time"
)

// GenerateId 生成实例ID
func GenerateId(obj interface{}) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(letterRunes))]
	}
	var obj_id = ""
	switch obj.(type) {
	case *model.Instance:
		obj_id = "mysql-" + string(b)
	case *model.Project:
		obj_id = "proj-" + string(b)
	}
	return obj_id

}
