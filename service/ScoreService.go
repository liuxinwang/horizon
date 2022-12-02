package service

import (
	"github.com/gin-gonic/gin"
	"horizon/model"
	"net/http"
)

// ScoreSelectByList 查询实例巡检扣分列表
func ScoreSelectByList(c *gin.Context) {
	// 变量初始化
	inspId := c.Query("InspId")
	var scores []model.Score
	// 执行查询
	model.Db.Where("insp_id = ?", inspId).Order("deduction desc").Find(&scores)
	// 返回结果集
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": &scores, "err": ""})
}
