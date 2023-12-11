package model

import (
	"fmt"
	json "github.com/json-iterator/go"
	"gorm.io/gorm"
	"horizon/config"
	"horizon/notification"
	"time"
)

type DataMigrateJob struct {
	ID           uint      `gorm:"primaryKey;comment:主键ID" json:"id"`
	Name         string    `gorm:"type:varchar(50);not null;comment:任务名称" json:"name"`
	Describe     string    `gorm:"type:varchar(255);not null;comment:描述" json:"describe"`
	SourceInstId string    `gorm:"type:varchar(20);not null;comment:源实例ID" json:"sourceInstId"`
	SourceDb     string    `gorm:"type:varchar(100);not null;comment:源数据库" json:"sourceDb"`
	TargetInstId string    `gorm:"type:varchar(20);not null;comment:目的实例ID" json:"targetInstId"`
	TargetDb     string    `gorm:"type:varchar(100);not null;comment:目的数据库" json:"targetDb"`
	Status       string    `gorm:"type:enum('NotStart', 'Running', 'Error', 'Finished');default:'NotStart';not null;comment:任务状态" json:"status"`
	UserName     string    `gorm:"type:varchar(50);not null;comment:用户名" json:"userName"`
	CreatedAt    time.Time `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
}

type DataMigrateJobDetail struct {
	ID               uint      `gorm:"primaryKey;comment:主键ID" json:"id"`
	DataMigrateJobId uint      `gorm:"type:int;comment:迁移任务主键ID" json:"dataMigrateJobId"`
	TableName        string    `gorm:"type:varchar(100);not null;comment:表名称" json:"tableName"`
	Status           string    `gorm:"type:enum('NotStart', 'Running', 'Error', 'Finished');default:'NotStart';not null;comment:同步状态" json:"status"`
	EstimateRows     uint      `gorm:"type:int;comment:预估行数" json:"estimateRows"`
	CompletedRows    uint      `gorm:"type:int;comment:已完成行数" json:"completedRows"`
	ErrorMsg         string    `gorm:"type:varchar(1000);not null;comment:错误信息" json:"errorMsg"`
	CreatedAt        time.Time `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt        time.Time `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
}

func (dm *DataMigrateJob) AfterCreate(tx *gorm.DB) (err error) {
	title := "数据迁移工单通知"
	var dataMigrateJobUser User
	result := tx.Where("user_name = ?", dm.UserName).First(&dataMigrateJobUser)
	if result.Error != nil {
		return result.Error
	}

	if dm.Status == "NotStart" {
		markdown := notification.DingContentMarkdown{
			Title: title,
			Text: fmt.Sprintf(
				"用户：%s 提交的 数据迁移【%s】，正在等待审批，请确认！[审核](%s/dataManger/dataMigrateJobDetail/%d)",
				dataMigrateJobUser.NickName, dm.Name,
				config.Conf.General.HomeAddress, dm.ID),
		}
		at := notification.DingContentAt{
			AtMobiles: []string{},
			AtUserIds: []string{},
			IsAtAll:   false,
		}
		content := notification.DingContent{
			MsgType:         "markdown",
			ContentMarkdown: markdown,
			ContentAt:       at,
		}
		marshal, err := json.Marshal(content)
		if err != nil {
			return err
		}
		notification.SendDingDing(string(marshal))
	}
	return nil
}

func (dm *DataMigrateJob) AfterUpdate(tx *gorm.DB) (err error) {
	title := "数据迁移工单通知"
	var dataMigrateJobUser User
	result := tx.Where("user_name = ?", dm.UserName).First(&dataMigrateJobUser)
	if result.Error != nil {
		return result.Error
	}

	if dm.Status == "Running" {
		markdown := notification.DingContentMarkdown{
			Title: title,
			Text: fmt.Sprintf(
				"用户：@%s 提交的 数据迁移【%s】，正在执行。[工单详情](%s/dataManger/dataMigrateJobDetail/%d)",
				dataMigrateJobUser.Phone, dm.Name, config.Conf.General.HomeAddress, dm.ID),
		}
		at := notification.DingContentAt{
			AtMobiles: []string{dataMigrateJobUser.Phone},
			AtUserIds: []string{},
			IsAtAll:   false,
		}
		content := notification.DingContent{
			MsgType:         "markdown",
			ContentMarkdown: markdown,
			ContentAt:       at,
		}
		marshal, err := json.Marshal(content)
		if err != nil {
			return err
		}
		notification.SendDingDing(string(marshal))
	}
	if dm.Status == "Finished" {
		markdown := notification.DingContentMarkdown{
			Title: title,
			Text: fmt.Sprintf(
				"用户：@%s 提交的 数据迁移【%s】，已执行成功。[工单详情](%s/dataManger/dataMigrateJobDetail/%d)",
				dataMigrateJobUser.Phone, dm.Name, config.Conf.General.HomeAddress, dm.ID),
		}
		at := notification.DingContentAt{
			AtMobiles: []string{dataMigrateJobUser.Phone},
			AtUserIds: []string{},
			IsAtAll:   false,
		}
		content := notification.DingContent{
			MsgType:         "markdown",
			ContentMarkdown: markdown,
			ContentAt:       at,
		}
		marshal, err := json.Marshal(content)
		if err != nil {
			return err
		}
		notification.SendDingDing(string(marshal))
	}
	if dm.Status == "Error" {
		markdown := notification.DingContentMarkdown{
			Title: title,
			Text: fmt.Sprintf(
				"用户：@%s 提交的 SQL工单【%s】，执行失败！[工单详情](%s/dataManger/dataMigrateJobDetail/%d)",
				dataMigrateJobUser.Phone, dm.Name, config.Conf.General.HomeAddress, dm.ID),
		}
		at := notification.DingContentAt{
			AtMobiles: []string{dataMigrateJobUser.Phone},
			AtUserIds: []string{},
			IsAtAll:   false,
		}
		content := notification.DingContent{
			MsgType:         "markdown",
			ContentMarkdown: markdown,
			ContentAt:       at,
		}
		marshal, err := json.Marshal(content)
		if err != nil {
			return err
		}
		notification.SendDingDing(string(marshal))
	}
	return
}
