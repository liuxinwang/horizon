package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"horizon/model"
	"net/http"
)

// RolePermissionInsert 新增角色权限
func RolePermissionInsert(c *gin.Context) {
	// 参数映射到对象
	type permission struct {
		model.Menu
		Selected     datatypes.JSON `json:"selected"`
		SelectedData datatypes.JSON `json:"selectedData"`
	}
	var permBody struct {
		Role        model.Role   `json:"role"`
		Permissions []permission `json:"permissions"`
	}
	if err := c.ShouldBind(&permBody); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	var rolePermissions []model.RolePermission
	parentIds := make(map[uint]uint)
	for _, permission := range permBody.Permissions {
		if len(permission.Selected) > 0 && permission.Selected.String() != "[]" {
			// 添加选择的节点的父节点
			if _, ok := parentIds[permission.ParentId]; !ok {
				parentIds[permission.ParentId] = permission.ParentId
				rolePermissions = append(rolePermissions, model.RolePermission{
					RoleId:     permBody.Role.ID,
					MenuId:     permission.ParentId,
					ActionData: datatypes.JSON("[]"),
					ActionList: datatypes.JSON("[]"),
				})
			}
			// 添加选择的节点
			rolePermissions = append(rolePermissions, model.RolePermission{
				RoleId:     permBody.Role.ID,
				MenuId:     permission.ID,
				ActionData: permission.SelectedData,
				ActionList: permission.Selected,
			})
		}
	}

	model.Db.Debug().Delete(&model.RolePermission{}, "role_id = ?", permBody.Role.ID)
	result := model.Db.Debug().Create(&rolePermissions)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
	}
}
