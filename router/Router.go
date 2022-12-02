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
	user := api.Group("/user")
	user.GET("/info", authMiddleware.MiddlewareFunc(), handler.GetUser)

	// instance group
	instance := api.Group("/instance", authMiddleware.MiddlewareFunc())
	instance.GET("", handler.InstanceGet)
	instance.GET("/:instId", handler.InstanceInstIdGet)
	instance.POST("", handler.InstancePost)
	instance.PUT("", handler.InstancePut)
	instance.DELETE("/:id", handler.InstanceDelete)

	// inspection group
	inspection := api.Group("/inspection", authMiddleware.MiddlewareFunc())
	inspection.GET("", handler.InspectionGet)
	inspection.GET("/:id", handler.InspectionDetailGet)

	// score group
	score := api.Group("/score", authMiddleware.MiddlewareFunc())
	score.GET("", handler.ScoreGet)
}
