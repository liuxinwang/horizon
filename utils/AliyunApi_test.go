// This file is auto-generated, don't edit it. Thanks.
package utils

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"horizon/config"
	"horizon/model"
	"testing"
)

func TestListClusters(t *testing.T) {
	clusters, err := ListClusters(1, 10, "")
	var totalCount int64
	var configCenterInstances []model.ConfigCenterInstance
	if err != nil {
		return
	}
	body := clusters["body"]
	if b, ok := body.(map[string]interface{}); ok {
		instanceData := b["Data"]
		if tc, isInt := b["TotalCount"].(json.Number); isInt {
			totalCount, _ = tc.Int64()
		}
		if d, ok2 := instanceData.([]interface{}); ok2 {
			for _, cluster := range d {
				var configCenterInstance model.ConfigCenterInstance
				err := mapstructure.Decode(cluster, &configCenterInstance)
				if err != nil {
					return
				}
				configCenterInstances = append(configCenterInstances, configCenterInstance)
			}
		}
	}
	fmt.Println("totalCount", totalCount)
	fmt.Println("configCenterInstances", configCenterInstances)
}

func TestListNacosConfigs(t *testing.T) {
	configs, err := ListNacosConfigs(1, 10, "", "", "")
	var totalCount int64
	var configCenterInstanceConfigs []model.ConfigCenterInstanceConfig
	if err != nil {
		return
	}
	body := configs["body"]
	if b, ok := body.(map[string]interface{}); ok {
		instanceConfigData := b["Configurations"]
		if tc, isInt := b["TotalCount"].(json.Number); isInt {
			totalCount, _ = tc.Int64()
		}
		if d, ok2 := instanceConfigData.([]interface{}); ok2 {
			for _, cluster := range d {
				var configCenterInstanceConfig model.ConfigCenterInstanceConfig
				err := mapstructure.Decode(cluster, &configCenterInstanceConfig)
				if err != nil {
					return
				}
				configCenterInstanceConfigs = append(configCenterInstanceConfigs, configCenterInstanceConfig)
			}
		}
	}
	fmt.Println("totalCount", totalCount)
	fmt.Println("configCenterInstanceConfigs", configCenterInstanceConfigs)
}

func TestGetNacosConfig(t *testing.T) {
	configs, err := GetNacosConfig("", "", "", "")
	var configCenterInstanceConfig model.ConfigCenterInstanceConfig
	if err != nil {
		return
	}
	body := configs["body"]
	if b, ok := body.(map[string]interface{}); ok {
		instanceConfigData := b["Configuration"]
		err := mapstructure.Decode(instanceConfigData, &configCenterInstanceConfig)
		if err != nil {
			return
		}
	}
	fmt.Println("configCenterInstanceConfig", configCenterInstanceConfig)
}

func TestDesensitizedYamlConfig(t *testing.T) {
	content := ""
	// var data map[string]interface{}
	var dataNode yaml.Node
	err := yaml.Unmarshal([]byte(content), &dataNode)
	if err != nil {
		return
	}

	flag := 0
	iterNode(&dataNode, &flag)
	marshal, err := yaml.Marshal(&dataNode)
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Println(marshal)
}

func iterNode(node *yaml.Node, flag *int) {
	if node.Content != nil && len(node.Content) > 0 {
		for _, n := range node.Content {
			iterNode(n, flag)
		}
	} else {
		if node.Value == "password" || node.Value == "appSecret" {
			*flag = 1
			return
		}
		if *flag == 1 {
			node.Value = "******"
			*flag = 0
		}
	}
}

func TestMain(m *testing.M) {
	fmt.Println("begin")
	help := HelpInit()
	config.InitConfig(help.ConfigFile)
	model.InitDb()
	m.Run()
	fmt.Println("end")
}
