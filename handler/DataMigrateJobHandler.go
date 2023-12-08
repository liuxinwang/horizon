package handler

import (
	"github.com/gin-gonic/gin"
	"horizon/service"
)

func DataMigrateJobGet(c *gin.Context) {
	service.DataMigrateJobSelectByList(c)
}

func DataMigrateJobIdGet(c *gin.Context) {
	service.DataMigrateJobSelectById(c)
}

func DataMigrateJobPost(c *gin.Context) {
	service.DataMigrateJobInsert(c)
}

func DataMigrateJobExecutePost(c *gin.Context) {
	service.DataMigrateJobExecuteUpdate(c)
}

func DataMigrateJobDetailGet(c *gin.Context) {
	service.DataMigrateJobDetailSelectByList(c)
}
