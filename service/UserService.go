package service

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"horizon/model"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type userListResult struct {
	ID        uint      `json:"id"`
	UserName  string    `json:"userName"`
	NickName  string    `json:"nickName"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	RoleId    string    `json:"roleId"`
	RoleName  string    `json:"roleName"`
}

type Perm struct {
	RoleId       string         `json:"roleId"`
	PermissionId string         `json:"permissionId"`
	Actions      datatypes.JSON `json:"actions"`
	ActionList   datatypes.JSON `json:"actionList"`
}
type RolePerm struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Describe string `json:"describe"`
	Perms    []Perm `gorm:"foreignKey:RoleId;references:ID" json:"permissions"`
}
type UserNameResult struct {
	User *model.User `json:"user"`
	Role *RolePerm   `json:"role"`
}

// UserSelectByList 查询用户列表
func UserSelectByList(c *gin.Context) {
	// 变量初始化
	pageNo, _ := strconv.Atoi(c.Query("pageNo"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	var Db = model.Db
	var totalCount int64
	data := gin.H{"totalCount": 0, "data": &[]userListResult{}, "pageNo": pageNo, "pageSize": pageSize, "totalPage": 0}

	// 查询条件处理
	if userName, isExist := c.GetQuery("userName"); isExist == true && strings.Trim(userName, " ") != "" {
		Db = Db.Where("users.user_name like ?", "%"+userName+"%")
	}
	if nickName, isExist := c.GetQuery("nickName"); isExist == true && strings.Trim(nickName, " ") != "" {
		Db = Db.Where("users.nick_name like ?", "%"+nickName+"%")
	}
	if status, isExist := c.GetQuery("status"); isExist == true && status != "0" {
		Db = Db.Where("users.status = ?", status)
	}

	var userListResults []userListResult
	// 执行查询
	Db.Table("users").Select("users.*, roles.id as role_id, roles.name as role_name").
		Joins("left join user_roles on users.id = user_roles.user_id").
		Joins("left join roles on user_roles.role_id = roles.id").
		Order("created_at desc").
		Limit(pageSize).Offset((pageNo-1)*pageSize - 1).Find(&userListResults)
	Db.Table("users").Select("users.*, roles.id as role_id, roles.name as role_name").
		Joins("left join user_roles on users.id = user_roles.user_id").
		Joins("left join roles on user_roles.role_id = roles.id").Count(&totalCount)
	// Db.Preload("UserRoles").Order("created_at desc").Limit(pageSize).Offset((pageNo-1)*pageSize - 1).Find(&users)
	// Db.Model(&model.User{}).Count(&totalCount)

	// 处理结果集并返回
	totalPage := math.Ceil(float64(totalCount) / float64(pageSize))
	if totalCount > 0 {
		data = gin.H{"totalCount": totalCount, "data": &userListResults, "pageNo": pageNo, "pageSize": pageSize, "totalPage": totalPage}
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": data, "err": ""})
}

func UserUpdate(c *gin.Context) {
	// 参数映射到对象
	var userBody struct {
		ID       uint   `json:"id"`
		UserName string `json:"userName"`
		NickName string `json:"nickName"`
		RoleId   string `json:"roleId"`
	}
	if err := c.ShouldBind(&userBody); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	tx := model.Db.Begin()
	// 执行更新
	updMap := map[string]interface{}{"nick_name": userBody.NickName}
	if err := tx.Model(model.User{}).Where("user_name = ?", userBody.UserName).Updates(updMap).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	if err := tx.Where("user_id = ?", userBody.ID).Delete(&model.UserRole{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	if err := tx.Create(&model.UserRole{UserId: userBody.ID, RoleId: userBody.RoleId}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
}

func UserUpdatePassword(c *gin.Context) {
	// 参数映射到对象
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	// 密码处理
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	user.Password = string(bcryptPassword)
	// 执行更新
	updMap := map[string]interface{}{"password": user.Password}
	model.Db.Model(model.User{}).Where("user_name = ?", user.UserName).Updates(updMap)
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
}

func UserSelectByUserName(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	var result UserNameResult
	model.Db.Where("user_name = ?", claims["UserName"].(string)).First(&result.User)
	result.User.Password = ""

	roleSql := "select r.id as id, r.name, r.describe\n" +
		"from user_roles ur \n" +
		"inner join users u on ur.user_id = u.id\n" +
		"inner join roles r on ur.role_id = r.id\n" +
		"where u.user_name = ?"
	model.Db.Raw(roleSql, claims["UserName"].(string)).Scan(&result.Role)

	permissionsSql := "select r.id as role_id, m.name as permission_id, rp.action_data as actions, rp.action_list\n" +
		"from user_roles ur \n" +
		"inner join users u on ur.user_id = u.id\n" +
		"inner join roles r on ur.role_id = r.id \n" +
		"inner join role_permissions rp on ur.role_id = rp.role_id\n" +
		"inner join menus m on rp.menu_id = m.id\n" +
		"where u.user_name = ?"
	model.Db.Raw(permissionsSql, claims["UserName"].(string)).Scan(&result.Role.Perms)

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "success",
		"data": &result,
		"err":  "",
	})
}
