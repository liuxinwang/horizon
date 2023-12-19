package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"horizon/model"
	"horizon/utils"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// ConfigCenterInstanceSelectByList 查询实例列表
func ConfigCenterInstanceSelectByList(c *gin.Context) {
	// 变量初始化
	pageNo, _ := strconv.Atoi(c.Query("pageNo"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	clusterAliasName := c.Query("clusterAliasName")
	var configCenterInstances []model.ConfigCenterInstance
	var totalCount int64
	data := gin.H{"totalCount": 0, "data": &[]model.ConfigCenterInstance{}, "pageNo": pageNo, "pageSize": pageSize, "totalPage": 0}

	clusters, err := utils.ListClusters(pageNo, pageSize, clusterAliasName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": data, "err": err.Error()})
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
					c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": data, "err": err.Error()})
					return
				}
				// 查询实例空间
				namespaces, err := configCenterInstanceNamespaceSelectByList(configCenterInstance.InstanceId)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": data, "err": err.Error()})
					return
				}
				configCenterInstance.Namespaces = namespaces
				configCenterInstances = append(configCenterInstances, configCenterInstance)
			}
		}
	}

	// 处理结果集并返回
	totalPage := math.Ceil(float64(totalCount) / float64(pageSize))
	if totalCount > 0 {
		data = gin.H{"totalCount": totalCount, "data": &configCenterInstances, "pageNo": pageNo, "pageSize": pageSize, "totalPage": totalPage}
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": data, "err": ""})
}

// configCenterInstanceNamespaceSelectByList 查询实例命名空间列表
func configCenterInstanceNamespaceSelectByList(instanceId string) ([]model.ConfigCenterNamespace, error) {
	var configCenterNamespaces []model.ConfigCenterNamespace
	namespaces, err := utils.ListEngineNamespaces(instanceId)
	if err != nil {
		return nil, err
	}
	body := namespaces["body"]
	if b, ok := body.(map[string]interface{}); ok {
		namespaceData := b["Data"]
		if d, ok2 := namespaceData.([]interface{}); ok2 {
			for _, namespace := range d {
				var configCenterNamespace model.ConfigCenterNamespace
				err := mapstructure.Decode(namespace, &configCenterNamespace)
				if err != nil {
					return nil, err
				}
				configCenterNamespace.NamespaceId = configCenterNamespace.NamespaceShowName
				configCenterNamespaces = append(configCenterNamespaces, configCenterNamespace)
			}
		}
	}
	return configCenterNamespaces, nil
}

// ConfigCenterInstanceConfigSelectByList 查询实例配置列表
func ConfigCenterInstanceConfigSelectByList(c *gin.Context) {
	// 变量初始化
	pageNo, _ := strconv.Atoi(c.Query("pageNo"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	instanceId := c.Query("instanceId")
	namespaceId := c.Query("namespaceId")
	dataId := c.Query("dataId")
	var configCenterInstanceConfigs []model.ConfigCenterInstanceConfig
	var totalCount int64
	data := gin.H{"totalCount": 0, "data": &[]model.ConfigCenterInstanceConfig{}, "pageNo": pageNo, "pageSize": pageSize, "totalPage": 0}

	configs, err := utils.ListNacosConfigs(pageNo, pageSize, instanceId, namespaceId, dataId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": data, "err": err.Error()})
		return
	}
	body := configs["body"]
	if b, ok := body.(map[string]interface{}); ok {
		instanceConfigData := b["Configurations"]
		if tc, isInt := b["TotalCount"].(json.Number); isInt {
			totalCount, _ = tc.Int64()
		}
		if d, ok2 := instanceConfigData.([]interface{}); ok2 {
			for _, config := range d {
				var configCenterInstanceConfig model.ConfigCenterInstanceConfig
				err := mapstructure.Decode(config, &configCenterInstanceConfig)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": data, "err": err.Error()})
					return
				}
				configCenterInstanceConfig.InstanceId = instanceId
				configCenterInstanceConfig.NamespaceId = namespaceId
				configCenterInstanceConfigs = append(configCenterInstanceConfigs, configCenterInstanceConfig)
			}
		}
	}

	// 处理结果集并返回
	totalPage := math.Ceil(float64(totalCount) / float64(pageSize))
	if totalCount > 0 {
		data = gin.H{"totalCount": totalCount, "data": &configCenterInstanceConfigs, "pageNo": pageNo, "pageSize": pageSize, "totalPage": totalPage}
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": data, "err": ""})
}

// ConfigCenterInstanceConfigSelectByData 查看实例配置详情
func ConfigCenterInstanceConfigSelectByData(c *gin.Context) {
	var configCenterInstanceConfig model.ConfigCenterInstanceConfig
	err := c.ShouldBind(&configCenterInstanceConfig)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	config, err := utils.GetNacosConfig(
		configCenterInstanceConfig.InstanceId,
		configCenterInstanceConfig.NamespaceId,
		configCenterInstanceConfig.DataId,
		configCenterInstanceConfig.Group)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": nil, "err": err.Error()})
		return
	}
	body := config["body"]
	if b, ok := body.(map[string]interface{}); ok {
		instanceConfigData := b["Configuration"]
		err := mapstructure.Decode(instanceConfigData, &configCenterInstanceConfig)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": nil, "err": err.Error()})
			return
		}
	}
	desensitizedContent, err := desensitizedYamlConfig(configCenterInstanceConfig.Content)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": nil, "err": err.Error()})
		return
	}
	configCenterInstanceConfig.Content = strings.ReplaceAll(desensitizedContent, "    ", "  ")
	// 返回结果集
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": &configCenterInstanceConfig, "err": ""})
}

func desensitizedYamlConfig(content string) (string, error) {
	var dataNode yaml.Node
	err := yaml.Unmarshal([]byte(content), &dataNode)
	if err != nil {
		return "", err
	}

	flag := 0
	iterNode(&dataNode, &flag)
	marshal, err := yaml.Marshal(&dataNode)
	if err != nil {
		return "", err
	}
	return string(marshal), err
}

func iterNode(node *yaml.Node, flag *int) {
	if node.Content != nil && len(node.Content) > 0 {
		for _, n := range node.Content {
			iterNode(n, flag)
		}
	} else {
		if node.Value == "password" ||
			node.Value == "appSecret" ||
			strings.Contains(strings.ToLower(node.Value), "secret") {
			*flag = 1
			return
		}
		if *flag == 1 {
			desensitizedNode(node)
			*flag = 0
		}
	}
}

func desensitizedNode(node *yaml.Node) {
	node.Value = "******"
}
