package service

import (
	"github.com/gin-gonic/gin"
	"horizon/model"
	"math"
	"net/http"
	"strconv"
)

// MenuSelectByList 查询菜单列表
func MenuSelectByList(c *gin.Context) {
	// 变量初始化
	pageNo, _ := strconv.Atoi(c.Query("pageNo"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	var Db = model.Db
	var menus []model.Menu
	var totalCount int64
	data := gin.H{"totalCount": 0, "data": &[]model.Menu{}, "pageNo": pageNo, "pageSize": pageSize, "totalPage": 0}

	// 查询条件处理

	// 执行查询
	Db.Order("id").Limit(pageSize).Offset((pageNo-1)*pageSize - 1).Find(&menus)
	Db.Model(&model.Menu{}).Count(&totalCount)

	// 处理结果集并返回
	totalPage := math.Ceil(float64(totalCount) / float64(pageSize))
	if totalCount > 0 {
		data = gin.H{"totalCount": totalCount, "data": &menus, "pageNo": pageNo, "pageSize": pageSize, "totalPage": totalPage}
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": data, "err": ""})
}
