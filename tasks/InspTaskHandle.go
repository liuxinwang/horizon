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
			availableResult, _ := promApi.MetricMemoryAvailable()
			score.ScoreMemory(metric, result, availableResult)
		case model.SwapUse:
			result, _ = promApi.MetricSwap()
			score.ScoreSwap(metric, result)
		case model.DiskUtilization:
			result, _ = promApi.MetricDisk()
			resultAvailable, _ := promApi.MetricDiskAvailable()
			resultGrowth, _ := promApi.MetricDisk7DayAverageDailyGrowth()
			score.ScoreDisk(metric, resultAvailable, resultGrowth)
		case model.IOPSUtilization:
			result, _ = promApi.MetricIOPS()
			score.ScoreIOPS(metric, result)
		case model.Deadlock:
			result, _ = promApi.MetricDeadlock()
			score.ScoreDeadlock(metric, result)
		case model.SlowSQLNum:
			result, _ = promApi.MetricSlowSQLNum()
			slowTotalNum := promApi.MetricSlowSQLTotalNum()
			score.ScoreSlowSQLNum(metric, slowTotalNum)
		case model.IncrementIdOverflow:
			result, _ = promApi.MetricIncrementIdOverflow()
			score.ScoreIncrementIdOverflow(metric, result)
		case model.LockWait:
			result, _ = promApi.MetricLockWait()
			score.ScoreLockWait(metric, result)
		case model.BigTableNum:
			result, _ = promApi.MetricBigTableNum()
			score.ScoreBigTableNum(metric, result)
		case model.ThreadsRunningNum:
			result, _ = promApi.MetricThreadsRunningNum()
			score.ScoreThreadsRunningNum(metric, result)
		case model.ThreadsConnected:
			result, _ = promApi.MetricThreadsConnected()
			score.ScoreThreadsConnected(metric, result)
		case model.IBPCacheHitsRate:
			result, _ = promApi.MetricIBPCacheHitsRate()
			score.ScoreIBPCacheHitsRate(metric, result)
		case model.QPS:
			result, _ = promApi.MetricQPS()
		case model.TPS:
			result, _ = promApi.MetricTPS()
		case model.HighRiskAccount:
			result, _ = promApi.MetricHighRiskAccount()
			score.ScoreHighRiskAccount(metric, result)
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
	}
	// 计算总评分
	totalScore := calculateTotalScore(inspection.InspId)
	model.Db.Model(&inspection).Update("score", totalScore)
	// 评健康等级
	healthLevel := calculateScoreLevel(totalScore)
	model.Db.Model(&inspection).Update("level", healthLevel)
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
