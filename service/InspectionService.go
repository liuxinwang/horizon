package service

import (
	"github.com/gin-gonic/gin"
	"horizon/model"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// InspectionSelectByList 查询实例巡检列表
func InspectionSelectByList(c *gin.Context) {
	// 变量初始化
	pageNo, _ := strconv.Atoi(c.Query("pageNo"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	var Db = model.Db
	var inspections []model.Inspection
	var totalCount int64
	data := gin.H{"total": 0, "data": &[]model.Inspection{}, "pageNo": pageNo, "pageSize": pageSize, "totalPage": 0}

	// 查询条件处理
	if instId, isExist := c.GetQuery("InstId"); isExist == true && strings.Trim(instId, " ") != "" {
		Db = Db.Where("inst_id like ?", "%"+instId+"%")
	}
	if instName, isExist := c.GetQuery("InstName"); isExist == true && strings.Trim(instName, " ") != "" {
		Db = Db.Where("inst_name like ?", "%"+instName+"%")
	}

	// 执行查询
	Db.Order("created_at desc").Limit(pageSize).Offset((pageNo-1)*pageSize - 1).Find(&inspections)
	Db.Model(&model.Inspection{}).Count(&totalCount)

	// 处理结果集并返回
	totalPage := math.Ceil(float64(totalCount) / float64(pageSize))
	if totalCount > 0 {
		data = gin.H{"totalCount": totalCount, "data": &inspections, "pageNo": pageNo, "pageSize": pageSize, "totalPage": totalPage}
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": data, "err": ""})
}

// InspectionSelectByInspId 查看实例巡检详情
func InspectionSelectByInspId(c *gin.Context) {
	// 变量初始化
	inspId := c.Param("id")
	var inspDetails []model.InspDetail
	// 执行查询
	model.Db.Where("insp_id = ?", inspId).Find(&inspDetails)
	// 返回结果集
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": &inspDetails, "err": ""})
}
