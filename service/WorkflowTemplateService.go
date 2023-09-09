package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"horizon/model"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// WorkflowTemplateSelectByList 查询列表
func WorkflowTemplateSelectByList(c *gin.Context) {
	// 变量初始化
	pageNo, _ := strconv.Atoi(c.Query("pageNo"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	var Db = model.Db
	var workflowTemplates []model.WorkflowTemplate
	var totalCount int64
	data := gin.H{"totalCount": 0, "data": &[]model.WorkflowTemplate{}, "pageNo": pageNo, "pageSize": pageSize, "totalPage": 0}

	// 查询条件处理
	if name, isExist := c.GetQuery("name"); isExist == true && strings.Trim(name, " ") != "" {
		Db = Db.Where("name like ?", "%"+name+"%")
	}

	// 执行查询
	Db.Preload("WorkflowTemplateDetails", func(db *gorm.DB) *gorm.DB {
		return db.Order("id asc")
	}).Order("created_at desc").Limit(pageSize).Offset((pageNo-1)*pageSize - 1).Find(&workflowTemplates)
	Db.Model(&model.WorkflowTemplate{}).Count(&totalCount)

	// 处理结果集并返回
	totalPage := math.Ceil(float64(totalCount) / float64(pageSize))
	if totalCount > 0 {
		data = gin.H{"totalCount": totalCount, "data": &workflowTemplates, "pageNo": pageNo, "pageSize": pageSize, "totalPage": totalPage}
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": data, "err": ""})
}

// WorkflowTemplateInsert 新增
func WorkflowTemplateInsert(c *gin.Context) {
	// 参数映射到对象
	var workflowTemplate model.WorkflowTemplate
	if err := c.ShouldBind(&workflowTemplate); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	type Result struct {
		MaxCode uint
	}
	var rs Result
	// 获取ID
	model.Db.Model(&model.WorkflowTemplate{}).Select("MAX(code) as max_code").Scan(&rs)
	if rs.MaxCode != 0 {
		workflowTemplate.Code = rs.MaxCode + 1
	} else {
		workflowTemplate.Code = 10001
	}

	result := model.Db.Create(&workflowTemplate)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
	}
}

func WorkflowTemplateUpdate(c *gin.Context) {
	// 参数映射到对象
	var workflowTemplate model.WorkflowTemplate
	if err := c.ShouldBind(&workflowTemplate); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	// 执行更新
	updMap := map[string]interface{}{"name": workflowTemplate.Name}
	model.Db.Model(model.WorkflowTemplate{}).Where("id = ?", workflowTemplate.ID).Updates(updMap)
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
}

func WorkflowTemplateDelete(c *gin.Context) {
	id := c.Param("id")
	result := model.Db.Delete(&model.WorkflowTemplate{}, id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": "记录不存在"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
	}
}

func WorkflowTemplateConfigInsert(c *gin.Context) {
	var permBody struct {
		model.WorkflowTemplate
		NodeNames      []string `json:"nodeNames"`
		ProjectRoleIds []string `json:"projectRoleIds"`
	}
	if err := c.ShouldBind(&permBody); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	tx := model.Db.Begin()
	// first delete and add WorkflowTemplateDetail
	result := tx.Where("workflow_template_id = ?", permBody.WorkflowTemplate.ID).Delete(&model.WorkflowTemplateDetail{})
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
		tx.Rollback()
		return
	}
	for i, nodeName := range permBody.NodeNames {
		var workflowTemplateDetail model.WorkflowTemplateDetail
		workflowTemplateDetail.WorkflowTemplateId = permBody.WorkflowTemplate.ID
		workflowTemplateDetail.WorkflowTemplateCode = permBody.WorkflowTemplate.Code
		workflowTemplateDetail.SerialNumber = uint(i + 1)
		workflowTemplateDetail.NodeName = nodeName
		workflowTemplateDetail.ProjectRoleId = permBody.ProjectRoleIds[i]
		result := tx.Create(&workflowTemplateDetail)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
}
