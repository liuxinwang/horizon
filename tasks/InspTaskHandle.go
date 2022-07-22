package tasks

import (
	"fmt"
	"gorm.io/gorm"
	"horizon/model"
	"horizon/utils"
	"time"
)

// InspTaskRunning 实例巡检入口
func InspTaskRunning(instId string) {
	var instance model.Instance
	model.Db.Where("inst_id = ?", instId).First(&instance)
	result, inspection := initInspection(instance)
	if result.Error != nil {
		return
	}
	// 初始化prometheus client
	yday := time.Now().AddDate(0, 0, -1)
	startTime := time.Date(yday.Year(), yday.Month(), yday.Day(), 0, 0, 0, 0, yday.Location())
	endTime := time.Date(yday.Year(), yday.Month(), yday.Day(), 23, 59, 59, 999, yday.Location())
	promApi := utils.Prom{
		Api:               utils.PromAPI(),
		StartTime:         startTime,
		EndTime:           endTime,
		Instance:          instance,
		NodeExporterInst:  fmt.Sprintf("%s:9100", instance.Ip),
		MySQLExporterInst: fmt.Sprintf("%s:1%d", instance.Ip, instance.Port),
	}
	var score Score
	score.inspection = *inspection

	// 遍历指标
	rows, _ := model.Db.Model(&model.Metric{}).Rows()
	defer rows.Close()

	for rows.Next() {
		var metric model.Metric
		// ScanRows 方法用于将一行记录扫描至结构体
		model.Db.ScanRows(rows, &metric)
		var result []byte
		// 获取指标数据
		switch metric.Key {
		case model.CpuUtilization:
			result, _ = promApi.MetricCPU()
			score.ScoreCPU(metric, result)
		case model.MemoryUtilization:
			result, _ = promApi.MetricMemory()
		case model.SwapUse:
			result, _ = promApi.MetricSwap()
		case model.DiskUtilization:
			result, _ = promApi.MetricDisk()
		case model.IOPSUtilization:
			result, _ = promApi.MetricIOPS()
		case model.Deadlock:
			result, _ = promApi.MetricDeadlock()
		case model.SlowSQLNum:
			result, _ = promApi.MetricSlowSQLNum()
		case model.IncrementIdOverflow:
			result, _ = promApi.MetricIncrementIdOverflow()
		case model.LockWait:
			result, _ = promApi.MetricLockWait()
		case model.BigTableNum:
			result, _ = promApi.MetricBigTableNum()
		case model.ThreadsRunningNum:
			result, _ = promApi.MetricThreadsRunningNum()
		case model.ThreadsConnected:
			result, _ = promApi.MetricThreadsConnected()
		case model.IBPCacheHitsRate:
			result, _ = promApi.MetricIBPCacheHitsRate()
		case model.QPS:
			result, _ = promApi.MetricQPS()
		case model.TPS:
			result, _ = promApi.MetricTPS()
		case model.HighRiskAccount:
			result, _ = promApi.MetricHighRiskAccount()
		case model.HAStatus:
			// result, _ = promApi.MetricHAStatus()
			continue
		case model.ReplicationStatus:
			result, _ = promApi.MetricReplicationStatus()
		case model.ReplicationDelay:
			result, _ = promApi.MetricReplicationDelay()
		case model.BackupStatus:
			// result, _ = promApi.MetricBackupStatus()
			continue
		case model.NetworkTrafficIn:
			result, _ = promApi.MetricNetworkTrafficIn()
		case model.NetworkTrafficOut:
			result, _ = promApi.MetricNetworkTrafficOut()
		default:
			fmt.Printf("default")
		}
		// 记录数据到DB
		model.Db.Create(&model.InspDetail{
			InspId: inspection.InspId, Metric: metric.Key, Result: result,
		})
		// 单指标评分
	}
	// 计算总评分
	// 评健康等级
}

// 初始化巡检记录
func initInspection(instance model.Instance) (*gorm.DB, *model.Inspection) {
	InspId := instance.InstId + "-" + time.Now().Format("20060102")
	inspection := model.Inspection{
		InspId:   InspId,
		InstId:   instance.InstId,
		InstName: instance.Name,
	}
	model.Db.Where("insp_id = ?", inspection.InspId).Delete(&model.Inspection{})
	result := model.Db.Create(&inspection)
	return result, &inspection
}
