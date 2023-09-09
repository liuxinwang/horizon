package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"horizon/handler"
	"horizon/utils"
	"io"
	"log"
	"net/http"
	"os"
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
	// logo handle
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	gin.DisableConsoleColor()
	// 记录到文件。
	f, _ := os.Create("horizon.log")
	gin.DefaultWriter = io.MultiWriter(f)
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
	user.GET("", handler.GetUserList)
	user.PUT("", handler.PutUser)
	user.PUT("resetPassword", handler.PutResetPassword)

	// role group
	role := api.Group("/role", authMiddleware.MiddlewareFunc())
	role.GET("", handler.RoleGet)
	role.POST("", handler.RolePost)
	role.PUT("", handler.RolePut)
	role.DELETE("/:id", handler.RoleDelete)
	role.POST("permission", handler.RolePermissionPost)

	// menu group
	menu := api.Group("/menu", authMiddleware.MiddlewareFunc())
	menu.GET("", handler.MenuGet)

	// instance group
	instance := api.Group("/instance", authMiddleware.MiddlewareFunc())
	instance.GET("", handler.InstanceGet)
	instance.GET("/:instId", handler.InstanceInstIdGet)
	instance.POST("", handler.InstancePost)
	instance.PUT("", handler.InstancePut)
	instance.DELETE("/:id", handler.InstanceDelete)
	instance.GET("db/:instId", handler.InstanceDdGet)

	// inspection group
	inspection := api.Group("/inspection", authMiddleware.MiddlewareFunc())
	inspection.GET("", handler.InspectionGet)
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
	project.POST("", handler.ProjectPost)
	project.PUT("", handler.ProjectPut)
	project.DELETE("/:id", handler.ProjectDelete)
	project.POST("resource/config", handler.ProjectResourceConfigPost)
	project.GET("role", handler.ProjectRoleGet)
	project.GET("user/:userName", handler.ProjectUserNameGet)
	project.GET("datasource/:projId", handler.ProjectDataSourceGet)

	// workflow group
	workflow := sqlAudit.Group("/workflow", authMiddleware.MiddlewareFunc())
	workflow.GET("", handler.WorkflowGet)
	workflow.GET("/:id", handler.WorkflowIdGet)
	workflow.POST("", handler.WorkflowPost)
	workflow.PUT("", handler.WorkflowPut)
	workflow.DELETE("/:id", handler.WorkflowDelete)
	workflow.GET("progress/:id", handler.WorkflowIdProgressGet)
	workflow.POST("audit", handler.WorkflowAuditPost)
	workflow.POST("cancel", handler.WorkflowCancelPost)
	workflow.POST("execute", handler.WorkflowExecutePost)

	// workflowTemplate group
	workflowTemplate := sqlAudit.Group("/workflowTemplate", authMiddleware.MiddlewareFunc())
	workflowTemplate.GET("", handler.WorkflowTemplateGet)
	workflowTemplate.POST("", handler.WorkflowTemplatePost)
	workflowTemplate.PUT("", handler.WorkflowTemplatePut)
	workflowTemplate.DELETE("/:id", handler.WorkflowTemplateDelete)
	workflowTemplate.POST("/config", handler.WorkflowTemplateConfigPost)
}
