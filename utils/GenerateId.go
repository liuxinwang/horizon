package utils

import (
	"fmt"
	"horizon/model"
	"math/rand"
	"strings"
	"time"
)

// GenerateId 生成实例ID
func GenerateId(obj interface{}) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(letterRunes))]
	}
	var objId = ""
	switch objValue := obj.(type) {
	case *model.Instance:
		objId = fmt.Sprintf("%s-%s", strings.ToLower(string(objValue.Type)), string(b))
	case *model.Project:
		objId = "proj-" + string(b)
	}
	return objId

}
