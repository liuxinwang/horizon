package service

import (
	"github.com/gin-gonic/gin"
	"horizon/model"
	"horizon/utils"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// ProjectSelectByList 查询实例列表
func ProjectSelectByList(c *gin.Context) {
	// 变量初始化
	pageNo, _ := strconv.Atoi(c.Query("pageNo"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	var Db = model.Db
	var projects []model.Project
	var totalCount int64
	data := gin.H{"totalCount": 0, "data": &[]model.Project{}, "pageNo": pageNo, "pageSize": pageSize, "totalPage": 0}

	// 查询条件处理
	if name, isExist := c.GetQuery("Name"); isExist == true && strings.Trim(name, " ") != "" {
		Db = Db.Where("name like ?", "%"+name+"%")
	}
	if status, isExist := c.GetQuery("Status"); isExist == true && status != "0" {
		Db = Db.Where("status = ?", status)
	}

	// 执行查询
	Db.Preload("ProjectDatasources").
		Preload("ProjectUsers").
		Order("created_at desc").Limit(pageSize).Offset((pageNo-1)*pageSize - 1).Find(&projects)
	Db.Model(&model.Project{}).Count(&totalCount)

	// 处理结果集并返回
	totalPage := math.Ceil(float64(totalCount) / float64(pageSize))
	if totalCount > 0 {
		data = gin.H{"totalCount": totalCount, "data": &projects, "pageNo": pageNo, "pageSize": pageSize, "totalPage": totalPage}
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": data, "err": ""})
}

// ProjectSelectByInstId 查看实例信息
func ProjectSelectByInstId(c *gin.Context) {
	var Db = model.Db
	var project model.Project
	// 执行查询
	Db.Where("proj_id = ?", c.Param("instId")).Find(&project)
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": &project, "err": ""})
}

// ProjectInsert 新增
func ProjectInsert(c *gin.Context) {
	// 参数映射到对象
	var project model.Project
	if err := c.ShouldBind(&project); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	// 获取ID
	project.ProjId = utils.GenerateId(&project)
	result := model.Db.Create(&project)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
	}
}

func ProjectUpdate(c *gin.Context) {
	// 参数映射到对象
	var project model.Project
	if err := c.ShouldBind(&project); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	// 执行更新
	updMap := map[string]interface{}{"name": project.Name, "describe": project.Describe}
	model.Db.Model(model.Project{}).Where("proj_id = ?", project.ProjId).Updates(updMap)
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
}

func ProjectDelete(c *gin.Context) {
	id := c.Param("id")
	result := model.Db.Delete(&model.Project{}, id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": "项目不存在"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
	}
}

func ProjectResourceConfigInsert(c *gin.Context) {
	var permBody struct {
		model.Project
		UsersRole   []string `json:"projUsersRole"`
		Users       []string `json:"projUsers"`
		Datasources []string `json:"datasources"`
	}
	if err := c.ShouldBind(&permBody); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	tx := model.Db.Begin()
	// update project workflow template
	tx.Model(&model.Project{}).
		Where("proj_id = ?", permBody.ProjId).
		UpdateColumn("workflow_template_code", permBody.WorkflowTemplateCode)
	// first delete and add ProjectDatasources
	result := tx.Where("proj_id = ?", permBody.Project.ProjId).Delete(&model.ProjectDatasource{})
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
		tx.Rollback()
		return
	}
	for _, instId := range permBody.Datasources {
		var projectDatasource model.ProjectDatasource
		projectDatasource.ProjId = permBody.Project.ProjId
		projectDatasource.InstId = instId
		result := tx.Create(&projectDatasource)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
			tx.Rollback()
			return
		}
	}
	// first delete and add ProjectUser
	result = tx.Where("proj_id = ?", permBody.Project.ProjId).Delete(&model.ProjectUser{})
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
		tx.Rollback()
		return
	}
	for i, userName := range permBody.Users {
		var projectUser model.ProjectUser
		projectUser.ProjId = permBody.Project.ProjId
		projectUser.RoleId = permBody.UsersRole[i]
		projectUser.UserName = userName
		result := tx.Create(&projectUser)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
}

// ProjectRoleSelectByList 查询列表
func ProjectRoleSelectByList(c *gin.Context) {
	// 变量初始化
	pageNo, _ := strconv.Atoi(c.Query("pageNo"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	var Db = model.Db
	var projectRoles []model.ProjectRole
	var totalCount int64
	data := gin.H{"totalCount": 0, "data": &[]model.Project{}, "pageNo": pageNo, "pageSize": pageSize, "totalPage": 0}

	// 执行查询
	Db.Order("created_at").Limit(pageSize).Offset((pageNo-1)*pageSize - 1).Find(&projectRoles)
	Db.Model(&model.ProjectRole{}).Count(&totalCount)

	// 处理结果集并返回
	totalPage := math.Ceil(float64(totalCount) / float64(pageSize))
	if totalCount > 0 {
		data = gin.H{"totalCount": totalCount, "data": &projectRoles, "pageNo": pageNo, "pageSize": pageSize, "totalPage": totalPage}
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": data, "err": ""})
}

// ProjectUserSelectByUserName 查询列表
func ProjectUserSelectByUserName(c *gin.Context) {
	type result struct {
		ProjId string `json:"projId"`
		Name   string `json:"name"`
	}
	var results []result
	// 执行查询
	model.Db.Select("projects.*").Model(&model.ProjectUser{}).
		Joins("inner join projects on project_users.proj_id = projects.proj_id").
		Where("user_name = ?", c.Param("userName")).Scan(&results)
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": &results, "err": ""})
}

// ProjectDataSourceSelectByProjId 查询列表
func ProjectDataSourceSelectByProjId(c *gin.Context) {
	type result struct {
		InstId string `json:"instId"`
		Name   string `json:"name"`
	}
	var results []result
	// 执行查询
	model.Db.Select("instances.*").Model(&model.ProjectDatasource{}).
		Joins("inner join instances on project_datasources.inst_id = instances.inst_id").
		Where("proj_id = ?", c.Param("projId")).Scan(&results)
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": &results, "err": ""})
}
