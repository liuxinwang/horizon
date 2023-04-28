package service

import (
	"github.com/gin-gonic/gin"
	"horizon/model"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// RoleSelectByList 查询角色列表
func RoleSelectByList(c *gin.Context) {
	// 变量初始化
	pageNo, _ := strconv.Atoi(c.Query("pageNo"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	var Db = model.Db
	var roles []model.Role
	var totalCount int64
	data := gin.H{"totalCount": 0, "data": &[]model.Role{}, "pageNo": pageNo, "pageSize": pageSize, "totalPage": 0}

	// 查询条件处理
	if key, isExist := c.GetQuery("id"); isExist == true && strings.Trim(key, " ") != "" {
		Db = Db.Where("id like ?", "%"+key+"%")
	}
	if name, isExist := c.GetQuery("name"); isExist == true && strings.Trim(name, " ") != "" {
		Db = Db.Where("name like ?", "%"+name+"%")
	}

	// 执行查询
	Db.Preload("UserRoles").Preload("RolePermissions").
		Preload("RolePermissions.Menu").
		Order("created_at desc").Limit(pageSize).
		Offset((pageNo-1)*pageSize - 1).Find(&roles)
	Db.Model(&model.Role{}).Count(&totalCount)

	// 处理结果集并返回
	totalPage := math.Ceil(float64(totalCount) / float64(pageSize))
	if totalCount > 0 {
		data = gin.H{"totalCount": totalCount, "data": &roles, "pageNo": pageNo, "pageSize": pageSize, "totalPage": totalPage}
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": data, "err": ""})
}

// RoleInsert 新增角色
func RoleInsert(c *gin.Context) {
	// 参数映射到对象
	var role model.Role
	if err := c.ShouldBind(&role); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	result := model.Db.Create(&role)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
	}
}

func RoleUpdate(c *gin.Context) {
	// 参数映射到对象
	var role model.Role
	if err := c.ShouldBind(&role); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	// 执行更新
	updMap := map[string]interface{}{"name": role.Name}
	model.Db.Model(model.Role{}).Where("id = ?", role.ID).Updates(updMap)
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
}

func RoleDelete(c *gin.Context) {
	id := c.Param("id")
	result := model.Db.Delete(&model.Role{ID: id})
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": "数据不存在"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
	}
}
