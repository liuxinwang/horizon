package handler

import (
	"github.com/gin-gonic/gin"
	"horizon/service"
)

func ProjectGet(c *gin.Context) {
	service.ProjectSelectByList(c)
}

func ProjectProjIdGet(c *gin.Context) {
	service.ProjectSelectByInstId(c)
}

func ProjectPost(c *gin.Context) {
	service.ProjectInsert(c)
}

func ProjectPut(c *gin.Context) {
	service.ProjectUpdate(c)
}

func ProjectDelete(c *gin.Context) {
	service.ProjectDelete(c)
}

func ProjectResourceConfigPost(c *gin.Context) {
	service.ProjectResourceConfigInsert(c)
}

func ProjectRoleGet(c *gin.Context) {
	service.ProjectRoleSelectByList(c)
}

func ProjectUserNameGet(c *gin.Context) {
	service.ProjectUserSelectByUserName(c)
}

func ProjectDataSourceGet(c *gin.Context) {
	service.ProjectDataSourceSelectByProjId(c)
}
