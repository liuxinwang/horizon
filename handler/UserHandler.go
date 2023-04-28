package handler

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"horizon/model"
	"horizon/service"
	"net/http"
)

func GetUserList(c *gin.Context) {
	service.UserSelectByList(c)
}

func PutUser(c *gin.Context) {
	service.UserUpdate(c)
}

func PutResetPassword(c *gin.Context) {
	service.UserUpdatePassword(c)
}

func GetUser(c *gin.Context) {
	service.UserSelectByUserName(c)
}

func GetUserNav(c *gin.Context) {
	type Data struct {
		ID   int
		Name string
		Age  int
	}
	var menus []model.Menu
	claims := jwt.ExtractClaims(c)
	sql := "select m.* " +
		"from users u " +
		"inner join user_roles ur on u.id = ur.user_id " +
		"inner join role_permissions rp on ur.role_id = rp.role_id " +
		"inner join menus m on rp.menu_id = m.id " +
		"where u.user_name = ?"
	model.Db.Raw(sql, claims["UserName"].(string)).Scan(&menus)
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "success",
		"data": &menus,
		"err":  "",
	})
}
