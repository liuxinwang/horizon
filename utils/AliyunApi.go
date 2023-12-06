// This file is auto-generated, don't edit it. Thanks.
package utils

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"horizon/config"
)

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *openapi.Client, _err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// Endpoint 请参考 https://api.aliyun.com/product/mse
	config.Endpoint = tea.String("mse.cn-shanghai.aliyuncs.com")
	_result = &openapi.Client{}
	_result, _err = openapi.NewClient(config)
	return _result, _err
}

/**
 * API 相关
 * @param path params
 * @return OpenApi.Params
 */
func CreateApiInfo(action string) (_result *openapi.Params) {
	params := &openapi.Params{
		// 接口名称
		Action: tea.String(action),
		// 接口版本
		Version: tea.String("2019-05-31"),
		// 接口协议
		Protocol: tea.String("HTTPS"),
		// 接口 HTTP 方法
		Method:   tea.String("POST"),
		AuthType: tea.String("AK"),
		Style:    tea.String("RPC"),
		// 接口 PATH
		Pathname: tea.String("/"),
		// 接口请求体内容格式
		ReqBodyType: tea.String("json"),
		// 接口响应体内容格式
		BodyType: tea.String("json"),
	}
	_result = params
	return _result
}

func ListClusters(pageNum int, pageSize int, clusterAliasName string) (rs map[string]interface{}, _err error) {
	client, _err := CreateClient(tea.String(config.Conf.Aliyun.AccessKeyId), tea.String(config.Conf.Aliyun.AccessKeySecret))
	if _err != nil {
		return rs, _err
	}

	action := "ListClusters"
	params := CreateApiInfo(action)
	// query params
	queries := map[string]interface{}{}
	queries["RegionId"] = tea.String("cn-shanghai")
	queries["PageNum"] = tea.Int(pageNum)
	queries["PageSize"] = tea.Int(pageSize)
	if clusterAliasName != "" {
		queries["ClusterAliasName"] = tea.String(clusterAliasName)
	}
	// runtime options
	runtime := &util.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}
	// 复制代码运行请自行打印 API 的返回值
	// 返回值为 Map 类型，可从 Map 中获得三类数据：响应体 body、响应头 headers、HTTP 返回的状态码 statusCode。
	rs, _err = client.CallApi(params, request, runtime)
	if _err != nil {
		return nil, _err
	}
	return rs, _err
}

func ListEngineNamespaces(instanceId string) (rs map[string]interface{}, _err error) {
	client, _err := CreateClient(tea.String(config.Conf.Aliyun.AccessKeyId), tea.String(config.Conf.Aliyun.AccessKeySecret))
	if _err != nil {
		return rs, _err
	}

	action := "ListEngineNamespaces"
	params := CreateApiInfo(action)
	// query params
	queries := map[string]interface{}{}
	queries["InstanceId"] = tea.String(instanceId)
	// runtime options
	runtime := &util.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}
	// 复制代码运行请自行打印 API 的返回值
	// 返回值为 Map 类型，可从 Map 中获得三类数据：响应体 body、响应头 headers、HTTP 返回的状态码 statusCode。
	rs, _err = client.CallApi(params, request, runtime)
	if _err != nil {
		return nil, _err
	}
	return rs, _err
}

func ListNacosConfigs(pageNum int, pageSize int, instanceId string, namespaceId string, dataId string) (rs map[string]interface{}, _err error) {
	client, _err := CreateClient(tea.String(config.Conf.Aliyun.AccessKeyId), tea.String(config.Conf.Aliyun.AccessKeySecret))
	if _err != nil {
		return rs, _err
	}

	action := "ListNacosConfigs"
	params := CreateApiInfo(action)
	// query params
	queries := map[string]interface{}{}
	queries["RegionId"] = tea.String("cn-shanghai")
	queries["PageNum"] = tea.Int(pageNum)
	queries["PageSize"] = tea.Int(pageSize)
	queries["InstanceId"] = tea.String(instanceId)
	queries["NamespaceId"] = tea.String(namespaceId)
	if dataId != "" {
		queries["DataId"] = tea.String(dataId)
	}
	// runtime options
	runtime := &util.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}
	// 复制代码运行请自行打印 API 的返回值
	// 返回值为 Map 类型，可从 Map 中获得三类数据：响应体 body、响应头 headers、HTTP 返回的状态码 statusCode。
	rs, _err = client.CallApi(params, request, runtime)
	if _err != nil {
		return nil, _err
	}
	return rs, _err
}

func GetNacosConfig(instanceId string, namespaceId string, dataId string, group string) (rs map[string]interface{}, _err error) {
	client, _err := CreateClient(tea.String(config.Conf.Aliyun.AccessKeyId), tea.String(config.Conf.Aliyun.AccessKeySecret))
	if _err != nil {
		return rs, _err
	}

	action := "GetNacosConfig"
	params := CreateApiInfo(action)
	// query params
	queries := map[string]interface{}{}
	queries["InstanceId"] = tea.String(instanceId)
	queries["NamespaceId"] = tea.String(namespaceId)
	queries["DataId"] = tea.String(dataId)
	queries["Group"] = tea.String(group)
	// runtime options
	runtime := &util.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}
	// 复制代码运行请自行打印 API 的返回值
	// 返回值为 Map 类型，可从 Map 中获得三类数据：响应体 body、响应头 headers、HTTP 返回的状态码 statusCode。
	rs, _err = client.CallApi(params, request, runtime)
	if _err != nil {
		return nil, _err
	}
	return rs, _err
}
