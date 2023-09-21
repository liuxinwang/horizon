package tasks

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"horizon/model"
	"strings"
)

var Cron = cron.New()

func InitTasks() {
	// 清除task数据，用于重新生成
	model.Db.Where("id > ?", 0).Delete(&model.Task{})

	// 循环instances生成cron
	rows, _ := model.Db.Model(&model.Instance{}).Where("status = ? AND insp_status = ? ", "Running", "Enabled").Rows()
	defer rows.Close()
	for rows.Next() {
		var instance model.Instance
		// ScanRows 方法用于将一行记录扫描至结构体
		model.Db.ScanRows(rows, &instance)
		// 每天00:10:00并发巡检
		//Cron.AddFunc("10 0 * * *", func() { InspTaskRunning(instance.InstId) })
		entryId, _ := Cron.AddFunc("10 0 * * *", func() { InspTaskRunning(instance.InstId) })
		// 记录任务到库
		model.Db.Create(&model.Task{
			ID:   uint(entryId),
			Name: fmt.Sprintf("inst-insp-%s", instance.InstId),
			Args: strings.Join([]string{instance.InstId}, ","),
		})
	}
	// 开始cron
	Cron.Start()

	// 开始定时工单任务
	go WorkflowTaskRunning()
}
