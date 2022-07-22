package utils

import (
	"errors"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"horizon/config"
	"horizon/model"
	"log"
	"time"
)

func JWTAuthMiddleware() *jwt.GinJWTMiddleware {
	identityKey := "UserName"
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "horizon zone",
		Key:         []byte(config.Conf.General.SecretKey),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		// 此函数在成功验证（登录）后调用
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		// 从嵌入在 jwt 令牌中的声明中获取用户身份，并将此身份值传递给Authorizator
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &model.User{
				UserName: claims[identityKey].(string),
			}
		},
		// 此函数应验证给定 gin 上下文的用户凭据（即密码与给定用户电子邮件的哈希密码匹配，以及任何其他身份验证逻辑）
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals model.User
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			// 获取用户信息
			var user model.User
			result := model.Db.Where("user_name = ?", loginVals.UserName).First(&user)
			// 用户名验证
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, jwt.ErrFailedAuthentication
			}
			// 密码验证
			passwordValidation := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginVals.Password))
			if passwordValidation != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			// 返回用户信息
			user.Password = ""
			return &user, nil
		},
		// 给定用户身份值（data参数）和 gin 上下文，此函数应检查用户是否有权到达此端点
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*model.User); ok && v.UserName == "admin" {
				return true
			}
			return false
		},
		// 在登录、授权用户出现任何错误，或者请求中没有令牌的情况或情况下，将发生以下情况
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}
	return authMiddleware
}
