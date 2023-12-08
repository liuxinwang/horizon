package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"horizon/handler"
	"horizon/utils"
	"log"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.NoRoute(func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	commonRouter(r)
	return r
}

func InitRouterPack() *gin.Engine {
	// init gin
	r := gin.Default()
	// index handle
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	// frontend history mode handle
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	commonRouter(r)
	return r
}

func commonRouter(r *gin.Engine) {
	authMiddleware := utils.JWTAuthMiddleware()
	r.GET("/ping", handler.Ping)

	api := r.Group("/api")

	// auth group
	auth := api.Group("/auth")
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.POST("/login", authMiddleware.LoginHandler)
	auth.POST("/logout", authMiddleware.LogoutHandler)

	// user group
	user := api.Group("/user", authMiddleware.MiddlewareFunc())
	user.GET("/info", handler.GetUser)
	user.GET("/nav", handler.GetUserNav)
	user.GET("/query", handler.GetUserList)
	user.PUT("/edit", handler.PutUser)
	user.PUT("resetPassword", handler.PutResetPassword)

	// role group
	role := api.Group("/role", authMiddleware.MiddlewareFunc())
	role.GET("query", handler.RoleGet)
	role.POST("add", handler.RolePost)
	role.PUT("edit", handler.RolePut)
	role.DELETE("delete/:id", handler.RoleDelete)
	role.POST("permission", handler.RolePermissionPost)

	// menu group
	menu := api.Group("/menu", authMiddleware.MiddlewareFunc())
	menu.GET("", handler.MenuGet)

	// instance group
	instance := api.Group("/instance", authMiddleware.MiddlewareFunc())
	instance.GET("query", handler.InstanceGet)
	instance.GET("/:instId", handler.InstanceInstIdGet)
	instance.POST("add", handler.InstancePost)
	instance.PUT("edit", handler.InstancePut)
	instance.DELETE("delete/:id", handler.InstanceDelete)
	instance.GET("db/:instId", handler.InstanceDbGet)
	instance.GET("db/table", handler.InstanceDbTableGet)

	// inspection group
	inspection := api.Group("/inspection", authMiddleware.MiddlewareFunc())
	inspection.GET("query", handler.InspectionGet)
	inspection.GET("/:id", handler.InspectionDetailGet)

	// score group
	score := api.Group("/score", authMiddleware.MiddlewareFunc())
	score.GET("", handler.ScoreGet)

	// sql audit group
	sqlAudit := api.Group("/sqlaudit", authMiddleware.MiddlewareFunc())
	// project group
	project := sqlAudit.Group("/project", authMiddleware.MiddlewareFunc())
	project.GET("", handler.ProjectGet)
	project.GET("/:projId", handler.ProjectProjIdGet)
	project.POST("add", handler.ProjectPost)
	project.PUT("edit", handler.ProjectPut)
	project.DELETE("/:id", handler.ProjectDelete)
	project.POST("resource/config", handler.ProjectResourceConfigPost)
	project.GET("role", handler.ProjectRoleGet)
	project.GET("user/:userName", handler.ProjectUserNameGet)
	project.GET("datasource/:projId", handler.ProjectDataSourceGet)

	// workflow group
	workflow := sqlAudit.Group("/workflow", authMiddleware.MiddlewareFunc())
	workflow.GET("", handler.WorkflowGet)
	workflow.GET("/:id", handler.WorkflowIdGet)
	workflow.POST("add", handler.WorkflowPost)
	workflow.PUT("edit", handler.WorkflowPut)
	workflow.DELETE("/:id", handler.WorkflowDelete)
	workflow.GET("progress/:id", handler.WorkflowIdProgressGet)
	workflow.POST("audit", handler.WorkflowAuditPost)
	workflow.POST("cancel", handler.WorkflowCancelPost)
	workflow.POST("execute", handler.WorkflowExecutePost)
	workflow.POST("scheduledExecute", handler.WorkflowScheduledExecutionPost)
	workflow.GET("sqlDetail", handler.WorkflowIdSqlDetailGet)

	// workflowTemplate group
	workflowTemplate := sqlAudit.Group("/workflowTemplate", authMiddleware.MiddlewareFunc())
	workflowTemplate.GET("", handler.WorkflowTemplateGet)
	workflowTemplate.POST("add", handler.WorkflowTemplatePost)
	workflowTemplate.PUT("edit", handler.WorkflowTemplatePut)
	workflowTemplate.DELETE("/:id", handler.WorkflowTemplateDelete)
	workflowTemplate.POST("/config", handler.WorkflowTemplateConfigPost)

	// configCenter group
	configCenter := api.Group("/configCenter", authMiddleware.MiddlewareFunc())
	configCenterInstance := configCenter.Group("/instance", authMiddleware.MiddlewareFunc())
	configCenterInstance.GET("query", handler.ConfigCenterInstanceGet)

	configCenterInstanceConfig := configCenter.Group("/config", authMiddleware.MiddlewareFunc())
	configCenterInstanceConfig.GET("query", handler.ConfigCenterInstanceConfigGet)
	configCenterInstanceConfig.GET("detail", handler.ConfigCenterInstanceConfigDataGet)

	// dataManger group
	dataManger := api.Group("/dataManger", authMiddleware.MiddlewareFunc())
	dataMigrateJob := dataManger.Group("/dataMigrateJob", authMiddleware.MiddlewareFunc())
	dataMigrateJob.GET("query", handler.DataMigrateJobGet)
	dataMigrateJob.GET("/:id", handler.DataMigrateJobIdGet)
	dataMigrateJob.POST("add", handler.DataMigrateJobPost)
	dataMigrateJob.POST("execute", handler.DataMigrateJobExecutePost)
	dataMigrateJobDetail := dataManger.Group("/dataMigrateJobDetail", authMiddleware.MiddlewareFunc())
	dataMigrateJobDetail.GET("query", handler.DataMigrateJobDetailGet)
}
