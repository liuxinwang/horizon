package handler

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"horizon/model"
	"net/http"
)

func GetUser(c *gin.Context) {
	var user model.User
	claims := jwt.ExtractClaims(c)
	model.Db.Where("user_name = ?", claims["UserName"].(string)).First(&user)
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "success",
		"data": &user,
		"err":  "",
	})
}
