package service

import (
	"github.com/gin-gonic/gin"
	"horizon/model"
	"horizon/utils"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// WorkflowSelectByList 查询列表
func WorkflowSelectByList(c *gin.Context) {
	// 变量初始化
	pageNo, _ := strconv.Atoi(c.Query("pageNo"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	var Db = model.Db
	var workflows []model.Workflow
	var totalCount int64
	data := gin.H{"totalCount": 0, "data": &[]model.Workflow{}, "pageNo": pageNo, "pageSize": pageSize, "totalPage": 0}

	// 查询条件处理
	if name, isExist := c.GetQuery("Name"); isExist == true && strings.Trim(name, " ") != "" {
		Db = Db.Where("name like ?", "%"+name+"%")
	}
	if status, isExist := c.GetQuery("Status"); isExist == true && status != "0" {
		Db = Db.Where("status = ?", status)
	}

	// 执行查询
	Db.Order("created_at desc").Limit(pageSize).Offset((pageNo-1)*pageSize - 1).Find(&workflows)
	Db.Model(&model.Workflow{}).Count(&totalCount)

	// 处理结果集并返回
	totalPage := math.Ceil(float64(totalCount) / float64(pageSize))
	if totalCount > 0 {
		data = gin.H{"totalCount": totalCount, "data": &workflows, "pageNo": pageNo, "pageSize": pageSize, "totalPage": totalPage}
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": data, "err": ""})
}

// WorkflowSelectByInstId 查看实例信息
func WorkflowSelectByInstId(c *gin.Context) {
	var Db = model.Db
	var workflow model.Workflow
	// 执行查询
	Db.Where("proj_id = ?", c.Param("instId")).Find(&workflow)
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": &workflow, "err": ""})
}

// WorkflowInsert 新增
func WorkflowInsert(c *gin.Context) {
	// 参数映射到对象
	var workflow model.Workflow
	if err := c.ShouldBind(&workflow); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	// 获取ID
	workflow.ProjId = utils.GenerateId(&workflow)
	result := model.Db.Create(&workflow)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
	}
}

func WorkflowUpdate(c *gin.Context) {
	// 参数映射到对象
	var workflow model.Workflow
	if err := c.ShouldBind(&workflow); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	// 执行更新
	updMap := map[string]interface{}{"name": workflow.Name, "describe": workflow.Describe}
	model.Db.Model(model.Workflow{}).Where("proj_id = ?", workflow.ProjId).Updates(updMap)
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
}

func WorkflowDelete(c *gin.Context) {
	id := c.Param("id")
	result := model.Db.Delete(&model.Workflow{}, id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": "项目不存在"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
	}
}
