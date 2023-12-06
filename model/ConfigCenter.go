package model

type ConfigCenterInstance struct {
	MseVersion        string `json:"mseVersion"`
	InternetAddress   string `json:"internetAddress"`
	ResourceGroupId   string `json:"resourceGroupId"`
	InstanceId        string `json:"instanceId"`
	ClusterId         string `json:"clusterId"`
	CreateTime        string `json:"createTime"`
	ClusterType       string `json:"clusterType"`
	EndDate           string `json:"endDate"`
	ClusterAliasName  string `json:"clusterAliasName"`
	AppVersion        string `json:"appVersion"`
	VersionCode       string `json:"versionCode"`
	InstanceCount     int    `json:"instanceCount"`
	IntranetAddress   string `json:"intranetAddress"`
	CanUpdate         bool   `json:"canUpdate"`
	VpcId             string `json:"vpcId"`
	ChargeType        string `json:"chargeType"`
	ClusterName       string `json:"clusterName"`
	InitStatus        string `json:"initStatus"`
	MaintenancePeriod struct {
	} `json:"maintenancePeriod"`
	IntranetDomain string `json:"intranetDomain"`
	Tags           struct {
		AcsRmRgId string `json:"acs:rm:rgId"`
	} `json:"tags"`
	Namespaces []ConfigCenterNamespace `json:"namespaces"`
}

type ConfigCenterNamespace struct {
	Type              int    `json:"type"`
	Quota             int    `json:"quota"`
	ConfigCount       int    `json:"configCount"`
	SourceType        string `json:"sourceType"`
	NamespaceShowName string `json:"namespaceShowName"`
	Namespace         string `json:"namespace"`
	ServiceCount      int    `json:"serviceCount"`
	NamespaceId       string `json:"namespaceId" form:"namespaceId"`
}

type ConfigCenterInstanceConfig struct {
	Group            string `json:"group" form:"group"`
	Desc             string `json:"desc" form:"desc"`
	Type             string `json:"type" form:"type"`
	DataId           string `json:"dataId" form:"dataId"`
	EncryptedDataKey string `json:"encryptedDataKey" form:"encryptedDataKey"`
	Content          string `json:"content" form:"content"`
	AppName          string `json:"appName" form:"appName"`
	Md5              string `json:"md5" form:"md5"`
	InstanceId       string `json:"instanceId" form:"instanceId"`
	NamespaceId      string `json:"namespaceId" form:"namespaceId"`
}
