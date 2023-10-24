package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"horizon/config"
	"horizon/model"
	"horizon/utils"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// InstanceSelectByList 查询实例列表
func InstanceSelectByList(c *gin.Context) {
	// 变量初始化
	pageNo, _ := strconv.Atoi(c.Query("pageNo"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	var Db = model.Db
	var instances []model.Instance
	var totalCount int64
	data := gin.H{"totalCount": 0, "data": &[]model.Instance{}, "pageNo": pageNo, "pageSize": pageSize, "totalPage": 0}

	// 查询条件处理
	if instId, isExist := c.GetQuery("InstId"); isExist == true && strings.Trim(instId, " ") != "" {
		Db = Db.Where("inst_id like ?", "%"+instId+"%")
	}
	if name, isExist := c.GetQuery("Name"); isExist == true && strings.Trim(name, " ") != "" {
		Db = Db.Where("name like ?", "%"+name+"%")
	}
	if status, isExist := c.GetQuery("Status"); isExist == true && status != "0" {
		Db = Db.Where("status = ?", status)
	}
	if ip, isExist := c.GetQuery("Ip"); isExist == true && strings.Trim(ip, " ") != "" {
		Db = Db.Where("ip like ?", "%"+ip+"%")
	}

	// 执行查询
	Db.Order("created_at desc").Limit(pageSize).Offset((pageNo-1)*pageSize - 1).Find(&instances)
	Db.Model(&model.Instance{}).Count(&totalCount)

	// 处理结果集并返回
	totalPage := math.Ceil(float64(totalCount) / float64(pageSize))
	if totalCount > 0 {
		data = gin.H{"totalCount": totalCount, "data": &instances, "pageNo": pageNo, "pageSize": pageSize, "totalPage": totalPage}
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": data, "err": ""})
}

// InstanceSelectByInstId 查看实例信息
func InstanceSelectByInstId(c *gin.Context) {
	var Db = model.Db
	var instance model.Instance
	// 执行查询
	Db.Where("inst_id = ?", c.Param("instId")).Find(&instance)
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": &instance, "err": ""})
}

// InstanceInsert 新增实例
func InstanceInsert(c *gin.Context) {
	// 参数映射到对象
	var instance model.Instance
	if err := c.ShouldBind(&instance); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	// 加密密码
	instance.Password = encryptPassword(instance.Password)
	// 获取状态
	if _, err := getStatus(&instance); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	} else {
		getVersion(&instance)
	}
	// 获取实例ID
	instId := utils.GenerateId(&instance)
	instance.InstId = instId
	result := model.Db.Create(&instance)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
	}
}

func InstanceUpdate(c *gin.Context) {
	// 参数映射到对象
	var instance model.Instance
	if err := c.ShouldBind(&instance); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	// 密码处理
	var instance2 model.Instance
	model.Db.Select("password").Where("inst_id = ?", instance.InstId).First(&instance2)
	if instance.Password != instance2.Password {
		// 加密密码
		instance.Password = encryptPassword(instance.Password)
		// 验证状态
		if _, err := getStatus(&instance); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
			return
		} else {
			getVersion(&instance)
		}
	}
	// 执行更新
	updMap := map[string]interface{}{"name": instance.Name, "type": instance.Type,
		"env_type": instance.EnvType, "role": instance.Role, "ip": instance.Ip,
		"port": instance.Port, "user": instance.User, "password": instance.Password}
	model.Db.Model(model.Instance{}).Where("inst_id = ?", instance.InstId).Updates(updMap)
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
}

func InstanceDelete(c *gin.Context) {
	id := c.Param("id")
	result := model.Db.Delete(&model.Instance{}, id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": "数据不存在"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
	}
}

// InstanceDbSelectByInstId 查询列表
func InstanceDbSelectByInstId(c *gin.Context) {
	var instance model.Instance
	model.Db.Where("inst_id = ?", c.Param("instId")).First(&instance)
	databases, err := getDatabases(&instance)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": &databases, "err": ""})
}

// 获取实例状态
func getStatus(instance *model.Instance) (*gorm.DB, error) {
	dsn := "%s:%s@tcp(%s:%d)/information_schema?charset=utf8mb4&parseTime=True&loc=Local&timeout=1s"
	dsn = fmt.Sprintf(dsn, instance.User, utils.DecryptAES([]byte(config.Conf.General.SecretKey), instance.Password), instance.Ip, instance.Port)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		instance.Status = "Error"
	} else {
		instance.Status = "Running"
	}
	return db, err
}

// 生成实例ID
func generateId(instance *model.Instance) {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(letterRunes))]
	}
	instance.InstId = "mysql-" + string(b)
}

// 获取实例版本
func getVersion(instance *model.Instance) {
	dsn := "%s:%s@tcp(%s:%d)/information_schema?charset=utf8mb4&parseTime=True&loc=Local&timeout=1s"
	dsn = fmt.Sprintf(dsn, instance.User, utils.DecryptAES([]byte(config.Conf.General.SecretKey), instance.Password), instance.Ip, instance.Port)
	Db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("ERROR occurred:" + err.Error())
	}
	var version string
	Db.Raw("SELECT VERSION()").Scan(&version)
	instance.Version = version
}

// 连接密码加密
func encryptPassword(password string) string {
	// bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return utils.EncryptAES([]byte(config.Conf.General.SecretKey), password)
}

// 获取实例数据库
func getDatabases(instance *model.Instance) (interface{}, error) {
	dsn := "%s:%s@tcp(%s:%d)/information_schema?charset=utf8mb4&parseTime=True&loc=Local&timeout=1s"
	dsn = fmt.Sprintf(dsn, instance.User, utils.DecryptAES([]byte(config.Conf.General.SecretKey), instance.Password), instance.Ip, instance.Port)
	Db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	type result struct {
		Name string `json:"name"`
	}
	var results []result
	Db.Raw("SELECT SCHEMA_NAME as name " +
		"FROM SCHEMATA " +
		"WHERE SCHEMA_NAME NOT IN ('information_schema', 'mysql', 'performance_schema' , '')").Scan(&results)
	return results, nil
}
